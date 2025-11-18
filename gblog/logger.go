package main

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// 初始化日志配置
func InitLogger(env string) {
	var core zapcore.Core

	// 日志输出格式：JSON（生产）或控制台（开发）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // 级别大写（INFO, ERROR）
		EncodeTime:     customTimeEncoder,           // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 调用者信息（文件名:行号）
	}

	// 根据环境选择编码器
	var encoder zapcore.Encoder
	if env == "dev" {
		// 开发环境：控制台输出，更友好的格式
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// 生产环境：JSON格式，便于日志收集工具解析
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 日志级别：debug（开发）/ info（生产）
	level := zap.InfoLevel
	if env == "dev" {
		level = zap.DebugLevel
	}

	// 输出目标：控制台 + 文件（按大小/时间切割）
	writeSyncer := getLogWriter()
	core = zapcore.NewCore(encoder, writeSyncer, level)

	// 开发环境额外开启调用者信息和堆栈跟踪
	if env == "dev" {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		logger = zap.New(core, zap.AddCaller())
	}

	zap.ReplaceGlobals(logger) // 替换zap全局日志实例
}

// 自定义时间格式（毫秒级）
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 日志文件输出配置（自动切割、压缩、清理）
func getLogWriter() zapcore.WriteSyncer {
	// 使用lumberjack实现日志轮转
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/gblog.log", // 日志文件路径
		MaxSize:    10,                 // 单个文件最大10MB
		MaxBackups: 30,                 // 最多保留30个备份文件
		MaxAge:     7,                  // 保留7天
		Compress:   true,               // 压缩旧日志
	}

	// 同时输出到控制台和文件
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),        // 控制台
		zapcore.AddSync(lumberJackLogger), // 文件
	)
}
