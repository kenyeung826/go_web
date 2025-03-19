package util

import (
	"bytes"
	"fmt"
	fwMiddleware "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	AppLogEntryCtxKey = &ContextKey{Name: "AppLogContextKey"}
)

type LoggerInterface interface {
	Output(calldepth int, s string) error
}

type AppLogger struct {
	out *os.File
	*log.Logger
}

var GlobalLog *AppLogger

var globalLogMux sync.Mutex

// GetGlobalLog init
func GetGlobalLog() *AppLogger {
	globalLogMux.Lock()

	defer globalLogMux.Unlock()
	if GlobalLog == nil {
		errLogPath := os.Getenv("Go_Log_Path")
		globalLogFile, err := os.OpenFile(errLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		CheckError(err, nil)

		GlobalLog = &AppLogger{
			Logger: log.New(globalLogFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile),
			out:    globalLogFile,
		}
	}
	return GlobalLog
}

func CloseGlobalLog() {
	globalLogMux.Lock()

	defer globalLogMux.Unlock()
	if GlobalLog != nil {
		GlobalLog.out.Close()
		GlobalLog = nil
	}
}

type Outputter interface {
	Output(calldepth int, s string) error
}

// AppLogOutputter @implement Outputter
type DefaultOutputter struct {
	Logger LoggerInterface
	Buf    *bytes.Buffer
}

func (d DefaultOutputter) Output(calldepth int, s string) error {
	calldepth++
	s = fmt.Sprint(d.Buf.String(), s)
	fmt.Println("App Log Entry")
	return d.Logger.Output(calldepth, s)
}

// AppLogEntry @implement LogEntry
type AppLogEntry struct {
	Outputter Outputter
}

func (e *AppLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
}

func (e *AppLogEntry) Panic(v interface{}, stack []byte) {
	fmt.Println("app logger panic")
	fwMiddleware.PrintPrettyStack(v)
	ps := prettyStack{}
	s, _ := ps.parse(stack, v)
	e.Outputter.Output(2, string(s))
}

func (l *AppLogEntry) Print(v ...any) {
	s := fmt.Sprint(v...)
	fmt.Println("applogger")
	l.Outputter.Output(2, s)
}

// Printf calls l.Output to print to the AppLogger.
// Arguments are handled in the manner of [fmt.Printf].
func (l *AppLogEntry) Printf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	fmt.Println("applogger")
	l.Outputter.Output(2, s)
}

// Println calls l.Output to print to the AppLogger.
// Arguments are handled in the manner of [fmt.Println].
func (l *AppLogEntry) Println(v ...any) {
	s := fmt.Sprintln(v...)
	fmt.Println("applogger")
	l.Outputter.Output(2, s)
}

// Fatal is equivalent to l.Print() followed by a call to [os.Exit](1).
func (l *AppLogEntry) Fatal(v ...any) {
	l.Outputter.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to [os.Exit](1).
func (l *AppLogEntry) Fatalf(format string, v ...any) {
	l.Outputter.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to [os.Exit](1).
func (l *AppLogEntry) Fatalln(v ...any) {
	l.Outputter.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *AppLogEntry) lPanic(v ...any) {
	s := fmt.Sprint(v...)
	l.Outputter.Output(2, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *AppLogEntry) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	l.Outputter.Output(2, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *AppLogEntry) Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	l.Outputter.Output(2, s)
	panic(s)
}
