package middleware

import (
	"bytes"
	"net/http"
)

// @implement http.ResponseWriter
type bufferedResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	Buffer      *bytes.Buffer
	HeadersSent bool
}

func (brw *bufferedResponseWriter) Header() http.Header {
	return brw.ResponseWriter.Header()
}

func (brw *bufferedResponseWriter) Write(data []byte) (int, error) {
	return brw.Buffer.Write(data)
}

func (brw *bufferedResponseWriter) WriteHeader(statusCode int) {
	brw.StatusCode = statusCode
}

func (brw *bufferedResponseWriter) Send() {
	if !brw.HeadersSent {
		brw.ResponseWriter.WriteHeader(brw.StatusCode)
		brw.HeadersSent = true
	}
	// Write the buffered content to the original ResponseWriter
	brw.Buffer.WriteTo(brw.ResponseWriter)
}

func BufferedResponseHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		brw := &bufferedResponseWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
			Buffer:         &bytes.Buffer{},
			HeadersSent:    false,
		}
		next.ServeHTTP(brw, r)
		brw.Send()
	}

	return http.HandlerFunc(fn)
}
