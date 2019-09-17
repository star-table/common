package web

import (
	"net/http"
)

//Http请求处理代理
type HttpHandleProxy struct {
	Before func(http.ResponseWriter, *http.Request)
	After  func(http.ResponseWriter, *http.Request)
	Except func(http.ResponseWriter, *http.Request)
}

type ResponseProxyWriter struct {
	writer http.ResponseWriter
	Output []byte
}

func (this *ResponseProxyWriter) Header() http.Header {
	return this.writer.Header()
}
func (this *ResponseProxyWriter) Write(bytes []byte) (int, error) {
	this.Output = append(this.Output, bytes[0:len(bytes)]...)
	return this.writer.Write(bytes)
}
func (this *ResponseProxyWriter) WriteHeader(i int) {
	this.writer.WriteHeader(i)
}

func NewRespProxyWriter(w http.ResponseWriter) *ResponseProxyWriter {
	return &ResponseProxyWriter{
		writer: w,
		Output: []byte{},
	}
}

func (this *HttpHandleProxy) For(handle func(http.ResponseWriter,
	*http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if this.Except != nil {
			defer func() {
				if err := recover(); err != nil {
					this.Except(w, r)
				}
			}()
		}

		if this.Before != nil {
			this.Before(w, r)
		}

		proxy := NewRespProxyWriter(w)

		if handle != nil {
			handle(proxy, r)
		}

		if this.After != nil {
			this.After(proxy, r)
		}
	}
}
