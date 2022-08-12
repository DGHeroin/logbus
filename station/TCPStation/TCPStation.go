package main

import (
    "bufio"
    "encoding/json"
    "github.com/DGHeroin/logbus/bus"
    "io"
    "log"
    "net"
)

func main() {
    ln, err := net.Listen("tcp", ":50002")
    if err != nil {
        panic(err)
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            break
        }
        go func(conn net.Conn) {
            buf := bufio.NewReader(conn)
            dec := json.NewDecoder(buf)
            defer conn.Close()
            for {
                var buffer []bus.Data
                if err := dec.Decode(&buffer); err != nil {
                    if err == io.EOF {
                        break
                    }
                    log.Println(err)
                    break
                }
                for _, data := range buffer {
                    log.Println(data)
                }
            }
        }(conn)
    }
}
