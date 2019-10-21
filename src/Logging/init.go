package Gwisted

import (
    "runtime"
    "path/filepath"
    "fmt"
    "time"
    "strings"
    "strconv"
    "os"
)

var LogLevel = INFO

type ELogLevel int

const (
    NOTSET   ELogLevel = 0
    VERBOSE  ELogLevel = 1
    DEBUG    ELogLevel = 2
    INFO     ELogLevel = 3
    WARNING  ELogLevel = 4
    ERROR    ELogLevel = 5
    FATAL    ELogLevel = 6
)

func init() {
    envVar := os.Getenv("Melkweg-Logging")
    envVar = strings.ToUpper(envVar)

    if (envVar == "NOTSET") {
        LogLevel = 0
    } else if (envVar == "VERBOSE") {
        LogLevel = 1
    } else if (envVar == "DEBUG") {
        LogLevel = 2
    } else if (envVar == "INFO") {
        LogLevel = 3
    } else if (envVar == "WARNING") {
        LogLevel = 4
    } else if (envVar == "ERROR") {
        LogLevel = 5
    } else if (envVar == "FATAL") {
        LogLevel = 6
    } else if (envVar == "") {
        /* pass */
    } else {
        fmt.Printf("[%s] is not a valid log level")
    }
}

func GetGoroutineID() int {
    var buf [64]byte
    n := runtime.Stack(buf[:], false)
    idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
    id, err := strconv.Atoi(idField)
    if err != nil {
        panic(fmt.Sprintf("cannot get goroutine id: %v", err))
    }
    return id
}

func Verbose(format string, a ...interface{}) {
    if LogLevel <= VERBOSE {
        logAux("Verbose", format, a...)
    }
}

func Debug(format string, a ...interface{}) {
    if LogLevel <= DEBUG {
        logAux("Debug", format, a...)
    }
}

func Info(format string, a ...interface{}) {
    if LogLevel <= INFO {
        logAux("Info", format, a...)
    }
}

func Warning(format string, a ...interface{}) {
    if LogLevel <= WARNING {
        logAux("Warning", format, a...)
    }
}

func Error(format string, a ...interface{}) {
    if LogLevel <= ERROR {
        logAux("Error", format, a...)
    }
}

func Fatal(format string, a ...interface{}) {
    if LogLevel <= FATAL {
        logAux("Fatal", format, a...);
    }
}

func logAux(level string, format string, a ...interface{}) {
    _, path, lineno, _ := runtime.Caller(2);
    wd, _ := os.Getwd()
    relpath, _ := filepath.Rel(wd, path)

    t := time.Now();
    goid := GetGoroutineID()
    a = append([]interface{} { level, t.Format("2006-01-02 15:04:05.00"), goid, relpath, lineno }, a...);
    fmt.Printf("[%s] %s [T%d][%s:%d] " + format + "\n", a...);
}

