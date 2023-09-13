package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"pmis-backend-go/global"
)

//global.Logger和global.SugaredLogger已经在全局变量中定义，这里是修改他们的配置

// 自定义 生成编码器的函数
func newEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	//使用json格式来记录
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 自定义 写入同步器的函数
func newWriteSyncer(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	loggerRule := &lumberjack.Logger{
		Filename:   filename,  //日志文件的位置
		MaxSize:    maxSize,   //在进行切割之前，日志文件的最大大小（MB）
		MaxBackups: maxBackup, //保留旧文件的最大个数
		MaxAge:     maxAge,    //保留旧文件的最大天数
		Compress:   compress,  //是否压缩旧文件
	}
	return zapcore.AddSync(loggerRule)
}

func InitLogger() {
	encoder := newEncoder() //调用自定义的编码器函数，生成新的编码器
	//调用自定义的写入同步器函数，传入文件路径+名称、最大尺寸、最大备份数量、最大保存天数，生成新的写入同步器
	writeSyncer := newWriteSyncer(
		global.Config.Log.FileName,
		global.Config.Log.MaxSizeForLog,
		global.Config.Log.MaxBackup,
		global.Config.Log.MaxAge,
		global.Config.Log.Compress,
	)
	//声明zap的核心参数
	var core zapcore.Core
	mode := global.Config.App.Mode
	//如果是开发模式：
	if mode == "debug" {
		//生成开发模式下的、encoder默认配置文件，用于管理控制台的显示内容
		developmentEncoderConfig := zap.NewDevelopmentEncoderConfig()
		//调整格式
		developmentEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		//根据encoder配置文件，生成控制台的encoder
		consoleEncoder := zapcore.NewConsoleEncoder(developmentEncoderConfig)
		//生成zap的核心文件core
		core = zapcore.NewTee(
			// 往日志文件里面写
			zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel),
			// 在终端输出
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		//如果是生产环境就只写到日志里，不在终端输出
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	}
	global.Logger = zap.New(core, zap.AddCaller()) //根据zap的要求，生成一个日志记录器
	global.SugaredLogger = global.Logger.Sugar()   //使用加糖模式的日志记录器，牺牲点效率，但简单一些

	defer global.SugaredLogger.Sync()
}
