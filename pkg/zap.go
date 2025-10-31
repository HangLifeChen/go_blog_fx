package core

import (
	"log"
	"os"

	"go_blog/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化并返回一个基于配置设置的新 zap.Logger 实例
func InitLogger(config config.Config) *zap.Logger {
	zapCfg := config.ThirdPartyServices.Zap

	//1. 创建一个用于日志输出的 writeSyncer
	writeSyncer := getLogWriter(zapCfg.Filename, zapCfg.MaxSize, zapCfg.MaxBackups, zapCfg.MaxAge)

	//1.1 如果配置了控制台输出，则添加控制台输出
	if zapCfg.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}

	//2. 创建日志格式化的编码器
	encoder := getEncoder()

	//3. 根据配置确定日志级别
	var logLevel zapcore.Level
	//3.1 解析yaml文件中的日志级别
	if err := logLevel.UnmarshalText([]byte(zapCfg.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	//4. 创建核心和日志实例
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}

// getLogWriter 返回一个 zapcore.WriteSyncer，该写入器利用 lumberjack 包，实现日志的滚动记录
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,   // 日志文件的位置
		MaxSize:    maxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: maxBackups, // 保留旧文件的最大个数
		MaxAge:     maxAge,     // 保留旧文件的最大天数
	}
	return zapcore.AddSync(lumberJackLogger)
}

// getEncoder 返回一个为生产日志配置的 JSON 编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         // 使用 ISO8601 时间格式
	encoderConfig.TimeKey = "time"                                // 将时间键设置为 "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       // 将日志级别大写
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder // 将持续时间编码为秒
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       // 将调用者编码为短格式
	return zapcore.NewJSONEncoder(encoderConfig)
}
