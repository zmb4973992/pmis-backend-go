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
	v             = viper.New()
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
	Path          string
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

func Init() {
	v.AddConfigPath("./config/") //告诉viper，配置文件的路径在哪
	v.SetConfigName("config")    //告诉viper，配置文件的前缀是什么
	v.SetConfigType("yaml")      //告诉viper，配置文件的后缀是什么
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("MyConfig file error: %w \n", err))
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已修改:", e.Name)
		loadConfig()
	})
	v.WatchConfig()
	loadConfig()
}

func loadConfig() {
	Config.APPConfig.AppMode = v.GetString("App.AppMode")
	Config.APPConfig.HttpPort = v.GetString("App.HttpPort")

	Config.DBConfig.DbHost = v.GetString("Database.DbHost")
	Config.DBConfig.DbPort = v.GetString("Database.DbPort")
	Config.DBConfig.DbName = v.GetString("Database.DbName")
	Config.DBConfig.DbUsername = v.GetString("Database.DbUsername")
	Config.DBConfig.DbPassword = v.GetString("Database.DbPassword")
	Config.DBConfig.DSN =
		"sqlserver://" + Config.DBConfig.DbUsername + ":" +
			Config.DBConfig.DbPassword + "@" + Config.DBConfig.DbHost +
			":" + Config.DBConfig.DbPort + "?database=" + Config.DBConfig.DbName
	Config.DBConfig.OmittedColumns = v.GetStringSlice("Database.OmittedColumns")

	//配置里的密钥是string类型，jwt要求为[]byte类型，必须转换后才能使用
	Config.JWTConfig.SecretKey = []byte(v.GetString("JWT.SecretKey"))
	Config.JWTConfig.ValidityPeriod = v.GetInt("JWT.ValidityPeriod")

	Config.LogConfig.Path = v.GetString("Log.LogPath")
	Config.LogConfig.FileName = v.GetString("Log.LogPath") + "/status.log"
	Config.LogConfig.MaxSizeForLog = v.GetInt("Log.LogMaxSize")
	Config.LogConfig.MaxBackup = v.GetInt("Log.LogMaxBackup")
	Config.LogConfig.MaxAge = v.GetInt("Log.LogMaxAge")
	Config.LogConfig.Compress = v.GetBool("Log.LogCompress")

	Config.UploadConfig.FullPath = v.GetString("UploadFiles.FullPath") + "/"
	Config.UploadConfig.MaxSizeForUpload = v.GetInt64("UploadFiles.MaxSize") << 20

	Config.EmailConfig.OutgoingMailServer = v.GetString("Email.OutgoingMailServer")
	Config.EmailConfig.Port = v.GetInt("Email.Port")
	Config.EmailConfig.Account = v.GetString("Email.Account")
	Config.EmailConfig.Password = v.GetString("Email.Password")

	Config.PagingConfig.DefaultPageSize = v.GetInt("Paging.DefaultPageSize")
	Config.PagingConfig.MaxPageSize = v.GetInt("Paging.MaxPageSize")

}
