package logging

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var zapLog *zap.Logger

func initZap(config Config) {
	var coreArr []zapcore.Core

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()            //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder    //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder //按级别显示不同颜色
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder       //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//encoder := zapcore.NewJSONEncoder(encoderConfig)

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		if config.Debug {
			return lev < zap.ErrorLevel && lev >= zap.DebugLevel
		} else {
			return lev < zap.ErrorLevel && lev > zap.DebugLevel
		}
	})

	log.Println("输出日志目录-----------------------------------------------------")
	log.Println(config.InfoFile)
	//info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.InfoFile, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    100,             //文件大小限制,单位MB
		MaxBackups: 100,             //最大保留日志文件数量
		MaxAge:     30,              //日志文件保留天数
		Compress:   false,           //是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	//error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.ErrorFile, //日志文件存放目录
		MaxSize:    100,              //文件大小限制,单位MB
		MaxBackups: 5,                //最大保留日志文件数量
		MaxAge:     30,               //日志文件保留天数
		Compress:   false,            //是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	options := []zap.Option{zap.AddCaller()}
	if config.Debug {
		options = append(options, zap.AddStacktrace(lowPriority))
	} else {
		options = append(options, zap.AddStacktrace(zapcore.WarnLevel))
	}
	zapLog = zap.New(zapcore.NewTee(coreArr...), options...) //zap.AddCaller()为显示文件名和行号
}
