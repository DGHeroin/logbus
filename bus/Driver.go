package bus

import (
    "github.com/DGHeroin/logbus/utils"
    "sync"
    "time"
)

type (
    Driver interface {
        Add(data Data)
        Adds(...DataOption)
        Flush() error
        Close() error
    }
    driver struct {
        car         Car
        bufferMutex sync.RWMutex
        buffer      []Data
        cacheMutex  sync.RWMutex
        cacheBuffer [][]Data
        option      *options
    }
    options struct {
        batchSize     int
        cacheCapacity int
        interval      time.Duration
        retries       int
    }
    Options func(o *options)
)

func (d *driver) Adds(opts ...DataOption) {
    data := Data{}
    for _, opt := range opts {
        opt(&data)
    }
    d.Add(data)
}

func (d *driver) Add(data Data) {
    if data.Properties != nil {
        if _, ok := data.Properties["time"]; !ok {
            data.Properties["time"] = utils.GetTimeString()
        }
    }
    d.bufferMutex.Lock()
    d.buffer = append(d.buffer, data)
    d.bufferMutex.Unlock()
    d.checkFlush()
}
func (d *driver) Flush() error {
    d.growCacheBuffer(false)

    d.bufferMutex.Lock()
    defer d.bufferMutex.Unlock()

    if len(d.buffer) == 0 {
        return nil
    }

    sendBuffer := d.buffer
    d.buffer = d.cacheBuffer[0]
    d.cacheBuffer = d.cacheBuffer[1:]

    for i := 0; i <= d.option.retries; i++ {
        if err := d.car.Go(sendBuffer); err == nil {
            break
        }
    }

    return nil
}

func (d *driver) Close() error {
    return d.Flush()
}

func (d *driver) checkFlush() {
    if d.getBufferSize() >= d.option.batchSize {
        d.Flush()
    }
}

func (d *driver) getBufferSize() int {
    d.bufferMutex.RLock()
    defer d.bufferMutex.RUnlock()
    return len(d.buffer)
}

func (d *driver) getCacheBufferSize() int {
    d.bufferMutex.RLock()
    defer d.bufferMutex.RUnlock()
    return len(d.cacheBuffer)
}
func (d *driver) growCacheBuffer(fillToFull bool) {
    d.cacheMutex.Lock()
    defer d.cacheMutex.Unlock()
    sz := len(d.cacheBuffer)
    if !fillToFull {
        if sz > 0 {
            return
        }
    }
    for i := sz; i < d.option.cacheCapacity; i++ {
        d.cacheBuffer = append(d.cacheBuffer, make([]Data, 0, d.option.batchSize))
    }
}
func NewDriver(car Car, opts ...Options) Driver {
    o := defaultOption()
    for _, opt := range opts {
        opt(o)
    }
    d := &driver{
        car:    car,
        option: o,
    }
    // 初始化 buffer
    d.buffer = make([]Data, 0, d.option.batchSize)
    d.growCacheBuffer(true)
    if d.option.interval > 0 {
        go func() { // auto flush
            t := time.NewTicker(d.option.interval)
            defer t.Stop()
            for {
                <-t.C
                _ = d.Flush()
            }
        }()
    }
    return d
}

func defaultOption() *options {
    return &options{
        batchSize:     100,
        cacheCapacity: 5,
        interval:      0,
        retries:       3,
    }
}
func WithBatchSize(n int) Options {
    return func(o *options) {
        o.batchSize = n
    }
}
func WithCacheCapacity(n int) Options {
    return func(o *options) {
        o.cacheCapacity = n
    }
}
