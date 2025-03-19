package middleware

import (
	"app/util"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	fwMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewAppLogEntry(logger *util.AppLogger, r *http.Request) *util.AppLogEntry {
	buf := &bytes.Buffer{}
	reqID := fwMiddleware.GetReqID(r.Context())
	if reqID != "" {
		_, _ = fmt.Fprintf(buf, "[%s] ", reqID)
	}

	outputter := &util.DefaultOutputter{
		Logger: logger,
		Buf:    buf,
	}

	return &util.AppLogEntry{
		Outputter: outputter,
	}
}

// ServerLogEntry @implement LogEntry
type ServerLogEntry struct {
	RequestLogger fwMiddleware.LogEntry
	ErrorLogger   fwMiddleware.LogEntry
}

func NewServerLogEntry(requestLogger fwMiddleware.LogEntry, errorLogger fwMiddleware.LogEntry) *ServerLogEntry {
	return &ServerLogEntry{RequestLogger: requestLogger, ErrorLogger: errorLogger}
}

func (l *ServerLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.RequestLogger.Write(status, bytes, header, elapsed, extra)
}

func (l *ServerLogEntry) Panic(v interface{}, stack []byte) {
	l.ErrorLogger.Panic(v, stack)
}

func ServerLogger(next http.Handler) http.Handler {
	requestLogPath := os.Getenv("Request_Log_Path")
	appLogPath := os.Getenv("App_Log_Path")

	fn := func(w http.ResponseWriter, r *http.Request) {
		//ACCESS log
		requestLogFile, err := os.OpenFile(requestLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		util.CheckError(err, nil)
		multiWriter := io.MultiWriter(requestLogFile, os.Stdout)
		requestLogger := log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		rf := &fwMiddleware.DefaultLogFormatter{Logger: requestLogger, NoColor: true}
		re := rf.NewLogEntry(r)

		errorLogger := util.GetGlobalLog()
		ee := NewAppLogEntry(errorLogger, r)

		sle := NewServerLogEntry(re, ee)

		//application log
		appLogFile, err := os.OpenFile(appLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		util.CheckError(err, nil)
		appLogger := &util.AppLogger{
			Logger: log.New(appLogFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
		}
		ae := NewAppLogEntry(appLogger, r)

		ww := fwMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			sle.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			_ = requestLogFile.Close()
			_ = appLogFile.Close()
		}()

		r = r.WithContext(context.WithValue(r.Context(), fwMiddleware.LogEntryCtxKey, sle))
		r = r.WithContext(context.WithValue(r.Context(), util.AppLogEntryCtxKey, ae))

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}

func GetAppLogger(r *http.Request) *util.AppLogEntry {
	v := r.Context().Value(util.AppLogEntryCtxKey)
	entry, ok := v.(*util.AppLogEntry)
	if !ok {
		panic("No Key in context")
	}
	return entry
}
