package main

import (
    "encoding/json"
    logging "github.com/op/go-logging"
    . "Melkweg"
    "os"
)

var log = logging.MustGetLogger("Melkweg.Client")

func init() {
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    backendFormatter := logging.NewBackendFormatter(backend, logging.GlogFormatter)
    logging.SetBackend(backendFormatter)

    file, err := os.Open("config.json")
    if (err != nil) {
        panic("can't open config.json file")
    }
    defer file.Close()
    decoder := json.NewDecoder(file)
    config := GetConfigInstance()
    err = decoder.Decode(config)
    if (err != nil) {
        panic(err)
    }
}

