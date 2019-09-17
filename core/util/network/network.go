package network

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	_lock    sync.Mutex
	_localIp = ""
)

func GetIntranetIp() string {

	if _localIp != "" {
		return _localIp
	}

	_lock.Lock()

	defer _lock.Unlock()

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				_localIp = ipnet.IP.String()
				break
			}
		}
	}

	if _localIp == "" {
		_localIp = "127.0.0.1"
	}
	return _localIp
}
