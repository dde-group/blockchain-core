package logutils

/**
 * @Author: lee
 * @Description:
 * @File: zap
 * @Date: 2021/9/13 6:04 下午
 */

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gitlab.xbit.trade/blockchain/blockchain-core/utils/fileutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"time"
)

type ZapConfig struct {
	Directory    string `json:"directory"     yaml:"directory"   mapstructure:"directory"`
	ShowLine     bool   `json:"show-line"     yaml:"show-line"   mapstructure:"show-line"`
	ZapLevel     string `json:"zap-level"     yaml:"zap-level"   mapstructure:"zap-level"`
	Archive      string `json:"archive"     yaml:"archive"       mapstructure:"archive"`
	Format       string `json:"format"     yaml:"format"         mapstructure:"format"`
	LinkName     string `json:"link-name"     yaml:"link-name"   mapstructure:"link-name"`
	LogInConsole bool   `json:"log-in-console"     yaml:"log-in-console"    mapstructure:"log-in-console"`
	EncodeLevel  string `json:"encode-level"     yaml:"encode-level"        mapstructure:"encode-level"`
	WarnName     string `json:"warn-name"     yaml:"warn-name"   mapstructure:"warn-name"`
	ErrorName    string `json:"error-name"     yaml:"error-name"   mapstructure:"error-name"`
}

var DefaultZapConfig = ZapConfig{
	Directory:    "logs",
	ZapLevel:     "info",
	Archive:      "log",
	Format:       "console",
	LinkName:     "logs/latest_log.log",
	LogInConsole: true,
	EncodeLevel:  "LowercaseColorLevelEncoder",
	ShowLine:     true,
}

const LogCtxKey = "log_ctx_module"

type ZapLogModule struct {
	logger *zap.Logger
	config *ZapConfig
	ctx    context.Context
}

var _ ILogger = (*ZapLogModule)(nil)

func (m *ZapLogModule) WithContext(ctx context.Context) ILogger {
	if nil == ctx {
		return m
	}
	ret := &ZapLogModule{
		logger: m.logger,
		config: m.config,
		ctx:    ctx,
	}
	return ret
}
func (m *ZapLogModule) Info(msg string, fields ...zap.Field) {
	if nil != m.ctx {
		iCtxFields := m.ctx.Value(LogCtxKey)
		ctxFields, ok := iCtxFields.([]zap.Field)
		if ok && nil != ctxFields && len(ctxFields) > 0 {
			fields = append(fields, ctxFields...)
		}
	}
	m.logger.Info(msg, fields...)
}

func (m *ZapLogModule) Infof(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Error(msg string, fields ...zap.Field) {
	if nil != m.ctx {
		iCtxFields := m.ctx.Value(LogCtxKey)
		ctxFields, ok := iCtxFields.([]zap.Field)
		if ok && nil != ctxFields && len(ctxFields) > 0 {
			fields = append(fields, ctxFields...)
		}
	}
	m.logger.Error(msg, fields...)
}

func (m *ZapLogModule) Errorf(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Warn(msg string, fields ...zap.Field) {
	if nil != m.ctx {
		iCtxFields := m.ctx.Value(LogCtxKey)
		ctxFields, ok := iCtxFields.([]zap.Field)
		if ok && nil != ctxFields && len(ctxFields) > 0 {
			fields = append(fields, ctxFields...)
		}
	}
	m.logger.Warn(msg, fields...)
}

func (m *ZapLogModule) Warnf(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Debug(msg string, fields ...zap.Field) {
	if nil != m.ctx {
		iCtxFields := m.ctx.Value(LogCtxKey)
		ctxFields, ok := iCtxFields.([]zap.Field)
		if ok && nil != ctxFields && len(ctxFields) > 0 {
			fields = append(fields, ctxFields...)
		}
	}
	m.logger.Debug(msg, fields...)
}

func (m *ZapLogModule) Debugf(format string, vals ...interface{}) {

}

func (m *ZapLogModule) Fatal(msg string, fields ...zap.Field) {
	if nil != m.ctx {
		iCtxFields := m.ctx.Value(LogCtxKey)
		ctxFields, ok := iCtxFields.([]zap.Field)
		if ok && nil != ctxFields && len(ctxFields) > 0 {
			fields = append(fields, ctxFields...)
		}
	}
	m.logger.Fatal(msg, fields...)
}

func (m *ZapLogModule) Fatalf(format string, vals ...interface{}) {

}
func (m *ZapLogModule) DPanic(msg string, fields ...zap.Field) {
	if nil != m.ctx {
		iCtxFields := m.ctx.Value(LogCtxKey)
		ctxFields, ok := iCtxFields.([]zap.Field)
		if ok && nil != ctxFields && len(ctxFields) > 0 {
			fields = append(fields, ctxFields...)
		}
	}
	m.logger.DPanic(msg, fields...)
}

func (m *ZapLogModule) Panic(msg string, fields ...zap.Field) {
	if nil != m.ctx {
		iCtxFields := m.ctx.Value(LogCtxKey)
		ctxFields, ok := iCtxFields.([]zap.Field)
		if ok && nil != ctxFields && len(ctxFields) > 0 {
			fields = append(fields, ctxFields...)
		}
	}
	m.logger.Panic(msg, fields...)
}

func NewZapLogModule(config ZapConfig) (ILogger, error) {
	logger, err := newZapLogger(config)
	if nil != err {
		return nil, err
	}
	ret := ZapLogModule{
		logger: logger,
		config: &config,
	}

	return &ret, nil
}

func newZapLogger(config ZapConfig) (logger *zap.Logger, err error) {
	var zapConfig = config
	if err = fileutils.CreateDirectoryIfNotExist(zapConfig.Directory, os.ModePerm); nil != err {
		return nil, err
	}

	var level = zap.InfoLevel
	// 初始化配置文件的Level
	switch zapConfig.ZapLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dPanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	panicLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.PanicLevel
	})

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(&config, level), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore(&config, level), zap.AddStacktrace(panicLevel))
	}
	if zapConfig.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
		logger = logger.WithOptions(zap.AddCallerSkip(2))
	}

	return logger, nil
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(cfg *ZapConfig) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch {
	case cfg.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case cfg.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case cfg.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case cfg.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder(cfg *ZapConfig) zapcore.Encoder {
	if "json" == cfg.Format {
		return zapcore.NewJSONEncoder(getEncoderConfig(cfg))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(cfg))
}
func getEncoderCore(config *ZapConfig, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer, err := GetWriteSyncer(config) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}

	var cores []zapcore.Core
	cores = append(cores, zapcore.NewCore(getEncoder(config), writer, level))
	if config.WarnName != "" {
		warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == zapcore.WarnLevel
		})
		warnWriter, err := GetWarnWriteSyncer(config.WarnName, config)
		if nil != err {
			return
		}
		cores = append(cores, zapcore.NewCore(getEncoder(config), warnWriter, warnLevel))
	}

	if "" != config.ErrorName {
		errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		errWriter, err := GetErrorWriteSyncer(config.ErrorName, config)
		if nil != err {
			return
		}
		cores = append(cores, zapcore.NewCore(getEncoder(config), errWriter, errorLevel))
	}

	return zapcore.NewTee(cores...)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// zap logger中加入file-rotatelogs
func GetWriteSyncer(cfg *ZapConfig) (zapcore.WriteSyncer, error) {
	var filePath string
	filePath = path.Join(cfg.Directory, cfg.Archive+"-%Y-%m-%d.log")

	var linkName rotatelogs.Option
	if cfg.Archive == "" {
		linkName = rotatelogs.WithLinkName(cfg.LinkName)
	} else {
		linkName = rotatelogs.WithLinkName(cfg.LinkName + "_" + cfg.Archive)
	}

	fileWriter, err := rotatelogs.New(
		filePath,
		linkName,
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if cfg.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

func GetWarnWriteSyncer(name string, cfg *ZapConfig) (zapcore.WriteSyncer, error) {
	var filePath string
	filePath = path.Join(cfg.Directory, name+"-%Y-%m-%d.log")
	linkName := rotatelogs.WithLinkName(path.Join(cfg.Directory, name+".log"))
	fileWriter, err := rotatelogs.New(
		filePath,
		linkName,
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	return zapcore.AddSync(fileWriter), err
}

func GetErrorWriteSyncer(name string, cfg *ZapConfig) (zapcore.WriteSyncer, error) {
	var filePath string
	filePath = path.Join(cfg.Directory, name+"-%Y-%m-%d.log")
	linkName := rotatelogs.WithLinkName(path.Join(cfg.Directory, name+".log"))
	fileWriter, err := rotatelogs.New(
		filePath,
		linkName,
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	return zapcore.AddSync(fileWriter), err
}
