package logger

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/threadlocal"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"sync"
)

var (
	_lock   sync.Mutex
	_logger = map[string]*SysLogger{}
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

func (s *SysLogger) Info(msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByThreadLocal()
	s.log.Info("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}
func (s *SysLogger) InfoH(httpContext *model.HttpContext, msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Info("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) Infof(fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	s.Info(msg)
}

func (s *SysLogger) InfoHf(httpContext *model.HttpContext, fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	s.InfoH(httpContext, msg)
}

func (s *SysLogger) Error(msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByThreadLocal()
	s.log.Error("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) ErrorH(httpContext *model.HttpContext, msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Error("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) Errorf(fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	s.Error(msg)
}

func (s *SysLogger) ErrorHf(httpContext *model.HttpContext, fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	s.ErrorH(httpContext, msg)
}

func (s *SysLogger) Debug(msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByThreadLocal()
	s.log.Debug("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) DebugH(httpContext *model.HttpContext, msg string, fields ...zap.Field) {
	s.Init()
	traceId, fie := getTraceIdFieldByHttpContext(httpContext)
	s.log.Debug("[traceId="+traceId+"]"+msg, append(fields, fie)...)
}

func (s *SysLogger) Debugf(fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	s.Debug(msg)
}

func (s *SysLogger) DebugHf(httpContext *model.HttpContext, fmtstr string, args ...interface{}) {
	s.Init()
	msg := fmt.Sprintf(fmtstr, args...)
	s.DebugH(httpContext, msg)
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

		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoder),
			zapcore.NewMultiWriteSyncer(consoleOut, syncWriter),
			zap.NewAtomicLevelAt(level))
	}

	// 设置初始化字段
	filed := zap.Fields(zap.String(consts.LogTagKey, logConfig.Tag), zap.String(consts.LogAppKey, config.GetApplication().Name))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), filed)

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
