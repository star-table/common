package times

import (
	"fmt"
	"testing"
	"time"
)

func TestGetDateTimeStrByMillisecond(t *testing.T) {
	//ms := GetNowMillisecond()
	//
	//fmt.Println(ms)
	//
	//fmt.Println(GetDateTimeStrByMillisecond(ms))

	//Sleep(5)
	//fmt.Println("sleep")

	SleepMillisecond(100)
	fmt.Println("sleep ms")
	//fmt.Println(GetUnixTime("2019-01-02 15:04:05"))
}

func TestGetTime0(t *testing.T) {
	fmt.Println(GetBeiJingTime(), time.Now())
}

func TestGetMillisecond(t *testing.T) {
	fmt.Println(GetMillisecond(time.Now()))
}

func TestGetTodayTimeQuantum(t *testing.T) {
	t.Log(GetTodayTimeQuantum())
}

func TestGetYesterdayDate(t *testing.T) {
	t.Log(GetYesterdayDate())
}
