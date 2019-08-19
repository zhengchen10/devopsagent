package utils

import (
	"fmt"
	"time"
)

type Log struct {
	logFileName string
	logFile     File
	logLevel    int
}

//var logFile = ""
func (l *Log) InitLog(fileName string, logLevel string) {
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

}
func (l *Log) LogA(level, address, msg string) {
	newmsg := address + " " + msg
	l.Log(level, newmsg)
}

func (l *Log) Log(level, msg string) {
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
	if l.logLevel < 1 {
		l.Log("DEBUG", msg)
	}
}
func (l *Log) Info(msg string) {
	if l.logLevel < 2 {
		l.Log("INFO", msg)
	}
}
func (l *Log) Warn(msg string) {
	if l.logLevel < 3 {
		l.Log("WARN", msg)
	}
}

func (l *Log) Error(msg string) {
	if l.logLevel < 4 {
		l.Log("ERROR", msg)
	}
}

func (l *Log) DebugA(address, msg string) {
	if l.logLevel < 1 {
		l.LogA("DEBUG", address, msg)
	}
}

func (l *Log) InfoA(address, msg string) {
	if l.logLevel < 2 {
		l.LogA("INFO", address, msg)
	}
}

func (l *Log) WarnA(address, msg string) {
	if l.logLevel < 3 {
		l.LogA("WARN", address, msg)
	}
}

func (l *Log) ErrorA(address, msg string) {
	if l.logLevel < 4 {
		l.LogA("ERROR", address, msg)
	}
}
