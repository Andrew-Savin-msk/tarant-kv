package apiserver

import "net/http"

type responseWrighter struct {
	http.ResponseWriter
	code int
}

func (w *responseWrighter) WrightHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader((statusCode))
}
