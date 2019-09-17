package lock

import (
	"sync"
)

var lock = sync.Mutex{}
var lockMap = sync.Map{}

func Lock(key string) {
	mu, _ := lockMap.Load(key)
	if mu == nil {
		lock.Lock()
		defer lock.Unlock()
		mu, _ = lockMap.Load(key)
		if mu == nil {
			lockMap.Store(key, &sync.Mutex{})
			mu, _ = lockMap.Load(key)
		}
	}
	mu.(*sync.Mutex).Lock()
}
func Unlock(key string) {
	mu, _ := lockMap.Load(key)
	if mu != nil {
		mu.(*sync.Mutex).Unlock()
	}
}
