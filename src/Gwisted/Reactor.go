package Gwisted

type Reactor struct {
    ctrlCh <- chan int
}

func (self *Reactor) Start() {
    for {
        select {
        case <- self.ctrlCh:
            break
        }
    }
}

func (self *Reactor) Stop() {
    self.ctrlCh <- false
}


func (self *Reactor) ListenTCP(port int, factory *ProtocolFactory, backlog int) {
    // pass
}

func (self *Reactor) ConnectTCP(host string, port int, factory *ProtocolFactory, timeout int) {
    // pass
}
