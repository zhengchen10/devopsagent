package utils

import (
	"fmt"
	"time"
)

type Log struct {
	logFileName string
	logFile     File
}

//var logFile = ""
func (l *Log) InitLog(fileName string) {
	l.logFileName = fileName
	l.logFile = File{Path: l.logFileName}

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
	l.Log("DEBUG", msg)
}
func (l *Log) Info(msg string) {
	l.Log("INFO", msg)
}
func (l *Log) Warn(msg string) {
	l.Log("WARN", msg)
}

func (l *Log) Error(msg string) {
	l.Log("ERROR", msg)
}

func (l *Log) DebugA(address, msg string) {
	l.LogA("DEBUG", address, msg)
}

func (l *Log) InfoA(address, msg string) {
	l.LogA("INFO", address, msg)
}

func (l *Log) WarnA(address, msg string) {
	l.LogA("WARN", address, msg)
}

func (l *Log) ErrorA(address, msg string) {
	l.LogA("ERROR", address, msg)
}
