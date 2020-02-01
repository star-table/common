package times

import (
	"github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/common/core/types"
	"github.com/galaxy-book/common/core/util/date"
	"time"
)

func GetNowMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetMillisecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func GetNowSecond() int64 {
	return time.Now().Unix()
}

func GetNowNanoSecond() int64 {
	return time.Now().UnixNano()
}

func GetBeiJingTime() time.Time {
	timelocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timelocal
	return time.Now().Local()
}

func Sleep(second int64) {
	time.Sleep(time.Duration(second) * time.Second)
}

func SleepMillisecond(millSecond int64) {
	time.Sleep(time.Duration(millSecond) * time.Millisecond)
}

func GetDateTimeStrBySecond(s int64) string {
	return time.Unix(s, 0).Format(consts.AppTimeFormat)
}

func GetDateTimeStrByMillisecond(ms int64) string {
	second := ms / 1000
	return time.Unix(second, 0).Format(consts.AppTimeFormat)
}

func GetUnixTime(t types.Time) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")                                 //设置时区
	tt, _ := time.ParseInLocation(consts.AppTimeFormat, date.FormatTime(t), loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.Unix()
}

func GetWeeHours() string {
	timeStr := time.Now().Format(consts.AppDateFormat)
	t, _ := time.Parse(consts.AppDateFormat, timeStr)
	return date.Format(t)
}

//获取当天时间段：2019-08-12 00:00:00 - 2019-08-12 23:59:59
func GetTodayTimeQuantum() []time.Time {
	timeStr := time.Now().Format(consts.AppDateFormat)
	b, _ := time.Parse(consts.AppDateFormat, timeStr)
	a := b.Add(time.Duration(1000*60*60*24-1) * time.Millisecond)
	return []time.Time{b, a}
}

//2019-09-03
func GetYesterdayDate() string{
	timeStr := time.Now().Format(consts.AppDateFormat)
	b, _ := time.Parse(consts.AppDateFormat, timeStr)
	a := b.Add(time.Duration(-1) * time.Millisecond)
	return a.Format(consts.AppDateFormat)
}

//func GetTime0() time.Time {
//	return time.Unix(0, 0)
//}
