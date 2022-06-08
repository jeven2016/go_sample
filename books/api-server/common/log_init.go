package common

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

//https://blog.csdn.net/test1280/article/details/117266333
func InitLog(cfg Config) *zap.Logger {
	atomicLevel := setLogLevel(cfg.LogLevel)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 日志轮转
	writer := &lumberjack.Logger{
		// 日志名称
		Filename: "api-server.log",
		// 日志大小限制，单位MB
		MaxSize: 100,
		// 历史日志文件保留天数
		MaxAge: 30,
		// 最大保留历史日志数量
		MaxBackups: 10,
		// 本地时区
		LocalTime: true,
		// 历史日志文件压缩
		Compress: true,
	}

	syncers := []zapcore.WriteSyncer{zapcore.AddSync(writer)}
	if cfg.OutputConsole {
		//同时输出到控制台
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(syncers...),
		atomicLevel,
	)

	// 开启开发模式，堆栈跟踪
	options := []zap.Option{zap.AddCaller()}

	// 开启文件及行号
	if cfg.Dev {
		options = append(options, zap.Development())
	}
	// 设置默认携带的字段
	options = append(options, zap.Fields(zap.String("serviceName", cfg.ServiceName)))
	logger := zap.New(zapCore, options...)
	return logger
}

func setLogLevel(logLevel string) zap.AtomicLevel {
	atomicLevel := zap.NewAtomicLevel()
	switch logLevel {
	case "DEBUG":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "INFO":
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case "WARN":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "ERROR":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case "DPANIC":
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	case "PANIC":
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case "FATAL":
		atomicLevel.SetLevel(zapcore.FatalLevel)
	}

	return atomicLevel
}
