package bus

import (
    "github.com/DGHeroin/logbus/b"
    "time"
)

type (
    Car interface {
        Go(data []b.Data) error
        WaitFinish(timeout time.Duration) error
    }
)
