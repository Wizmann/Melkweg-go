package Melkweg

import (
    "sync"
)

type Config struct {
    ProxyConfigs []ProxyConfig
}

type ProxyConfig struct {
    Key string
    Timeout int
    HeartbeatTimeout int

    Name string

    ServerAddr string
    ServerPort int

    ClientPort int

    ClientOutgoingConnectionNum int
}

var mu sync.Mutex
var configInstance *Config

func newConfig() *Config {
    return &Config {}
}

func (self *ProxyConfig) GetKey() string {
    return self.Key
}

func (self *ProxyConfig) GetTimeout() int {
    return self.Timeout
}

func (self *ProxyConfig) GetName() string {
    return self.Name
}

func (self *ProxyConfig) GetHeartbeatTimeout() int {
    return self.HeartbeatTimeout
}

func (self *ProxyConfig) GetServerAddr() string {
    return self.ServerAddr
}

func (self *ProxyConfig) GetServerPort() int {
    return self.ServerPort
}

func (self *ProxyConfig) GetClientPort() int {
    return self.ClientPort
}

func (self *ProxyConfig) GetClientOutgoingConnectionNum() int {
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
