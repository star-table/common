package logger

import (
	"github.com/star-table/common/core/config"
	"testing"
)

//var _logpath string
//
//func TestMain(m *testing.M) {
//	fmt.Println("begin")
//
//	_logpath = filepath.Join(config.GetLogPath(), "log.log")
//	fmt.Println(config.GetLogPath())
//	fmt.Println(_logpath)
//
//	m.Run()
//
//	fmt.Println("end")
//}
//
//func TestInfo(t *testing.T) {
//
//	logger := GetLogger(_logpath, "info")
//
//	logger.Info("log info ")
//}

func TestLogger(m *testing.T) {
	config.LoadUnitTestConfig()

	var log = GetDefaultLogger()

	log.Info("2132113")

}
