package emt

import mqtt "github.com/eclipse/paho.mqtt.golang"

type PahoErrorLogger struct{}
type PahoDebugLogger struct{}

func (PahoErrorLogger) Println(v ...interface{})               { log.Error(v) }
func (PahoErrorLogger) Printf(format string, v ...interface{}) { log.Errorf(format, v) }
func (PahoDebugLogger) Println(v ...interface{})               { log.Debug(v) }
func (PahoDebugLogger) Printf(format string, v ...interface{}) { log.Debugf(format, v) }

// 打开MQTT底层库的日志输出
func init() {
	mqtt.ERROR = PahoErrorLogger{}
	mqtt.CRITICAL = PahoErrorLogger{}
	mqtt.WARN = PahoErrorLogger{}
	mqtt.DEBUG = PahoDebugLogger{}
}
