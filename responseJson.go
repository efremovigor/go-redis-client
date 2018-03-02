package main


type responseJson struct {
    Code      int      `json:"code"`
    Status    string   `json:"status"`
    Data      []string `json:"data"`
    DebugInfo []string `json:"debugInfo"`
}

func (res *responseJson) addDebugMsg (msg string)  {
    res.DebugInfo = append(res.DebugInfo,msg)
}

