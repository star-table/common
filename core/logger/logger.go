package logger

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/getsentry/sentry-go"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/sentry/client"
	"github.com/star-table/common/core/util/strs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"sync"
)

// 日志包的封装

// 常规用法，在一个服务的 main 文件中，调用 `var log = logger.GetDefaultLogger()`
// 在后续的业务代码中，可以直接调用 common 包下 logger 包的 Error()、 Info() 等方法。
// 例如：在 xx_service.go 业务代码文件中，声明一个 全局变量 log：`var log = *logger.GetDefaultLogger()`
// 而后，在对应的方法中调用日志方法：`log.Error(someErr.Error())`
var (
	_lock   sync.Mutex
	_logger = map[string]*SysLogger{}
	// 额外的日志配置
	extraLoggerOption = map[string]interface{}{}
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

type LogKafka struct {
	Producer sarama.SyncProducer
	Topic    string
}

type SysLogger struct {
	log          *zap.Logger
	name         string
	path         string
	level        zapcore.Level
	isConsoleOut bool
	kafkaTopic   string
}

func (lk *LogKafka) Write(p []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = lk.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = lk.Producer.SendMessage(msg)
	if err != nil {
		return
	}
	return

}

func getTraceIdFieldByThreadLocal() (string, zap.Field) {
	traceId := threadlocal.GetTraceId()
	return traceId, zap.String(consts.TraceIdLogKey, traceId)
}

func getTraceIdFieldByHttpContext(httpContext *model.HttpContext) (string, zap.Field) {
	traceId := httpContext.TraceId
	return traceId, zap.String(consts.TraceIdLogKey, traceId)
}

func (s *SysLogger) Info(msg interface{}, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByThreadLocal()
	str, ok := msg.(string)
	if !ok {
		str = strs.ObjectToString(msg)
	}
	s.log.Info("[traceId="+traceId+"]"+str, append(fields, fie)...)
}
func (s *SysLogger) InfoH(httpContext *model.HttpContext, msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Info("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) Infof(fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	traceId, fie := getTraceIdFieldByThreadLocal()
	s.log.Info("[traceId="+traceId+"]"+msg, fie)
}

func (s *SysLogger) InfoHf(httpContext *model.HttpContext, fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)

	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Info("[traceId="+traceId+"]"+msg, fie)
}

func (s *SysLogger) Error(msg interface{}, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByThreadLocal()
	str, ok := msg.(string)
	if !ok {
		str = strs.ObjectToString(msg)
	}
	s.log.Error("[traceId="+traceId+"]"+str, append(fields, fie)...)
}

func (s *SysLogger) ErrorH(httpContext *model.HttpContext, msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Error("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) Errorf(fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	traceId, fie := getTraceIdFieldByThreadLocal()
	s.log.Error("[traceId="+traceId+"]"+msg, fie)
}

func (s *SysLogger) ErrorHf(httpContext *model.HttpContext, fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)

	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Error("[traceId="+traceId+"]"+msg, fie)
}

func (s *SysLogger) Fatal(msg interface{}, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByThreadLocal()
	str, ok := msg.(string)
	if !ok {
		str = strs.ObjectToString(msg)
	}
	s.log.Fatal("[traceId="+traceId+"]"+str, append(fields, fie)...)
}

func (s *SysLogger) Debug(msg interface{}, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByThreadLocal()
	str, ok := msg.(string)
	if !ok {
		str = strs.ObjectToString(msg)
	}
	s.log.Debug("[traceId="+traceId+"]"+str, append(fields, fie)...)
}

func (s *SysLogger) DebugH(httpContext *model.HttpContext, msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Debug("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) Debugf(fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	traceId, fie := getTraceIdFieldByThreadLocal()
	s.log.Debug("[traceId="+traceId+"]"+msg, fie)
}

func (s *SysLogger) DebugHf(httpContext *model.HttpContext, fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Debug("[traceId="+traceId+"]"+msg, fie)
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// 获得默认Logger对象
func GetDefaultLogger() *SysLogger {
	return CreateLogger("default")
}

func GetMQLogger() *SysLogger {
	return CreateLogger("mq")
}

func GetRequestResponseLogger() *SysLogger {
	return CreateLogger("rr")
}

func CreateLogger(name string) *SysLogger {
	return &SysLogger{
		name: name,
	}
}

// 获得Logger对象
func (s *SysLogger) Init() {
	if s.log == nil {
		_lock.Lock()
		defer _lock.Unlock()
		if s.log == nil {
			s.InitLogger()
		}
	}
}

// 设置额外的配置，并在实例化 logger 时使用。但必须在实例化之前调用 SetExtraLoggerOption
func (s *SysLogger) SetExtraLoggerOption(key string, value interface{}) {
	extraLoggerOption[key] = value
}

func (s *SysLogger) InitLogger() *SysLogger {
	logConfig := config.GetLogConfig(s.name)
	if logConfig == nil {
		panic(s.name + " log config not exist!")
	}

	level := getLoggerLevel(logConfig.Level)
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.LogPath,                // ⽇志⽂件路径
		MaxSize:    s.getFileSizeByConfig(logConfig), // megabytes
		MaxBackups: s.getFileNumByConfig(logConfig),  //最多保留20个备份
		LocalTime:  true,
		Compress:   true, // 是否压缩 disabled by default
	})
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder.EncodeCaller = zapcore.FullCallerEncoder

	var core zapcore.Core

	if logConfig.IsConsoleOut == true {
		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		consoleOut := zapcore.Lock(os.Stdout)
		syncWriterList := []zapcore.WriteSyncer{
			consoleOut, syncWriter,
		}

		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoder),
			zapcore.NewMultiWriteSyncer(syncWriterList...),
			zap.NewAtomicLevelAt(level))
	}

	// 设置初始化字段
	filed := zap.Fields(zap.String(consts.LogTagKey, logConfig.Tag), zap.String(consts.LogAppKey, config.GetApplication().Name))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), filed)
	if s, ok := extraLoggerOption["sentryClient"]; ok {
		sentryClient, ok2 := s.(*sentry.Client)
		if ok2 {
			keyword := ""
			if text, ok3 := extraLoggerOption["dingTalkLogGroupKeyword"]; ok3 {
				textStr, ok4 := text.(string)
				if !ok4 {
					keyword = "【北极星警告】【告警日志】"
				} else {
					keyword = textStr
				}
			} else {
				keyword = "【北极星警告】【告警日志】"
			}
			sentryCfg := client.SentryCoreConfig{
				Level: zap.ErrorLevel,
				Tags: map[string]string{
					"source": "runx",
				},
				ExtraStringOption: map[string]string{
					"dingTalkLogGroupKeyword": keyword,
				},
			}
			sentryCore := client.NewSentryCore(sentryCfg, sentryClient)
			logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewTee(core, sentryCore)
			}))
		}
	}

	//s.log = logger.Sugar()
	s.log = logger
	s.path = logConfig.LogPath
	s.level = level
	s.isConsoleOut = logConfig.IsConsoleOut

	_logger[logConfig.LogPath] = s

	return _logger[logConfig.LogPath]
}

func (s *SysLogger) getFileSizeByConfig(logConfig *config.LogConfig) int {
	if logConfig.FileSize < 1 {
		return 1024
	}

	i, _ := strconv.Atoi(strconv.FormatInt(logConfig.FileSize, 10))
	return i
}

func (s *SysLogger) getFileNumByConfig(logConfig *config.LogConfig) int {
	if logConfig.FileNum < 1 {
		return 20
	}
	return logConfig.FileNum
}
