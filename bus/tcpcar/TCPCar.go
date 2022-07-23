package tcpcar

import (
    "encoding/json"
    "github.com/DGHeroin/logbus/b"
    "github.com/DGHeroin/logbus/bus"
    "github.com/DGHeroin/logbus/utils"
    "net"
    "sync"
    "time"
)

type (
    TCPCar struct {
        mu     sync.Mutex
        enc    *json.Encoder
        wg     sync.WaitGroup
        option options
    }
    options struct {
        address  string
        compress bool
        timeout  time.Duration
    }
    Options func(*options)
)

func (t *TCPCar) Go(buffer []b.Data) error {
    t.wg.Add(1)
    defer t.wg.Done()
    t.mu.Lock()

    if t.enc == nil {
        conn, err := net.Dial("tcp", t.option.address)
        if err != nil {
            t.mu.Unlock()
            return err
        }
        t.enc = json.NewEncoder(conn)
    }
    t.mu.Unlock()

    err := t.enc.Encode(&buffer)
    return err
}

func (t *TCPCar) WaitFinish(timeout time.Duration) error {
    return utils.WaitGroupWithout(&t.wg, timeout)
}

func New(address string, opts ...Options) bus.Car {
    o := defaultOption()
    o.address = address
    for _, opt := range opts {
        opt(&o)
    }

    return &TCPCar{
        option: o,
    }
}

func defaultOption() options {
    return options{
        address:  "",
        compress: false,
        timeout:  time.Second * 60,
    }
}
