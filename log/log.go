package log

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultLogger,stdout
func DefaultLogger() *zap.Logger {
	var coreArr []zapcore.Core

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05") //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                  //按级别显示不同颜色
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder                       //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)                           //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), lowPriority)   //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	errorCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), highPriority) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	coreArr = append(coreArr, infoCore)
	coreArr = append(coreArr, errorCore)
	return setLogLevel(zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()))
}

func setLogLevel(log *zap.Logger) *zap.Logger {
	switch viper.GetString("log.level") {
	case "DEBUG":
		return log.WithOptions(zap.IncreaseLevel(zapcore.DebugLevel))
	case "INFO":
		return log.WithOptions(zap.IncreaseLevel(zapcore.InfoLevel))
	case "ERROR":
		return log.WithOptions(zap.IncreaseLevel(zapcore.WarnLevel))
	case "DPANIC":
		return log.WithOptions(zap.IncreaseLevel(zapcore.DPanicLevel))
	case "PANIC":
		return log.WithOptions(zap.IncreaseLevel(zapcore.PanicLevel))
	case "FATAL":
		return log.WithOptions(zap.IncreaseLevel(zapcore.FatalLevel))
	default:
		return log.WithOptions(zap.IncreaseLevel(zapcore.ErrorLevel))
	}
}
