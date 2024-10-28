package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"ascale/pkg/conf/env"
	"ascale/pkg/stat/prom"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errProm = prom.BusinessErrCount
var fm sync.Map

const (
	_timeFormat = "2006-01-02T15:04:05.999999"

	// log level defined in level.go.
	_levelValue = "level_value"
	//  log level name: INFO, WARN...
	_level = "level"
	// log time.
	_time = "time"
	// request path.
	// _title = "title"
	// log file.
	_source = "source"
	// common log filed.
	_log = "log"
	// app name.
	_appID = "app_id"
	// container ID.
	_instanceID = "instance_id"
	// uniq ID from trace.
	_tid = "traceid"
	// request time.
	// _ts = "ts"
	// requester.
	_caller = "caller"
	// container environment: prod, pre, uat, fat.
	_deplyEnv = "env"
	// container area.
	_zone = "zone"
	// mirror flag
	_mirror = "mirror"
	// color.
	_color = "color"
	// cluster.
	_cluster = "cluster"
)

type Config struct {
	Family string
	Host   string

	// stdout
	Stdout bool

	Filter []string
}

var (
	l *zap.Logger
)

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func init() {
	host, _ := os.Hostname()
	l, _ = NewDevelopment(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Fields(zap.String(_appID, env.AppID)),
		zap.Fields(zap.String(_instanceID, host)),
	)
}

func ZapLogger() *zap.Logger {
	return l
}

func Init(conf *Config) {
	logger, err := NewProductionWithCore(WrapCore(
		ReportAllErrors(true),
		ServiceName(fmt.Sprintf("%s-%s", env.AppID, env.DeployEnv)),
	),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(fmt.Sprintf("initial log failed: %v", err))
	}

	l = logger

}

func Debug(msg string) {
	fields := appendFields()
	l.Debug(msg, fields...)
}

func Debugf(msg string, a ...interface{}) {
	fields := appendFields()
	l.Debug(fmt.Sprintf(msg, a...), fields...)
}

func DebugWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.Debug(msg, fields...)
}

// Info logs a message at InfoLevel
func Info(msg string) {
	fields := appendFields()
	l.Info(msg, fields...)
}

func Infof(msg string, a ...interface{}) {
	fields := appendFields()
	l.Info(fmt.Sprintf(msg, a...), fields...)
}

func InfoWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.Info(msg, fields...)
}

func Warn(msg string) {
	fields := appendFields()
	l.Warn(msg, fields...)
}

func Warnf(msg string, a ...interface{}) {
	fields := appendFields()
	l.Warn(fmt.Sprintf(msg, a...), fields...)
}

func WarnWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.Warn(msg, fields...)
}

func Error(msg string) {
	fields := appendFields()
	l.Error(msg, fields...)
}

func Errorf(msg string, a ...interface{}) {
	fields := appendFields()
	l.Error(fmt.Sprintf(msg, a...), fields...)
}

func ErrorWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.Error(msg, fields...)
}

func DPanic(msg string) {
	fields := appendFields()
	l.DPanic(msg, fields...)
}

func DPanicf(msg string, a ...interface{}) {
	fields := appendFields()
	l.DPanic(fmt.Sprintf(msg, a...), fields...)
}

func DPanicWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.DPanic(msg, fields...)
}

func Panic(msg string) {
	fields := appendFields()
	l.Panic(msg, fields...)
}

func Panicf(msg string, a ...interface{}) {
	fields := appendFields()
	l.Panic(fmt.Sprintf(msg, a...), fields...)
}

func PanicWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.Panic(msg, fields...)
}

func Fatal(msg string) {
	fields := appendFields()
	l.Fatal(msg, fields...)
}

func Fatalf(msg string, a ...interface{}) {
	fields := appendFields()
	l.Fatal(fmt.Sprintf(msg, a...), fields...)
}

func FatalWithFields(msg string, fields ...zapcore.Field) {
	fields = append(fields, appendFields()...)
	l.Fatal(msg, fields...)
}

func Close() {
	l.Sync()
}

func appendFields() []zapcore.Field {
	fields := make([]zapcore.Field, 0)
	fields = append(fields, zap.String(_appID, env.AppID))

	return fields
}
