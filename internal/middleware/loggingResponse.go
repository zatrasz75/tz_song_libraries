package middleware

import (
	"github.com/gofrs/uuid/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqID := req.URL.Query().Get("request_id")
		if reqID == "" {
			rID, _ := uuid.NewV4()
			reqID = rID.String()
		}
		req.Header.Set("X-Request-ID", reqID)
		w.Header().Set("X-Request-ID", reqID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		str := strings.Builder{}
		str.WriteString("<-- ip: ")
		str.WriteString(req.RemoteAddr)
		str.WriteString(", host: ")
		str.WriteString(req.Host)
		str.WriteString(" url: ")
		str.WriteString(req.URL.Path)
		str.WriteString(", method: ")
		str.WriteString(req.Method)
		str.WriteString(" status code: ")
		str.WriteString(strconv.Itoa(lrw.statusCode))
		str.WriteString(" ")
		str.WriteString(http.StatusText(lrw.statusCode))
		str.WriteString(", trace id: ")
		str.WriteString(reqID)

		log.Println(str.String())
	})
}
