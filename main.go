package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"time"
	"strings"
)

const DEV  = true

type responseJson struct {
	Code      int      `json:"code"`
	Status    string   `json:"status"`
	Data      []string `json:"data"`
	DebugInfo []string `json:"debugInfo"`
}

func (res *responseJson) addDebugMsg (msg string)  {
	res.DebugInfo = append(res.DebugInfo,msg)
}

type requestJson struct {
	Command struct {
		Method string        `json:"method"`
		Key    string        `json:"key"`
		Value  string        `json:"value"`
		Time   time.Duration `json:"time"`
	} `json:"command"`
}

var redisMappingDev = map[string]string{
	"10.20.29.11:6379":  "95.213.208.34:6380",
	"10.20.29.143:6379": "95.213.208.34:6381",
	"10.20.30.11:6379":  "95.213.208.34:6382",
	"10.20.28.10:6379":  "95.213.208.34:6383",
	"10.20.28.11:6379":  "95.213.208.34:6384",
	"10.20.28.12:6379":  "95.213.208.34:6385",
}

var redisConnMap = map[string]redisConnection{
	"95.213.208.34:6380" : createRedisConnection("95.213.208.34","6380",""),
	"95.213.208.34:6381" : createRedisConnection("95.213.208.34","6381",""),
	"95.213.208.34:6382" : createRedisConnection("95.213.208.34","6382",""),
	"95.213.208.34:6383" : createRedisConnection("95.213.208.34","6383",""),
	"95.213.208.34:6384" : createRedisConnection("95.213.208.34","6384",""),
	"95.213.208.34:6385" : createRedisConnection("95.213.208.34","6385",""),
}

func createRequestJson() requestJson {
	requestJson := requestJson{}
	requestJson.Command.Time = 0
	return requestJson
}

func rootHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	requestJson := createRequestJson()
	json.Unmarshal([]byte(request.Form.Get("json")), &requestJson)
	json.NewDecoder(request.Body).Decode(&requestJson)

	for _ ,conn := range redisConnMap {
		select {
		case conn.Request <-&requestJson:
			chanProcess(response, conn)
			return
		}
	}
}

func chanProcess(response http.ResponseWriter, con redisConnection) {
	var responseJson responseJson
	cache := <- con.Request
	redisOutput, err := con.redisProcess(cache)
	if err != nil {
		responseJson.addDebugMsg(err.Error())
		redisOutputSlice := strings.Split(err.Error(), " ")
		if redisOutputSlice[0] == "MOVED" {
			redirectConn := redisConnMap[redisMappingDev[redisOutputSlice[2]]]
			redirectConn = createRedisConnection(redirectConn.Host, redirectConn.Port, redirectConn.Password)
			redisOutput, err = redirectConn.redisProcess(cache)
			if err != nil {
				responseJson.addDebugMsg(err.Error())
			}
		}
	}
	responseJson.Data = redisOutput
	responseJson.Code = http.StatusOK
	responseJson.Status = http.StatusText(http.StatusOK)
	jsonData, _ := json.Marshal(responseJson)

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Length", fmt.Sprint(len(string(jsonData))))
	fmt.Fprint(response, string(jsonData))
}

func main() {
	http.HandleFunc("/", rootHandler)
		go func() {
			for i := 27800; i < 30000; i++ {

			resp, _ := http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"GET\",\"key\":\"foo"+string(i)+"\"}}"}})
			responseData, _ := ioutil.ReadAll(resp.Body)
			responseString := string(responseData)
			fmt.Println(responseString)

			//resp, _ = http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"SET\",\"key\":\"foo"+string(i)+"\",\"value\":\"супер-креветка\",\"time\":0}}"}})
			//
			//resp, _ = http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"GET\",\"key\":\"foo27802\"}}"}})
			//responseData, _ = ioutil.ReadAll(resp.Body)
			//responseString = string(responseData)
			//fmt.Println(responseString)
			}

		}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
