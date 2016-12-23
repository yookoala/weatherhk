package ctxlog

import (
	"io"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

var logger log.Logger

func init() {
	SetLogOutput(os.Stdout)
}

// SetLogOutput sets the output for log messages
func SetLogOutput(out io.Writer) {
	logger = log.NewLogfmtLogger(out)
}

// GetLoggers return loggers for general use
func GetLoggers(r *http.Request) (infoLog, errorLog log.Logger) {
	requestID := r.Header.Get("X-Request-ID")
	basicCtx := log.NewContext(logger).WithPrefix(
		"method", r.Method,
		"url", r.URL.EscapedPath(),
		"request_id", requestID)
	infoLog = basicCtx.WithPrefix("at", "info")
	errorLog = basicCtx.WithPrefix("at", "error")
	return
}
