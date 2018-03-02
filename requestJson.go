package main

import "time"

type requestJson struct {
    Command struct {
        Method string        `json:"method"`
        Key    string        `json:"key"`
        Value  string        `json:"value"`
        Time   time.Duration `json:"time"`
    } `json:"command"`
}

func createRequestJson() requestJson {
    requestJson := requestJson{}
    requestJson.Command.Time = 0
    return requestJson
}
