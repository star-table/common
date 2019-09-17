package file

import (
	"path"
	"runtime"
)

func GetCurrentPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath := path.Dir(filename)
		return abPath
	}
	return ""
}
