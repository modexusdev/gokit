package log

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	return rw.ResponseWriter.Write(data)
}

func (l *Logger) HTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			writer := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(writer, r)

			l.request(
				r.Method,
				r.URL.Path,
				writer.status,
				time.Since(start),
			)
		},
	)
}

func (l *Logger) request(
	method string,
	path string,
	status int,
	duration time.Duration,
	messages ...string,
) {
	now := time.Now().Format("15:04:05")

	level, levelColor := statusLevel(status)

	fmt.Printf(
		"%s[%s]%s "+
			"%s[%s]%s "+
			"%s%-8s%s "+
			"%s%-7s%s "+
			"%-28s "+
			"%s%d%s "+
			"%s%8s%s",

		// Time
		colorGray,
		now,
		colorReset,

		// Application
		colorBrightMagenta,
		l.name,
		colorReset,

		// Level
		levelColor,
		level,
		colorReset,

		// HTTP method
		methodColor(method),
		method,
		colorReset,

		// Path
		path,

		// Status
		statusColor(status),
		status,
		colorReset,

		// Duration
		colorGray,
		formatDuration(duration),
		colorReset,
	)

	message := strings.Join(messages, " ")

	if message != "" {
		fmt.Printf("  %s", message)
	}

	fmt.Println()
}

func statusLevel(status int) (string, string) {
	switch {
	case status >= 500:
		return "INTERNAL", colorBrightRed

	case status >= 400:
		return "ERROR", colorBrightYellow

	case status >= 300:
		return "REDIRECT", colorBrightCyan

	case status >= 200:
		return "INFO", colorBrightGreen

	default:
		return "INFO", colorGray
	}
}

func statusColor(status int) string {
	switch {
	case status >= 500:
		return colorBrightRed

	case status >= 400:
		return colorBrightYellow

	case status >= 300:
		return colorBrightCyan

	case status >= 200:
		return colorBrightGreen

	default:
		return colorGray
	}
}

func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return colorBrightCyan

	case http.MethodPost:
		return colorBrightMagenta

	case http.MethodPut:
		return colorBrightYellow

	case http.MethodPatch:
		return colorBrightYellow

	case http.MethodDelete:
		return colorBrightRed

	default:
		return colorWhite
	}
}

func formatDuration(duration time.Duration) string {
	if duration < time.Millisecond {
		return fmt.Sprintf(
			"%dµs",
			duration.Microseconds(),
		)
	}

	return fmt.Sprintf(
		"%.2fms",
		float64(duration.Microseconds())/1000,
	)
}
