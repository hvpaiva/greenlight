package uhttp

import "net/http"

type MetricsResponseWriter struct {
	http.ResponseWriter
	StatusCode    int
	HeaderWritten bool
}

func NewMetricsResponseWriter(w http.ResponseWriter) *MetricsResponseWriter {
	return &MetricsResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (mw *MetricsResponseWriter) WriteHeader(statusCode int) {
	mw.ResponseWriter.WriteHeader(statusCode)

	if !mw.HeaderWritten {
		mw.StatusCode = statusCode
		mw.HeaderWritten = true
	}
}

func (mw *MetricsResponseWriter) Write(b []byte) (int, error) {
	mw.HeaderWritten = true
	return mw.ResponseWriter.Write(b)
}

func (mw *MetricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.ResponseWriter
}
