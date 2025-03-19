package main

import (
	"app/util"
	"log"
	"os"
	"time"
)

const logPathPrefix = ""

var loggerMap map[string]*log.Logger = make(map[string]*log.Logger)

func New(name string) (*log.Logger, error) {
	logger, ok := loggerMap[name]
	if !ok {
		var curdate time.Time
		logPath := logPathPrefix + curdate.Format("") + ".log"
		out, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		util.CheckError(err, "Error opening log file")

		logger = log.New(out, "", log.LstdFlags)
		loggerMap[name] = logger
	}
	return logger, nil
}

func Close() {

}
