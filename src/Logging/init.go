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
    NOTSET LogLevel  = 0
    DEBUG LogLevel   = 1
    INFO LogLevel    = 2
    WARNING LogLevel = 3
    FATAL LogLevel   = 4
)

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
            logAux("WARNING", file, lineno, format, a...);
        }
    }
}

func Fatal(format string, a ...interface{}) {
    if logLevel <= FATAL {
        _, path, lineno, ok := runtime.Caller(1);
        _, file := filepath.Split(path)


        if (ok) {
            logAux("FATAL", file, lineno, format, a...);
        }
    }
}

func logAux(level string, file string, lineno int, format string, a ...interface{}) {
    t := time.Now();
    a = append([]interface{} { level, t.Format("2006-01-02 15:04:05.00"), file, lineno }, a...);
    fmt.Printf("[%s] %s [%s:%d] " + format + "\n", a...);
}

