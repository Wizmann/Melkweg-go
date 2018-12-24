package Melkweg

import (
    logging "github.com/op/go-logging"
    "os"
)

var log = logging.MustGetLogger("Melkweg.Client")

func init() {
    backend := logging.NewLogBackend(os.Stderr, "", 0)
    backendFormatter := logging.NewBackendFormatter(backend, logging.GlogFormatter)
    logging.SetBackend(backendFormatter)
}

