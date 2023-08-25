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

// 用于确定角色的数据范围
const (
	HisOrganization = iota + 1
	HisOrganizationAndInferiors
	AllOrganization
	CustomOrganization
)

// 这层只是中间的汇总层，只是包内引用、不展示，所以小写
type config struct {
	AppConfig
	DBConfig, DBConfigForLvmin, DBConfigForOldPmis DBConfig
	JWTConfig
	LogConfig
	UploadConfig
	DownloadConfig
	EmailConfig
	PagingConfig
	RateLimitConfig
	CaptchaConfig
	OssConfig
	LDAPConfig
}

type AppConfig struct {
	AppMode  string
	HttpPort string
}

type DBConfig struct {
	DbHost     string
	DbPort     string
	DbName     string
	DbUsername string
	DbPassword string
	DSN        string // Param Source Name 数据库连接字符串
}

type JWTConfig struct {
	SecretKey    string
	ValidityDays int
	Issuer       string
}

type LogConfig struct {
	FileName      string
	MaxSizeForLog int
	MaxBackup     int
	MaxAge        int
	Compress      bool
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

type OssConfig struct {
	Type string
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
	Config.AppConfig.AppMode = v.GetString("app.app-mode")
	allowedAppMode := []string{"debug", "test", "release"}
	if !isInSlice(Config.AppConfig.AppMode, allowedAppMode) {
		Config.AppConfig.AppMode = "debug"
	}
	Config.AppConfig.HttpPort = v.GetString("app.http-port")

	Config.DBConfig.DbHost = v.GetString("database.db-host")
	Config.DBConfig.DbPort = v.GetString("database.db-port")
	Config.DBConfig.DbName = v.GetString("database.db-name")
	Config.DBConfig.DbUsername = v.GetString("database.db-username")
	Config.DBConfig.DbPassword = v.GetString("database.db-password")
	Config.DBConfig.DSN = "sqlserver://" + Config.DBConfig.DbUsername +
		":" + Config.DBConfig.DbPassword + "@" + Config.DBConfig.DbHost +
		":" + Config.DBConfig.DbPort + "?database=" + Config.DBConfig.DbName

	Config.DBConfigForLvmin.DbHost = v.GetString("database2.db-host")
	Config.DBConfigForLvmin.DbPort = v.GetString("database2.db-port")
	Config.DBConfigForLvmin.DbName = v.GetString("database2.db-name")
	Config.DBConfigForLvmin.DbUsername = v.GetString("database2.db-username")
	Config.DBConfigForLvmin.DbPassword = v.GetString("database2.db-password")
	Config.DBConfigForLvmin.DSN = "sqlserver://" + Config.DBConfigForLvmin.DbUsername +
		":" + Config.DBConfigForLvmin.DbPassword + "@" + Config.DBConfigForLvmin.DbHost +
		":" + Config.DBConfigForLvmin.DbPort + "?database=" + Config.DBConfigForLvmin.DbName +
		"&encrypt=disable" //老版本数据库不支持加密连接，不加这个会报错

	Config.DBConfigForOldPmis.DbHost = v.GetString("database3.db-host")
	Config.DBConfigForOldPmis.DbPort = v.GetString("database3.db-port")
	Config.DBConfigForOldPmis.DbName = v.GetString("database3.db-name")
	Config.DBConfigForOldPmis.DbUsername = v.GetString("database3.db-username")
	Config.DBConfigForOldPmis.DbPassword = v.GetString("database3.db-password")
	Config.DBConfigForOldPmis.DSN = "sqlserver://" + Config.DBConfigForOldPmis.DbUsername +
		":" + Config.DBConfigForOldPmis.DbPassword + "@" + Config.DBConfigForOldPmis.DbHost +
		":" + Config.DBConfigForOldPmis.DbPort + "?database=" + Config.DBConfigForOldPmis.DbName

	Config.JWTConfig.SecretKey = v.GetString("jwt.secret-key")
	Config.JWTConfig.ValidityDays = v.GetInt("jwt.validity-days")
	Config.JWTConfig.Issuer = v.GetString("jwt.issuer")

	Config.LogConfig.FileName = v.GetString("log.log-path") + "/status.log"
	Config.LogConfig.MaxSizeForLog = v.GetInt("log.log-max-size")
	Config.LogConfig.MaxBackup = v.GetInt("log.log-max-backup")
	Config.LogConfig.MaxAge = v.GetInt("log.log-max-age")
	Config.LogConfig.Compress = v.GetBool("log.log-compress")

	Config.UploadConfig.StoragePath = v.GetString("upload.storage-path") + "/"
	Config.UploadConfig.MaxSize = v.GetInt64("upload.max-size") << 20

	Config.DownloadConfig.RelativePath = v.GetString("download.relative-path") + "/"
	var err error
	Config.DownloadConfig.LocalIP, err = GetLocalIP()
	if err != nil {
		SugaredLogger.Panicln("获取本地IP失败：", err)
	}
	Config.DownloadConfig.FullPath = "http://" + Config.DownloadConfig.LocalIP +
		":" + Config.AppConfig.HttpPort + Config.DownloadConfig.RelativePath

	Config.EmailConfig.OutgoingMailServer = v.GetString("email.outgoing-mail-server")
	Config.EmailConfig.Port = v.GetInt("email.port")
	Config.EmailConfig.Account = v.GetString("email.account")
	Config.EmailConfig.Password = v.GetString("email.password")

	Config.PagingConfig.DefaultPageSize = v.GetInt("paging.default-page-size")
	Config.PagingConfig.MaxPageSize = v.GetInt("paging.max-page-size")

	Config.RateLimitConfig.Limit = v.GetFloat64("rate-limit.limit")
	Config.RateLimitConfig.Burst = v.GetInt("rate-limit.burst")

	Config.CaptchaConfig.DigitLength = v.GetInt("captcha.digit-length")
	Config.CaptchaConfig.ImageWidth = v.GetInt("captcha.image-width")
	Config.CaptchaConfig.ImageHeight = v.GetInt("captcha.image-height")
	Config.CaptchaConfig.MaxSkew = v.GetFloat64("captcha.max-skew")
	Config.CaptchaConfig.DotCount = v.GetInt("captcha.dot-count")

	Config.OssConfig.Type = v.GetString("oss.type")

	Config.LDAPConfig.Server = v.GetString("ldap.server")
	Config.LDAPConfig.BaseDN = v.GetString("ldap.base-dn")
	Config.LDAPConfig.Filter = v.GetString("ldap.filter")
	Config.LDAPConfig.Suffix = v.GetString("ldap.suffix")
	Config.LDAPConfig.Account = v.GetString("ldap.account")
	Config.LDAPConfig.Password = v.GetString("ldap.password")
	Config.LDAPConfig.PermittedOUs = v.GetStringSlice("ldap.permitted-OUs")
	Config.LDAPConfig.Attributes = v.GetStringSlice("ldap.attributes")
}

// byte 是 uint8 的别名,rune 是 int32 的别名
type typeForSliceComparing interface {
	bool | string | int | int64 | int32 | int16 | int8 |
		uint | uint64 | uint32 | uint16 | uint8 |
		float64 | float32
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
