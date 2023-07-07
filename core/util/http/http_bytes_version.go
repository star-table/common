package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/times"
)

func Post(url string, params map[string]interface{}, body []byte, headerOptions ...HeaderOption) ([]byte, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(body))
	req.Header.Set("Content-Type", defaultContentType)

	if err != nil {
		return nil, 0, err
	}

	//tracing Id
	req.Header.Set(consts.TraceIdKey, threadlocal.GetTraceId())
	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}

	headers, _ := json.Marshal(req.Header)
	log.Infof("http type: POST|request [%s] starting|request body [%s]|request headers [%q]", fullUrl, body, headers)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)

	log.Infof("http type: POST| request [%s] successful| request body [%s]|request headers [%q]|response status code [%d]| response body [%s]|time-consuming [%s]", fullUrl, body, headers, httpCode, respBody, timeConsuming)
	return respBody, httpCode, err
}

func Get(url string, params map[string]interface{}, headerOptions ...HeaderOption) ([]byte, int, error) {
	fullUrl := url + ConvertToQueryParams(params)
	req, err := http.NewRequest("GET", fullUrl, nil)

	if err != nil {
		return nil, 0, err
	}
	//tracing Id
	req.Header.Set(consts.TraceIdKey, threadlocal.GetTraceId())
	for _, headerOption := range headerOptions {
		req.Header.Set(headerOption.Name, headerOption.Value)
	}
	headers, _ := json.Marshal(req.Header)
	log.Infof("http type: GET|request [%s] starting|request headers [%q]", fullUrl, headers)

	start := times.GetNowMillisecond()
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	end := times.GetNowMillisecond()
	timeConsuming := strconv.FormatInt(end-start, 10)

	respBody, httpCode, err := responseHandle(resp, err)
	log.Infof("http type: GET| request [%s] successful|request headers [%q]|response status code [%d]|response body [%s]| time-consuming [%s]", fullUrl, headers, httpCode, respBody, timeConsuming)
	return respBody, httpCode, err
}

func ConvertToQueryParams(params map[string]interface{}) string {
	if params == nil || len(params) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for k, v := range params {
		if v == nil {
			continue
		}
		if fv, ok := v.(float64); ok {
			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatFloat(fv, 'f', -1, 64)))
		} else if fv, ok := v.(float32); ok {
			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatFloat(float64(fv), 'f', -1, 32)))
		} else if iv, ok := v.(int64); ok {
			buffer.WriteString(fmt.Sprintf("%s=%s&", k, strconv.FormatInt(iv, 10)))
		} else {
			buffer.WriteString(fmt.Sprintf("%s=%v&", k, v))
		}
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}

func responseHandle(resp *http.Response, err error) ([]byte, int, error) {
	if err != nil {
		log.Error(err)
		return nil, 500, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, resp.StatusCode, err
	}
	return b, resp.StatusCode, nil
}
