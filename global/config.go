package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//这里是全局变量，给其他所有包引用
var (
	DB            *gorm.DB
	Config        config
	Logger        *zap.Logger        //zap的标准logger，速度更快，但是输入麻烦，用于取代gin的logger
	SugaredLogger *zap.SugaredLogger // zap的加糖logger，速度慢一点点，但是输入方便，自己用
	v1            = viper.New()
)

//这层只是中间的汇总层，只是包内引用、不展示，所以小写
type config struct {
	APPConfig
	DBConfig
	JWTConfig
	LogConfig
	UploadConfig
	EmailConfig
	PagingConfig
	SqlConfig
}

type APPConfig struct {
	AppMode  string
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

type SqlConfig struct {
	SqlStatement string
}

func Init() {
	v1.AddConfigPath("./config/") //告诉viper，配置文件的路径在哪
	v1.SetConfigName("config")    //告诉viper，配置文件的前缀是什么
	v1.SetConfigType("yaml")      //告诉viper，配置文件的后缀是什么
	if err := v1.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Config file error: %w \n", err))
	}
	v1.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
		loadConfig()
	})
	v1.WatchConfig()

	loadConfig()
}

func loadConfig() {
	Config.APPConfig.AppMode = v1.GetString("app.app_mode")
	Config.APPConfig.HttpPort = v1.GetString("app.http_port")

	Config.DBConfig.DbHost = v1.GetString("database.db_host")
	Config.DBConfig.DbPort = v1.GetString("database.db_port")
	Config.DBConfig.DbName = v1.GetString("database.db_name")
	Config.DBConfig.DbUsername = v1.GetString("database.db_username")
	Config.DBConfig.DbPassword = v1.GetString("database.db_password")
	Config.DBConfig.DSN =
		"sqlserver://" + Config.DBConfig.DbUsername + ":" +
			Config.DBConfig.DbPassword + "@" + Config.DBConfig.DbHost +
			":" + Config.DBConfig.DbPort + "?database=" + Config.DBConfig.DbName
	Config.DBConfig.OmittedColumns = v1.GetStringSlice("database.omitted_columns")

	//配置里的密钥是string类型，jwt要求为[]byte类型，必须转换后才能使用
	Config.JWTConfig.SecretKey = []byte(v1.GetString("jwt.secret_key"))
	Config.JWTConfig.ValidityPeriod = v1.GetInt("jwt.validity_period")

	Config.LogConfig.FileName = v1.GetString("log.log_path") + "/status.log"
	Config.LogConfig.MaxSizeForLog = v1.GetInt("log.log_max_size")
	Config.LogConfig.MaxBackup = v1.GetInt("log.log_max_backup")
	Config.LogConfig.MaxAge = v1.GetInt("log.log_max_age")
	Config.LogConfig.Compress = v1.GetBool("log.log_compress")

	Config.UploadConfig.FullPath = v1.GetString("upload_files.full_path") + "/"
	Config.UploadConfig.MaxSizeForUpload = v1.GetInt64("upload_files.max_size") << 20

	Config.EmailConfig.OutgoingMailServer = v1.GetString("email.outgoing_mail_server")
	Config.EmailConfig.Port = v1.GetInt("email.port")
	Config.EmailConfig.Account = v1.GetString("email.account")
	Config.EmailConfig.Password = v1.GetString("email.password")

	Config.PagingConfig.DefaultPageSize = v1.GetInt("paging.default_page_size")
	Config.PagingConfig.MaxPageSize = v1.GetInt("paging.max_page_size")
}
