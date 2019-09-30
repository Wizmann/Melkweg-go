package Gwisted

import (
    "runtime"
    "path/filepath"
    "fmt"
    "time"
    "os"
)

var logLevel = NOTSET

type LogLevel int

const (
    NOTSET   LogLevel = 0
    VERBOSE  LogLevel = 1
    DEBUG    LogLevel = 2
    INFO     LogLevel = 3
    WARNING  LogLevel = 4
    ERROR    LogLevel = 5
    FATAL    LogLevel = 6
)

func Verbose(format string, a ...interface{}) {
    if logLevel <= VERBOSE {
        logAux("Verbose", format, a...)
    }
}

func Debug(format string, a ...interface{}) {
    if logLevel <= DEBUG {
        logAux("Debug", format, a...)
    }
}

func Info(format string, a ...interface{}) {
    if logLevel <= INFO {
        logAux("Info", format, a...)
    }
}

func Warning(format string, a ...interface{}) {
    if logLevel <= WARNING {
        logAux("Warning", format, a...)
    }
}

func Error(format string, a ...interface{}) {
    if logLevel <= ERROR {
        logAux("Error", format, a...)
    }
}

func Fatal(format string, a ...interface{}) {
    if logLevel <= FATAL {
        logAux("Fatal", format, a...);
    }
}

func logAux(level string, format string, a ...interface{}) {
    _, path, lineno, _ := runtime.Caller(2);
    wd, _ := os.Getwd()
    relpath, _ := filepath.Rel(wd, path)

    t := time.Now();
    a = append([]interface{} { level, t.Format("2006-01-02 15:04:05.00"), relpath, lineno }, a...);
    fmt.Printf("[%s] %s [%s:%d] " + format + "\n", a...);
}

