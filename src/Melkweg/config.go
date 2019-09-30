package Melkweg

import (
    "sync"
)

type Config struct {
    Key string
    Timeout int
    HeartbeatTimeout int

    ServerAddr string
    ServerPort int

    ClinetPort int

    ClientOutgoingConnectionNum int
}

var mu sync.Mutex
var configInstance *Config

func newConfig() *Config {
    return &Config {}
}

func (self *Config) GetKey() string {
    return self.Key
}

func (self *Config) GetTimeout() int {
    return self.Timeout
}

func (self *Config) GetHeartbeatTimeout() int {
    return self.HeartbeatTimeout
}

func (self *Config) GetServerAddr() string {
    return self.ServerAddr
}

func (self *Config) GetServerPort() int {
    return self.ServerPort
}

func (self *Config) GetClientPort() int {
    return self.ClinetPort
}

func (self *Config) GetClientOutgoingConnectionNum() int {
    return self.ClientOutgoingConnectionNum
}

func GetConfigInstance() *Config {
    mu.Lock()
    defer mu.Unlock()

    if (configInstance == nil) {
        configInstance = newConfig()
    }
    return configInstance
}
