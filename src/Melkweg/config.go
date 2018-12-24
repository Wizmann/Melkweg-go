package Melkweg

import (
    "sync"
)

type Config struct {
    key string
    timeout int

    serverAddr string
    serverPort int

    clientOutgoingConnectionNum int
}

var mu sync.Mutex
var configInstance *Config

// TODO: read config from file
func newConfig() *Config {
    return &Config {}
}

func (self *Config) GetKey() string {
    return self.key
}

func (self *Config) GetTimeout() int {
    return self.timeout
}

func (self *Config) GetServerAddr() string {
    return self.serverAddr
}

func (self *Config) GetServerPort() int {
    return self.serverPort
}

func (self *Config) GetClientOutgoingConnectionNum() int {
    return self.clientOutgoingConnectionNum
}

func GetConfigInstance() *Config {
    mu.Lock()
    defer mu.Unlock()

    if (configInstance == nil) {
        configInstance = newConfig()
    }
    return configInstance
}
