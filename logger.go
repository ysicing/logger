// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package logger

import (
	"fmt"
	"github.com/ysicing/go-utils/extime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"runtime"
	"time"
)

const (
	DefaultMaxSize  = 10 // MB
	DefaultBackups  = 3
	DefaultMaxAge   = 7 // days
	DefaultCompress = true
)

var (
	Logger *zap.Logger
)

//func NewLogger(cfg *LogConfig) *zap.Logger {
//	encoder := getEncoder() // 编码器
//	errPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
//		return level >= zap.ErrorLevel
//	})
//	debugPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
//		return level < zap.ErrorLevel && level >= zap.DebugLevel
//	})
//	if cfg.Simple {
//		//writeSyncer := getLogWriter() // 写日志
//		errCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(getLogWriter("err"), zapcore.AddSync(os.Stdout)), errPriority)
//		debugCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(getLogWriter("debug"), zapcore.AddSync(os.Stdout)), debugPriority)
//		return zap.New(zapcore.NewTee(debugCore, errCore), zap.AddCaller())
//	}
//	writeSyncer := getLogWriterSimple()                                                                                        // 写日志
//	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel) // 如何写，写到哪, 什么级别写
//	return zap.New(zapcore.NewTee(core), zap.AddCaller())
//}

func InitLogger(cfg *LogConfig) {
	encoder := getEncoder() // 编码器
	errPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zap.ErrorLevel && level >= zap.DebugLevel
	})
	if cfg.Simple {
		writeSyncer := getLogWriterSimple()                                                                                        // 写日志
		core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel) // 如何写，写到哪, 什么级别写
		Logger = zap.New(zapcore.NewTee(core), zap.AddCaller())
		return
	}
	//writeSyncer := getLogWriter() // 写日志
	errCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(getLogWriter("err"), zapcore.AddSync(os.Stdout)), errPriority)
	debugCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(getLogWriter("debug"), zapcore.AddSync(os.Stdout)), debugPriority)
	Logger = zap.New(zapcore.NewTee(debugCore, errCore), zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder //zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func getLogWriter(loglevel string) zapcore.WriteSyncer {
	var logpath string
	if runtime.GOOS == "linux" {
		logpath = fmt.Sprintf("/var/log/gologger/%v/%v.log", extime.GetToday(), loglevel)
	} else {
		logpath = fmt.Sprintf("/tmp/gologger/%s/%v.log", extime.GetTodayHour(), loglevel)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    DefaultMaxSize,
		MaxBackups: DefaultBackups,
		MaxAge:     DefaultMaxAge,
		Compress:   DefaultCompress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getLogWriterSimple() zapcore.WriteSyncer {
	var logpath string
	if runtime.GOOS == "linux" {
		logpath = fmt.Sprintf("/var/log/gologger/%v/debug.log", extime.GetToday())
	} else {
		logpath = fmt.Sprintf("/tmp/gologger/%s/debug.log", extime.GetTodayHour())
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    DefaultMaxSize,
		MaxBackups: DefaultBackups,
		MaxAge:     DefaultMaxAge,
		Compress:   DefaultCompress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Debug(msg interface{}) {
	Logger.Sugar().Debug(msg)
}

func Info(msg interface{}) {
	Logger.Sugar().Info(msg)
}

func Error(msg interface{}) {
	Logger.Sugar().Error(msg)
}

func Exit(msg interface{}) {
	Logger.Sugar().Error(msg, "\t === err exit === \n")
	os.Exit(-1)
}
