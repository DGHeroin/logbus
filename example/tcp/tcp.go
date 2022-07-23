package main

import (
    "github.com/DGHeroin/logbus/b"
    "github.com/DGHeroin/logbus/bus"
    "github.com/DGHeroin/logbus/bus/tcpcar"
    "github.com/DGHeroin/logbus/utils"
)

func genData() b.Data {
    return b.Data{
        AccountId: "124",
        Event:     "普通",
        Properties: map[string]interface{}{
            "code": 11,
            "name": 66,
            "time": utils.GetTimeString(),
        },
    }
}
func main() {
    car := tcpcar.New("127.0.0.1:50002")
    d := bus.NewDriver(car)
    for i := 0; i < 200; i++ {
        d.Add(genData())
    }
    d.Adds(b.WithEvent("我的事件"), b.WithField("年纪", 24))
    d.Close()
}
