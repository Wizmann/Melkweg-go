package Client

import (
    "encoding/json"
    . "Melkweg"
    "os"
)

func init() {
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
