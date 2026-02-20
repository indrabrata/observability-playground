package common

import "net/http"

type Interceptor struct {
	// Note : By embedding `http.ResponseWriter` directly in the struct, Go automatically promotes all of its methods to Interceptor
	http.ResponseWriter
	StatusCode int
}

func NewInterceptor(w http.ResponseWriter) *Interceptor {
	return &Interceptor{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (crw *Interceptor) WriteHeader(statusCode int) {
	crw.StatusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}

func (crw *Interceptor) Write(b []byte) (int, error) {
	return crw.ResponseWriter.Write(b)
}
