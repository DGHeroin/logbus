package httpcar

import (
    "bytes"
    "encoding/json"
    "github.com/DGHeroin/logbus/bus"
    "github.com/DGHeroin/logbus/utils"
    "io/ioutil"
    "net/http"
    "strconv"
    "sync"
    "time"
)

type (
    HTTPCar struct {
        wg   sync.WaitGroup
        opts *options
    }
    options struct {
        serverUrl string
        compress  bool
        appId     string
        timeout   time.Duration
    }
    Options func(*options)
)

func (c *HTTPCar) WaitFinish(timeout time.Duration) error {
    return utils.WaitGroupWithout(&c.wg, timeout)
}

// 发车
func (c *HTTPCar) Go(buffer []bus.Data) error {
    c.wg.Add(1)
    defer c.wg.Done()

    jdata, err := json.Marshal(buffer)
    if err != nil {
        return err
    }
    if _, _, err = c.send(string(jdata), len(buffer)); err != nil {
        return err
    }
    return nil
}
func (c *HTTPCar) send(data string, size int) (statusCode int, code int, err error) {
    var encodedData string
    var compressType = "gzip"
    if c.opts.compress {
        encodedData, err = utils.GZipEncodeData(data)
    } else {
        encodedData = data
        compressType = "none"
    }
    if err != nil {
        return 0, 0, err
    }
    postData := bytes.NewBufferString(encodedData)

    var resp *http.Response
    req, _ := http.NewRequest("POST", c.opts.serverUrl, postData)
    req.Header["appid"] = []string{c.opts.appId}
    req.Header.Set("user-agent", "LogBus")
    req.Header.Set("compress", compressType)
    req.Header["data-count"] = []string{strconv.Itoa(size)}
    client := &http.Client{Timeout: c.opts.timeout}
    resp, err = client.Do(req)

    if err != nil {
        return 0, 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        var result struct {
            Code int
        }

        err = json.Unmarshal(body, &result)
        if err != nil {
            return resp.StatusCode, 1, err
        }

        return resp.StatusCode, result.Code, nil
    } else {
        return resp.StatusCode, -1, nil
    }
}
func New(opts ...Options) bus.Car {
    o := defaultOption()
    for _, opt := range opts {
        opt(o)
    }

    return &HTTPCar{
        opts: o,
    }
}

func defaultOption() *options {
    return &options{
        serverUrl: "http://127.0.0.1:12345/logbus",
        compress:  false,
        appId:     "test-app",
        timeout:   time.Second * 60,
    }
}
func WithServerUrl(u string) Options {
    return func(o *options) {
        o.serverUrl = u
    }
}
