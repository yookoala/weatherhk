package ctxlog

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
)

var logger log.Logger

var isHeroku bool

func init() {
	SetLogOutput(os.Stdout)
	if strings.ToLower(os.Getenv("ON_HEROKU")) == "true" {
		isHeroku = true
	}
}

// SetLogOutput sets the output for log messages
func SetLogOutput(out io.Writer) {
	logger = log.NewLogfmtLogger(out)
}

func currentTimestamp() interface{} {
	now := time.Now()
	return now.Format("2006-01-02T15:04:05") +
		fmt.Sprintf(".%09d", now.Nanosecond()) +
		now.Format(" -0700 MST")
}

// GetLoggers return loggers for general use
func GetLoggers(r *http.Request) (infoLog, errorLog log.Logger) {
	requestID := r.Header.Get("X-Request-ID")
	if isHeroku {
		basicCtx := log.NewContext(logger).WithPrefix(
			"method", r.Method,
			"url", r.URL.EscapedPath(),
			"request_id", requestID)
		infoLog = basicCtx.WithPrefix("at", "info")
		errorLog = basicCtx.WithPrefix("at", "error")
		return
	}
	basicCtx := log.NewContext(logger).WithPrefix(
		"ts", log.Valuer(currentTimestamp),
		"method", r.Method,
		"url", r.URL.EscapedPath(),
		"request_id", requestID)
	infoLog = basicCtx.WithPrefix("at", "info")
	errorLog = basicCtx.WithPrefix("at", "error")

	return
}
