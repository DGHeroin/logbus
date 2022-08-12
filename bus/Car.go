package bus

import (
    "time"
)

type (
    Car interface {
        Go(data []Data) error
        WaitFinish(timeout time.Duration) error
    }
)
