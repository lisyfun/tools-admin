package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"tools-admin/backend/common/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

var atomicLevel = zap.NewAtomicLevel()

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

var pid string

const layout = "2006-01-02 15:04:05.000"
const fileExt = ".log"

func init() {

	logConfig := config.Config.Logger

	pid = strconv.Itoa(os.Getpid())
	baseLogger, _ := zap.NewProduction()

	logger = baseLogger.With(zap.String("pid", pid)).Sugar()

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		EncodeLevel:   customLevelEncoder,
		EncodeTime:    customTimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeName:    zapcore.FullNameEncoder,
	}

	filePath := logConfig.Path
	maxAge := logConfig.MaxAge

	infoWriter := getWriter(filePath+string(os.PathSeparator)+config.Config.Server.Name+fileExt, maxAge)

	level := getLoggerLevel(logConfig.Level)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(infoWriter)),
		level,
	)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func Debug(msg string, tags ...any) {
	logger.Debugf(msg, tags...)
}

func Info(msg string, tags ...any) {
	logger.Infof(msg, tags...)
}

func Warn(msg string, tags ...any) {
	logger.Warnf(msg, tags...)
}

func Error(msg string, tags ...any) {
	logger.Errorf(msg, tags...)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%s", t.Format(layout)))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%s", fmt.Sprintf("|%s| |%s|", pid, level.CapitalString())))
}

func getWriter(filename string, maxAge int64) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每8小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, fileExt, "", -1)+"_%Y%m%d.log", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
