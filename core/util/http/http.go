package http

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/star-table/common/core/logger"
)

var log = logger.GetDefaultLogger()

const defaultContentType = "application/json"

var httpClient = &http.Client{}

type HeaderOption struct {
	Name  string
	Value string
}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(30) * time.Second,
	}
}

//func Post(url string, params map[string]interface{}, body string, headerOptions ...HeaderOption) (string, int, error) {
//	fullUrl := url + ConvertToQueryParams(params)
//	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(body))
//	req.Header.Set("Content-Type", defaultContentType)
//
//	if err != nil {
//		return consts.BlankString, 0, err
//	}
//
//	//tracing Id
//	req.Header.Set(consts.TraceIdKey, threadlocal.GetTraceId())
//	for _, headerOption := range headerOptions {
//		req.Header.Set(headerOption.Name, headerOption.Value)
//	}
//
//	headers := json.ToJsonIgnoreError(req.Header)
//	log.Infof("http type: POST|request [%s] starting|request body [%s]|request headers [%s]", fullUrl, body, headers)
//
//	start := times.GetNowMillisecond()
//	resp, err := httpClient.Do(req)
//
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//
//	end := times.GetNowMillisecond()
//	timeConsuming := strconv.FormatInt(end-start, 10)
//
//	respBody, httpCode, err := responseHandle(resp, err)
//
//	log.Infof("http type: POST| request [%s] successful| request body [%s]|request headers [%s]|response status code [%d]| response body [%s]|time-consuming [%s]", fullUrl, body, headers, httpCode, respBody, timeConsuming)
//	return respBody, httpCode, err
//}
//
//func Get(url string, params map[string]interface{}, headerOptions ...HeaderOption) (string, int, error) {
//	fullUrl := url + ConvertToQueryParams(params)
//	req, err := http.NewRequest("GET", fullUrl, nil)
//
//	if err != nil {
//		return consts.BlankString, 0, err
//	}
//	//tracing Id
//	req.Header.Set(consts.TraceIdKey, threadlocal.GetTraceId())
//	for _, headerOption := range headerOptions {
//		req.Header.Set(headerOption.Name, headerOption.Value)
//	}
//	headers := json.ToJsonIgnoreError(req.Header)
//	log.Infof("http type: GET|request [%s] starting|request headers [%s]", fullUrl, headers)
//
//	start := times.GetNowMillisecond()
//	resp, err := httpClient.Do(req)
//
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//
//	end := times.GetNowMillisecond()
//	timeConsuming := strconv.FormatInt(end-start, 10)
//
//	respBody, httpCode, err := responseHandle(resp, err)
//	log.Infof("http type: GET| request [%s] successful|request headers [%s]|response status code [%d]|response body [%s]| time-consuming [%s]", fullUrl, headers, httpCode, respBody, timeConsuming)
//	return respBody, httpCode, err
//}
//
//func responseHandle(resp *http.Response, err error) (string, int, error) {
//	if err != nil {
//		log.Error(err)
//		return "", 500, err
//	}
//	b, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Error(err)
//		return "", resp.StatusCode, err
//	}
//	return string(b), resp.StatusCode, nil
//}
//
//func ConvertToQueryParams(params map[string]interface{}) string {
//	paramsJson := json.ToJsonIgnoreError(params)
//	params = map[string]interface{}{}
//	_ = json.FromJson(paramsJson, &params)
//
//	if &params == nil || len(params) == 0 {
//		return ""
//	}
//	var buffer bytes.Buffer
//	buffer.WriteString("?")
//	for k, v := range params {
//		if v == nil {
//			continue
//		}
//		if fv, ok := v.(float64); ok {
//			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatFloat(fv, 'f', -1, 64)))
//		} else if fv, ok := v.(float32); ok {
//			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatFloat(float64(fv), 'f', -1, 32)))
//		} else if iv, ok := v.(int64); ok {
//			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatInt(iv, 10)))
//		} else {
//			buffer.WriteString(fmt.Sprintf("%s=%v&", k, v))
//		}
//	}
//	buffer.Truncate(buffer.Len() - 1)
//	return buffer.String()
//}
