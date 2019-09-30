package Gwisted

import (
    "runtime"
    "path/filepath"
    "fmt"
    "time"
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
        _, path, lineno, ok := runtime.Caller(1);
        _, file := filepath.Split(path)


        if (ok) {
            logAux("Verbose", file, lineno, format, a...);
        }
    }
}

func Debug(format string, a ...interface{}) {
    if logLevel <= DEBUG {
        _, path, lineno, ok := runtime.Caller(1);
        _, file := filepath.Split(path)


        if (ok) {
            logAux("Debug", file, lineno, format, a...);
        }
    }
}

func Info(format string, a ...interface{}) {
    if logLevel <= INFO {
        _, path, lineno, ok := runtime.Caller(1);
        _, file := filepath.Split(path)


        if (ok) {
            logAux("Info", file, lineno, format, a...);
        }
    }
}

func Warning(format string, a ...interface{}) {
    if logLevel <= WARNING {
        _, path, lineno, ok := runtime.Caller(1);
        _, file := filepath.Split(path)


        if (ok) {
            logAux("Warning", file, lineno, format, a...);
        }
    }
}

func Error(format string, a ...interface{}) {
    if logLevel <= ERROR {
        _, path, lineno, ok := runtime.Caller(1);
        _, file := filepath.Split(path)


        if (ok) {
            logAux("Error", file, lineno, format, a...);
        }
    }
}

func Fatal(format string, a ...interface{}) {
    if logLevel <= FATAL {
        _, path, lineno, ok := runtime.Caller(1);
        _, file := filepath.Split(path)


        if (ok) {
            logAux("Fatal", file, lineno, format, a...);
        }
    }
}

func logAux(level string, file string, lineno int, format string, a ...interface{}) {
    t := time.Now();
    a = append([]interface{} { level, t.Format("2006-01-02 15:04:05.00"), file, lineno }, a...);
    fmt.Printf("[%s] %s [%s:%d] " + format + "\n", a...);
}

