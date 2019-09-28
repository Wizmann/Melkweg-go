package Gwisted
import (
    "fmt"
    "net"
    logging "Logging"
)

type IConnector interface {
    Connect() (*net.TCPConn, error)
}

type ClientConnector struct {
    IConnector

    host string
    port int
    timeout int
}

func NewClientConnector(host string, port int, timeout int) *ClientConnector {
    return &ClientConnector {
        host: host,
        port: port,
        timeout: timeout,
    }
}

func (self *ClientConnector) Connect() (*net.TCPConn, error) {
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", self.host, self.port))
    if (err != nil) {
        logging.Fatal("dial TCP error on ", self.host, ":", self.port)
        return nil, err
    }
    return conn.(*net.TCPConn), nil
}

