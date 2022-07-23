package b

type (
    Data struct {
        AccountId  string                 `json:"accountId"`
        Event      string                 `json:"event"`
        Properties map[string]interface{} `json:"properties"`
    }
    DataOption func(d *Data)
)

func WithEvent(e string) DataOption {
    return func(d *Data) {
        d.Event = e
    }
}

func WithField(k string, v interface{}) DataOption {
    return func(d *Data) {
        if d.Properties == nil {
            d.Properties = map[string]interface{}{}
        }
        d.Properties[k] = v
    }
}
