package main

import (
    "encoding/json"
    "github.com/DGHeroin/logbus/bus"
    "github.com/DGHeroin/logbus/utils"
    "io/ioutil"
    "log"
    "net/http"
)

func JSONResponse(w http.ResponseWriter, code int, r map[string]interface{}) {
    data, _ := json.Marshal(&r)
    w.WriteHeader(code)
    w.Write(data)
}

func main() {
    http.HandleFunc("/logbus", func(w http.ResponseWriter, r *http.Request) {
        var (
            appId    = r.Header.Get("appId")
            compress = r.Header.Get("compress")
            payload  string
        )
        log.Println(r.Header)
        // check app id
        //
        switch compress {
        case "none":
            data, err := ioutil.ReadAll(r.Body)
            if err != nil {
                JSONResponse(w, 200, map[string]interface{}{
                    "code": 1001,
                    "err":  "read payload",
                })
                return
            }
            payload = string(data)
        case "gzip":
            data, err := utils.GZipDecodeData(r.Body)
            if err != nil {
                JSONResponse(w, 200, map[string]interface{}{
                    "code": 1001,
                    "err":  "read gzip payload",
                })
                return
            }
            payload = data
        }
        // decode payload
        var buffer []bus.Data
        if err := json.Unmarshal([]byte(payload), &buffer); err != nil {
            JSONResponse(w, 200, map[string]interface{}{
                "code": 1001,
                "err":  "decode payload error",
            })
            return
        }
        for _, data := range buffer {
            log.Println("得到", appId, data)
        }

        JSONResponse(w, 200, map[string]interface{}{
            "code": 0,
        })
    })
    http.ListenAndServe(":12345", nil)
}
