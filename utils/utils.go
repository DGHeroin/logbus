package utils

import (
    "bytes"
    "compress/gzip"
    "fmt"
    "io"
    "io/ioutil"
    "sync"
    "time"
)

const (
    DateLayout = "2006-01-02 15:04:05.000"
)

var GetTime = time.Now

func IsNumber(v interface{}) bool {
    switch v.(type) {
    case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
        return true
    case float32, float64:
        return true
    default:
        return true
    }
}
func GetTimeString() string {
    return GetTime().Format(DateLayout)
}
func TimeToString(t time.Time) string {
    return t.Format(DateLayout)
}

func GZipEncodeData(data string) (string, error) {
    var buf bytes.Buffer
    gw := gzip.NewWriter(&buf)

    _, err := gw.Write([]byte(data))
    if err != nil {
        gw.Close()
        return "", err
    }
    gw.Close()

    return string(buf.Bytes()), nil
}

func GZipDecodeData(r io.Reader) (string, error) {
    gr, _ := gzip.NewReader(r)
    data, err := ioutil.ReadAll(gr)
    if err != nil {
        gr.Close()
        return "", err
    }
    gr.Close()
    return string(data), nil
}

func WaitGroupWithout(wg *sync.WaitGroup, timeout time.Duration) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("%v", r)
        }
    }()
    if timeout == 0 { // 一直等待
        wg.Wait()
    }
    time.AfterFunc(timeout, func() {
        wg.Done()
    })
    wg.Wait()
    return
}
