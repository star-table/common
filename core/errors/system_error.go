package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

var (
	_codes = make(map[int]ResultCodeInfo) // 注册返回编号信息
)

//result code info
type ResultCodeInfo struct {
	code int

	message string

	langCode string
	// 异常时，是否发送到 sentry
	isNeedSentry bool
}

type SystemErrorInfo interface {
	Error() string

	Code() int

	Message() string

	LangCode() string

	Equal(error) bool
}



func (sei ResultCodeInfo) Error() string {
	return "{\"code\":" + strconv.Itoa(sei.code) + ",\"message\":\"" + sei.message + "\"}"
}

func (sei ResultCodeInfo) Code() int {
	return sei.code
}

func (sei ResultCodeInfo) Message() string {
	return sei.message
}

func (sei ResultCodeInfo) LangCode() string {
	return sei.langCode
}

func (sei *ResultCodeInfo) SetCode(code int) {
	sei.code = code
}

func (sei *ResultCodeInfo) SetMessage(message string) {
	sei.message = message
}

func (sei *ResultCodeInfo) SetLangCode(langCode string) {
	sei.langCode = langCode
}

func (sei ResultCodeInfo) Equal(err error) bool {
	return EqualError(sei, err)
}

func GetResultCodeInfoByCode(code int) ResultCodeInfo {

	if v, ok := _codes[code]; ok {
		return v
	} else {
		return SystemError
	}
}

func getResultCodeInfoByString(e string) ResultCodeInfo {
	if e == "" {
		return OK
	}

	i, err := strconv.Atoi(e)
	if err != nil {
		return SystemError
	}
	if v, ok := _codes[i]; ok {
		return v
	} else {
		return SystemError
	}
}

func convertSystemErrorInfo(e error) SystemErrorInfo {
	if e == nil {
		return OK
	}
	ec, ok := errors.Cause(e).(SystemErrorInfo)
	if ok {
		return ec
	}
	return getResultCodeInfoByString(e.Error())
}

func EqualError(code ResultCodeInfo, err error) bool {
	return convertSystemErrorInfo(err).Code() == code.Code()
}

//add system ResultCodeInfo
func AddResultCodeInfo(code int, message string, langCode string) ResultCodeInfo {
	if code < 0 {
		panic(fmt.Sprintf("result code: code %d must greater than zero", code))
	}
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("result code: %d already exist", code))
	}

	rci := ResultCodeInfo{
		code:     code,
		message:  message,
		langCode: langCode,
	}

	_codes[code] = rci
	return rci
}

// 构造错误实例的构造函数。它构造的实例，带有发送异常到 sentry 的标识。
func AddResultCodeInfoWithSentry(code int, message string, langCode string) ResultCodeInfo {
	if code < 0 {
		panic(fmt.Sprintf("result code: code %d must greater than zero", code))
	}
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("result code: %d already exist", code))
	}

	rci := ResultCodeInfo{
		code:         code,
		message:      message,
		langCode:     langCode,
		isNeedSentry: true,
	}

	_codes[code] = rci
	return rci
}

func BuildSystemErrorInfo(resultCodeInfo ResultCodeInfo, e ...error) SystemErrorInfo {
	if e == nil || (e != nil && len(e) == 0){
		return resultCodeInfo
	}
	if e != nil{
		if reflect.TypeOf(e[0]).Name() == "ResultCodeInfo"{
			v := e[0].(ResultCodeInfo)
			resultCodeInfo = v
		}else{
			resultCodeInfo.message += "," + e[0].Error()
		}
	}
	return resultCodeInfo
}

func BuildSystemErrorInfoWithMessage(resultCodeInfo ResultCodeInfo, message string) SystemErrorInfo {
	resultCodeInfo.message += message
	return resultCodeInfo
}
