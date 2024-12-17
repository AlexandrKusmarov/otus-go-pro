package internalhttp

import (
	"log"
	"net/http"
	"time"
)

// middlewareLogger — middleware для логирования запросов
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Обёртка для отслеживания статуса ответа
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Выполнение обработчика
		next.ServeHTTP(ww, r)

		// Логирование запроса
		latency := time.Since(start)
		log.Printf("%s [%s] %s %s %s %d %d \"%s\"",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			ww.statusCode,
			latency.Milliseconds(),
			r.UserAgent(),
		)
	})
}

// responseWriterWrapper — обёртка для ResponseWriter
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
