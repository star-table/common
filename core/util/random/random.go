package random

import (
	"github.com/galaxy-book/common/core/util/times"
	"github.com/galaxy-book/common/core/util/uuid"
	"math/rand"
	"strconv"
	"time"
)

func Token() string {
	return uuid.NewUuid()
}

func RandomString(l int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandomFileName() string {
	return uuid.NewUuid() + strconv.FormatInt(times.GetNowMillisecond(), 10)
}
