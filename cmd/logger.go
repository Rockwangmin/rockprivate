package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	Trace           *log.Logger
	Info            *log.Logger
	Warning         *log.Logger
	Error           *log.Logger
	Audit           *log.Logger
	runtimeLog, err = os.OpenFile("iot-cli.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	auditLog, err2 = os.OpenFile("audit.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
)

func init() {
	if err != nil || err2 != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(io.MultiWriter(runtimeLog, os.Stderr),
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(io.MultiWriter(runtimeLog),
		"INFO: ",
		log.Ldate|log.Ltime)

	Warning = log.New(io.MultiWriter(runtimeLog),
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(io.MultiWriter(runtimeLog),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Audit = log.New(io.MultiWriter(auditLog),
		"Audit: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func GetLogger(level string) *log.Logger {
	switch level {
	case "Info":
		return Info
	case "Trace":
		return Trace
	case "Warn":
		return Warning
	case "Error":
		return Error
	case "Audit":
		return Error
	}
	return nil
	INFO
}

func NewLogger(level string, prefix string) *log.Logger {
	var flag int
	switch level {
	case "Info", "Warn", "Audit":
		flag = log.Ldate | log.Ltime
	case "Trace", "Error":
		flag = log.Ldate | log.Ltime | log.Lshortfile
	}
	return log.New(io.MultiWriter(runtimeLog),
		fmt.Sprintf("%5s : %10s: ", level, prefix),
		flag)
}
