package Melkweg

import (
    "sync"
)

var mu Sync.Mutex

type Config struct {
    key string
    timeout int
}

// TODO: read config from file
func newConfig() *Config {
    return &Config {}
}

func (self *Config) GetKey() string {
    return self.key
}

func (self *Config) GetTimeout() string {
    return self.timeout
}

func GetInstance() *singleton {
    mu.Lock()
    defer mu.Unlock()

    if (instance == nil) {
        instance = newConfig()
    }
    return instance
}
