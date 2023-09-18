package global

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net"
	"strings"
)

// 需要全局使用的变量都在这里声明，方便其他包调用
var (
	DB           *gorm.DB //自身主数据库
	DBForLvmin   *gorm.DB //率敏的数据库
	DBForOldPmis *gorm.DB //老的pmis数据库
	Config       config

	// Logger zap的标准logger，速度更快，但是输入麻烦，用于取代gin的logger
	Logger *zap.Logger
	// SugaredLogger zap的加糖logger，速度慢一点点，但是输入方便，自己用
	// https://pkg.go.dev/go.uber.org/zap#SugaredLogger
	SugaredLogger *zap.SugaredLogger
	v             = viper.New()
)

// 这层只是中间的汇总层，只是包内引用、不展示，所以小写
type config struct {
	App                          AppConfig
	Db, DbForLvmin, DbForOldPmis DBConfig
	Jwt                          JWTConfig
	Log                          LogConfig
	Upload                       UploadConfig
	Download                     DownloadConfig
	Email                        EmailConfig
	Paging                       PagingConfig
	RateLimit                    RateLimitConfig
	Captcha                      CaptchaConfig
	Ldap                         LDAPConfig
	ExchangeRate                 ExchangeRateConfig
}

type AppConfig struct {
	Mode     string
	HttpPort string
}

type DBConfig struct {
	Host     string
	Port     string
	DbName   string
	Username string
	Password string
	DSN      string // Param Source Name 数据库连接字符串
}

type JWTConfig struct {
	SecretKey    string
	ValidityDays int
	Issuer       string
}

type LogConfig struct {
	FileName  string
	MaxSize   int
	MaxBackup int
	MaxAge    int
	Compress  bool
}

type UploadConfig struct {
	StoragePath string
	MaxSize     int64
}

type DownloadConfig struct {
	LocalIP      string
	RelativePath string
	FullPath     string
}

type EmailConfig struct {
	OutgoingMailServer string
	Port               int
	Account            string
	Password           string
}

type PagingConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

type RateLimitConfig struct {
	Limit float64
	Burst int
}

type CaptchaConfig struct {
	DigitLength int
	ImageWidth  int
	ImageHeight int
	MaxSkew     float64
	DotCount    int
}

type LDAPConfig struct {
	Server       string
	BaseDN       string
	Filter       string
	Suffix       string
	Account      string
	Password     string
	PermittedOUs []string
	Attributes   []string
}

type ExchangeRateConfig struct {
	USD float64 //美元
	EUR float64 //欧元
	HKD float64 //港币
	SGD float64 //新加坡元
	MYR float64 //马来西亚币
}

func InitConfig() {
	v.AddConfigPath("./config/") //告诉viper，配置文件的路径在哪
	v.SetConfigName("config")    //告诉viper，配置文件的前缀是什么
	v.SetConfigType("yaml")      //告诉viper，配置文件的后缀是什么
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Config file error: %w \n", err))
	}

	//配置文件热更新
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
		loadConfig()
	})

	v.WatchConfig()
	loadConfig()
}

func loadConfig() {
	Config.App.Mode = v.GetString("App.App-mode")
	allowedAppMode := []string{"debug", "test", "release"}
	if !isInSlice(Config.App.Mode, allowedAppMode) {
		Config.App.Mode = "debug"
	}
	Config.App.HttpPort = v.GetString("App.http-port")

	Config.Db.Host = v.GetString("database.db-host")
	Config.Db.Port = v.GetString("database.db-port")
	Config.Db.DbName = v.GetString("database.db-name")
	Config.Db.Username = v.GetString("database.db-username")
	Config.Db.Password = v.GetString("database.db-password")
	Config.Db.DSN = "sqlserver://" + Config.Db.Username +
		":" + Config.Db.Password + "@" + Config.Db.Host +
		":" + Config.Db.Port + "?database=" + Config.Db.DbName

	Config.DbForLvmin.Host = v.GetString("database2.db-host")
	Config.DbForLvmin.Port = v.GetString("database2.db-port")
	Config.DbForLvmin.DbName = v.GetString("database2.db-name")
	Config.DbForLvmin.Username = v.GetString("database2.db-username")
	Config.DbForLvmin.Password = v.GetString("database2.db-password")
	Config.DbForLvmin.DSN = "sqlserver://" + Config.DbForLvmin.Username +
		":" + Config.DbForLvmin.Password + "@" + Config.DbForLvmin.Host +
		":" + Config.DbForLvmin.Port + "?database=" + Config.DbForLvmin.DbName +
		"&encrypt=disable" //老版本数据库不支持加密连接，不加这个会报错

	Config.DbForOldPmis.Host = v.GetString("database3.db-host")
	Config.DbForOldPmis.Port = v.GetString("database3.db-port")
	Config.DbForOldPmis.DbName = v.GetString("database3.db-name")
	Config.DbForOldPmis.Username = v.GetString("database3.db-username")
	Config.DbForOldPmis.Password = v.GetString("database3.db-password")
	Config.DbForOldPmis.DSN = "sqlserver://" + Config.DbForOldPmis.Username +
		":" + Config.DbForOldPmis.Password + "@" + Config.DbForOldPmis.Host +
		":" + Config.DbForOldPmis.Port + "?database=" + Config.DbForOldPmis.DbName

	Config.Jwt.SecretKey = v.GetString("Jwt.secret-key")
	Config.Jwt.ValidityDays = v.GetInt("Jwt.validity-days")
	Config.Jwt.Issuer = v.GetString("Jwt.issuer")

	Config.Log.FileName = v.GetString("Log.Log-path") + "/status.Log"
	Config.Log.MaxSize = v.GetInt("Log.Log-max-size")
	Config.Log.MaxBackup = v.GetInt("Log.Log-max-backup")
	Config.Log.MaxAge = v.GetInt("Log.Log-max-age")
	Config.Log.Compress = v.GetBool("Log.Log-compress")

	Config.Upload.StoragePath = v.GetString("Upload.storage-path") + "/"
	Config.Upload.MaxSize = v.GetInt64("Upload.max-size") << 20

	Config.Download.RelativePath = v.GetString("Download.relative-path") + "/"
	var err error
	Config.Download.LocalIP, err = GetLocalIP()
	if err != nil {
		SugaredLogger.Panicln("获取本地IP失败：", err)
	}
	Config.Download.FullPath = "http://" + Config.Download.LocalIP +
		":" + Config.App.HttpPort + Config.Download.RelativePath

	Config.Email.OutgoingMailServer = v.GetString("Email.outgoing-mail-server")
	Config.Email.Port = v.GetInt("Email.port")
	Config.Email.Account = v.GetString("Email.account")
	Config.Email.Password = v.GetString("Email.password")

	Config.Paging.DefaultPageSize = v.GetInt("Paging.default-page-size")
	Config.Paging.MaxPageSize = v.GetInt("Paging.max-page-size")

	Config.RateLimit.Limit = v.GetFloat64("rate-limit.limit")
	Config.RateLimit.Burst = v.GetInt("rate-limit.burst")

	Config.Captcha.DigitLength = v.GetInt("Captcha.digit-length")
	Config.Captcha.ImageWidth = v.GetInt("Captcha.image-width")
	Config.Captcha.ImageHeight = v.GetInt("Captcha.image-height")
	Config.Captcha.MaxSkew = v.GetFloat64("Captcha.max-skew")
	Config.Captcha.DotCount = v.GetInt("Captcha.dot-count")

	Config.Ldap.Server = v.GetString("Ldap.server")
	Config.Ldap.BaseDN = v.GetString("Ldap.base-dn")
	Config.Ldap.Filter = v.GetString("Ldap.filter")
	Config.Ldap.Suffix = v.GetString("Ldap.suffix")
	Config.Ldap.Account = v.GetString("Ldap.account")
	Config.Ldap.Password = v.GetString("Ldap.password")
	Config.Ldap.PermittedOUs = v.GetStringSlice("Ldap.permitted-OUs")
	Config.Ldap.Attributes = v.GetStringSlice("Ldap.attributes")

	Config.ExchangeRate.USD = v.GetFloat64("exchange-rate.USD")
	Config.ExchangeRate.EUR = v.GetFloat64("exchange-rate.EUR")
	Config.ExchangeRate.HKD = v.GetFloat64("exchange-rate.HKD")
	Config.ExchangeRate.SGD = v.GetFloat64("exchange-rate.SGD")
	Config.ExchangeRate.MYR = v.GetFloat64("exchange-rate.MYR")

}

// byte 是 uint8 的别名,rune 是 int32 的别名
type typeForSliceComparing interface {
	bool | string | int | int64 | int32 | int16 | int8 | uint | uint64 |
		uint32 | uint16 | uint8 | float64 | float32
}

// isInSlice 这里使用了泛型，至少需要1.18版本以上
// 校验单个内容是否包含在切片中
func isInSlice[T typeForSliceComparing](element T, slice []T) bool {
	for _, v := range slice {
		if element == v {
			return true
		}
	}
	return false
}

func GetLocalIP() (ip string, err error) {
	var localIps []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			Addresses, _ := netInterfaces[i].Addrs()
			for _, address := range Addresses {
				if IpNet, ok := address.(*net.IPNet); ok && !IpNet.IP.IsLoopback() {
					if IpNet.IP.To4() != nil {
						localIps = append(localIps, IpNet.IP.String())
					}
				}
			}
		}
	}

	for i := range localIps {
		if strings.Contains(localIps[i], "10.") {
			return localIps[i], nil
		}
	}

	return "", errors.New("未找到本机ip，请联系网络管理员")
}
