package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 需要全局使用的变量都在这里声明，方便其他包调用
var (
	DB     *gorm.DB
	Config config
	// Logger zap的标准logger，速度更快，但是输入麻烦，用于取代gin的logger
	Logger *zap.Logger
	// SugaredLogger zap的加糖logger，速度慢一点点，但是输入方便，自己用
	// https://pkg.go.dev/go.uber.org/zap#SugaredLogger
	SugaredLogger *zap.SugaredLogger
	v             = viper.New()
)

// 这层只是中间的汇总层，只是包内引用、不展示，所以小写
type config struct {
	AppConfig
	DBConfig
	JWTConfig
	LogConfig
	UploadConfig
	EmailConfig
	PagingConfig
	RateLimitConfig
}

type AppConfig struct {
	AppMode  string `yaml:"app-mode"`
	HttpPort string
}

type DBConfig struct {
	DbHost         string
	DbPort         string
	DbName         string
	DbUsername     string
	DbPassword     string
	DSN            string   // Data Source Name 数据库连接字符串
	OmittedColumns []string //绝对不传给前端的数据库字段名
}

type JWTConfig struct {
	SecretKey      []byte //这里不能用string，是jwt包的要求，否则报错
	ValidityPeriod int
}

type LogConfig struct {
	FileName      string
	MaxSizeForLog int
	MaxBackup     int
	MaxAge        int
	Compress      bool
}

type UploadConfig struct {
	FullPath         string
	MaxSizeForUpload int64
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
	Config.DBConfig.OmittedColumns = v.GetStringSlice("database.omitted-columns")

	//配置里的密钥是string类型，jwt要求为[]byte类型，必须转换后才能使用
	Config.JWTConfig.SecretKey = []byte(v.GetString("jwt.secret-key"))
	Config.JWTConfig.ValidityPeriod = v.GetInt("jwt.validity-period")

	Config.LogConfig.FileName = v.GetString("log.log-path") + "/status.log"
	Config.LogConfig.MaxSizeForLog = v.GetInt("log.log-max-size")
	Config.LogConfig.MaxBackup = v.GetInt("log.log-max-backup")
	Config.LogConfig.MaxAge = v.GetInt("log.log-max-age")
	Config.LogConfig.Compress = v.GetBool("log.log-compress")

	Config.UploadConfig.FullPath = v.GetString("upload-files.full-path") + "/"
	Config.UploadConfig.MaxSizeForUpload = v.GetInt64("upload-files.max-size") << 20

	Config.EmailConfig.OutgoingMailServer = v.GetString("email.outgoing-mail-server")
	Config.EmailConfig.Port = v.GetInt("email.port")
	Config.EmailConfig.Account = v.GetString("email.account")
	Config.EmailConfig.Password = v.GetString("email.password")

	Config.PagingConfig.DefaultPageSize = v.GetInt("paging.default-page-size")
	Config.PagingConfig.MaxPageSize = v.GetInt("paging.max-page-size")

	Config.RateLimitConfig.Limit = v.GetFloat64("rate-limit.limit")
	Config.RateLimitConfig.Burst = v.GetInt("rate-limit.burst")
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
