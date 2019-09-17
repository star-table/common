package logger

import (
	"fmt"
	"gitea.bjx.cloud/allstar/common/core/config"
	"gitea.bjx.cloud/allstar/common/core/model"
	"gitea.bjx.cloud/allstar/common/core/threadlocal"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strconv"
	"strings"
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
	log          *zap.SugaredLogger
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

//func getTraceId() string {
//	if traceId, ok := threadlocal.Mgr.GetValue(consts.TraceIdKey); ok {
//		if traceId != nil {
//			fmt.Println(traceId)
//			return traceId.(string)
//		}
//	}
//	return ""
//}

func insertTraceId(args ...interface{}) []interface{} {
	cl := []interface{}{fmt.Sprintf("traceId=%s;", threadlocal.GetTraceId())}
	return append(cl, args...)
}
func insertTraceIdf(template string, args ...interface{}) (string, []interface{}) {
	t := fmt.Sprintf("%s;%s", "%s", template)
	return t, insertTraceId(args...)
}
func insertTraceIdH(httpContext *model.HttpContext, args ...interface{}) []interface{} {
	cl := []interface{}{fmt.Sprintf("traceId=%s;", httpContext.TraceId)}
	return append(cl, args...)
}
func insertTraceIdHf(httpContext *model.HttpContext, template string, args ...interface{}) (string, []interface{}) {
	t := fmt.Sprintf("%s;%s", "%s", template)
	return t, insertTraceIdH(httpContext, args...)
}

func (s *SysLogger) Info(args ...interface{}) {
	s.Init()
	s.log.Info(insertTraceId(args...))
}
func (s *SysLogger) InfoH(httpContext *model.HttpContext, args ...interface{}) {
	s.Init()
	s.log.Info(insertTraceIdH(httpContext, args...))
}

func (s *SysLogger) Infof(template string, args ...interface{}) {
	s.Init()
	t, ags := insertTraceIdf(template, args...)
	s.log.Infof(t, ags...)
}

func (s *SysLogger) InfoHf(httpContext *model.HttpContext, template string, args ...interface{}) {
	s.Init()
	t, ags := insertTraceIdHf(httpContext, template, args...)
	s.log.Infof(t, ags...)
}

func (s *SysLogger) Error(args ...interface{}) {
	s.Init()
	s.log.Error(insertTraceId(args...))
}

func (s *SysLogger) ErrorH(httpContext *model.HttpContext, args ...interface{}) {
	s.Init()
	s.log.Error(insertTraceIdH(httpContext, args...))
}

func (s *SysLogger) Errorf(template string, args ...interface{}) {
	s.Init()
	t, ags := insertTraceIdf(template, args...)
	s.log.Errorf(t, ags...)
}

func (s *SysLogger) ErrorHf(httpContext *model.HttpContext, template string, args ...interface{}) {
	s.Init()
	t, ags := insertTraceIdHf(httpContext, template, args...)
	s.log.Errorf(t, ags...)
}

func (s *SysLogger) Debug(args ...interface{}) {
	s.Init()
	s.log.Debug(insertTraceId(args...))
}

func (s *SysLogger) DebugH(httpContext *model.HttpContext, args ...interface{}) {
	s.Init()
	s.log.Debug(insertTraceIdH(httpContext, args...))
}

func (s *SysLogger) Debugf(template string, args ...interface{}) {
	s.Init()
	t, ags := insertTraceIdf(template, args...)
	s.log.Debugf(t, ags...)
}

func (s *SysLogger) DebugHf(httpContext *model.HttpContext, template string, args ...interface{}) {
	s.Init()
	t, ags := insertTraceIdHf(httpContext, template, args...)
	s.log.Debugf(t, ags...)
}

//func (s *SysLogger) Fatal(args ...interface{}) {
//	s.log.Fatal(insertTraceId(args))
//}
//
//func (s *SysLogger) FatalH(httpContext *domains.HttpContext, args ...interface{}) {
//	s.log.Fatal(insertTraceIdH(httpContext, args))
//}
//
//func (s *SysLogger) Fatalf(template string, args ...interface{}) {
//	t, ags := insertTraceIdf(template, args)
//	s.log.Fatalf(t, ags)
//}
//
//func (s *SysLogger) FatalHf(httpContext *domains.HttpContext, template string, args ...interface{}) {
//	t, ags := insertTraceIdHf(httpContext, template, args)
//	s.log.Fatalf(t, ags)
//}

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

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	//consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	var allCore []zapcore.Core

	if logConfig.IsConsoleOut == true {
		allCore = append(allCore, zapcore.NewCore(zapcore.NewJSONEncoder(encoder), consoleDebugging, zap.NewAtomicLevelAt(level)))
	}

	if logConfig.EnableKafka == true {
		s.kafkaTopic = strings.Replace(logConfig.KafkaTopic, ":appName:", config.GetApplication().Name, -1)

		//设置日志输入到Kafka的配置
		kafkaConfig := sarama.NewConfig()
		//等待leader服务器保存成功后的响应
		kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
		// 随机的分区类型
		kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
		//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
		kafkaConfig.Producer.Return.Successes = true
		kafkaConfig.Producer.Return.Errors = true

		var err error = nil

		kl := LogKafka{
			Topic: s.kafkaTopic,
		}
		kl.Producer, err = sarama.NewSyncProducer(strings.Split(logConfig.KafkaNameServers, ","), kafkaConfig)
		if err != nil {
			fmt.Println(" log config kafka init error. ")
		}
		topicErrors := zapcore.AddSync(&kl)
		// 打印在kafka
		kafkaEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
		kafkaCore := zapcore.NewCore(kafkaEncoder, topicErrors, zap.NewAtomicLevelAt(level))
		allCore = append(allCore, kafkaCore)
	} else {
		s.kafkaTopic = ""
	}

	allCore = append(allCore, zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(level)))

	core := zapcore.NewTee(allCore...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	//level := getLoggerLevel(logConfig.Level)
	//syncWriter := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   logConfig.LogPath, // ⽇志⽂件路径
	//	MaxSize:    1024,     // megabytes
	//	MaxBackups: 20,       //最多保留20个备份
	//	LocalTime:  true,
	//	Compress:   true, // 是否压缩 disabled by default
	//})
	//encoder := zap.NewProductionEncoderConfig()
	//encoder.EncodeTime = zapcore.ISO8601TimeEncoder
	//core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(level))
	//logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	//
	s.log = logger.Sugar()
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
