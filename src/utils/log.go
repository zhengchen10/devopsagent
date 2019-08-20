package utils

import (
	"fmt"
	"time"
)

type Log struct {
	logFileName  string
	logFile      File
	logLevel     int
	logFileLevel int
}

//var logFile = ""
func (l *Log) InitLog(fileName string, logLevel string, logFileLevel string) {
	l.logFileName = fileName
	l.logFile = File{Path: l.logFileName}
	l.logLevel = 0
	if logLevel == "DEBUG" {
		l.logLevel = 0
	} else if logLevel == "INFO" {
		l.logLevel = 1
	} else if logLevel == "WARN" {
		l.logLevel = 2
	} else if logLevel == "ERROR" {
		l.logLevel = 3
	} else if logLevel == "NONE" {
		l.logLevel = 4
	}
	l.logFileLevel = 3
	if logFileLevel == "DEBUG" {
		l.logFileLevel = 0
	} else if logFileLevel == "INFO" {
		l.logFileLevel = 1
	} else if logFileLevel == "WARN" {
		l.logFileLevel = 2
	} else if logFileLevel == "ERROR" {
		l.logFileLevel = 3
	} else if logFileLevel == "NONE" {
		l.logFileLevel = 4
	}
}
func (l *Log) LogA(level, address, msg string, showConsole, showFile bool) {
	newmsg := address + " " + msg
	l.Log(level, newmsg, showConsole, showFile)
}

func (l *Log) Log(level, msg string, showConsole, showFile bool) {
	current := time.Now()
	year := current.Year()
	month := current.Month()
	day := current.Day()
	hour := current.Hour()     //小时
	minute := current.Minute() //分钟
	second := current.Second() //秒
	var logMsg string
	logMsg = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d [%s] %s", year, month, day, hour, minute, second, level, msg)
	fmt.Println(logMsg)
	l.logFile.AppendContent(logMsg + "\r\n")
}

func (l *Log) Debug(msg string) {
	l.Log("DEBUG", msg, l.logLevel < 1, l.logFileLevel < 1)
}

func (l *Log) Info(msg string) {
	l.Log("INFO", msg, l.logLevel < 2, l.logFileLevel < 2)
}

func (l *Log) Warn(msg string) {
	l.Log("WARN", msg, l.logLevel < 3, l.logFileLevel < 3)
}

func (l *Log) Error(msg string) {
	l.Log("ERROR", msg, l.logLevel < 4, l.logFileLevel < 4)
}

func (l *Log) DebugA(address, msg string) {
	l.LogA("DEBUG", address, msg, l.logLevel < 1, l.logFileLevel < 1)
}

func (l *Log) InfoA(address, msg string) {
	l.LogA("INFO", address, msg, l.logLevel < 2, l.logFileLevel < 2)
}

func (l *Log) WarnA(address, msg string) {
	l.LogA("WARN", address, msg, l.logLevel < 3, l.logFileLevel < 3)
}

func (l *Log) ErrorA(address, msg string) {
	l.LogA("ERROR", address, msg, l.logLevel < 4, l.logFileLevel < 4)
}
