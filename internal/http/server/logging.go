package server

import (
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
)

type interceptingWriter struct {
	http.ResponseWriter
	code    int
	written int64
}

// WriteHeader may not be explicitly called, so care must be taken to
// initialize w.code to its default value of http.StatusOK.
func (w *interceptingWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *interceptingWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	w.written += int64(n)
	return n, err
}

type loggingMiddleware struct {
	logger log.Logger
	next   http.Handler
}

func (l *loggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	iw := &interceptingWriter{w, http.StatusOK, 0}
	defer func(from time.Time) {
		_ = l.logger.Log(
			"path", r.URL.Path,
			"statusCode", iw.code,
			"responseBodySize", iw.written,
			"executionTime", time.Since(from),
		)
	}(time.Now())
	l.next.ServeHTTP(iw, r)
}

func NewLoggingMiddleware(logger log.Logger, next http.Handler) http.Handler {
	return &loggingMiddleware{
		logger: logger,
		next:   next,
	}
}
