package random

import (
	"fmt"
	"github.com/star-table/common/core/util/times"
	"github.com/star-table/common/core/util/uuid"
	"math/rand"
	"strconv"
	"time"
)

func Token() string {
	return uuid.NewUuid()
}

// 生成带有组织 id 和用户 id 的token
func GenTokenWithOrgId(orgId, userId int64) string {
	return fmt.Sprintf("o%[1]su%[2]st%[3]s", strconv.FormatInt(orgId, 10), strconv.FormatInt(userId, 10), Token())
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
