package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"github.com/go-redis/redis"
	"encoding/json"
	"time"
)

type responseJson struct {
	Code      int      `json:"code"`
	Status    string   `json:"status"`
	Data      []string `json:"data"`
	DebugInfo []string `json:"debugInfo"`
}

type requestJson struct {
	Command struct {
		Method string        `json:"method"`
		Key    string        `json:"key"`
		Value  string        `json:"value"`
		Time   time.Duration `json:"time"`
	} `json:"command"`
}

type redisConnection struct {
	Host string
	Port string
	Password string
	DB int
	Request chan *requestJson
}

var redisMappingDevList = map[string]string{
	"10.20.29.11:6379":  "95.213.208.34:6380",
	"10.20.29.143:6379": "95.213.208.34:6381",
	"10.20.30.11:6379":  "95.213.208.34:6382",
	"10.20.28.10:6379":  "95.213.208.34:6383",
	"10.20.28.11:6379":  "95.213.208.34:6384",
	"10.20.28.12:6379":  "95.213.208.34:6385",
}

var redisConn1 = createRedisConnection("95.213.208.34","6380","")
var redisConn2 = createRedisConnection("95.213.208.34","6380","")
var redisConn3 = createRedisConnection("95.213.208.34","6380","")
var redisConn4 = createRedisConnection("95.213.208.34","6380","")
var redisConn5 = createRedisConnection("95.213.208.34","6380","")
var redisConn6 = createRedisConnection("95.213.208.34","6380","")

func createRequestJson() requestJson {
	requestJson := requestJson{}
	requestJson.Command.Time = 0
	return requestJson
}

func createRedisConnection(host,port,password string) redisConnection {
	redisConnection := redisConnection{}
	redisConnection.Host = host
	redisConnection.Port = port
	redisConnection.Request = make(chan *requestJson, 1)
	redisConnection.DB = 0
	redisConnection.Password = password
	return redisConnection
}

func rootHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	var requestJson = createRequestJson()
	json.Unmarshal([]byte(request.Form.Get("json")), &requestJson)
	json.NewDecoder(request.Body).Decode(&requestJson)


	select {
	case redisConn1.Request <-&requestJson:
		chanProcess(response, redisConn1 )
	case redisConn2.Request <-&requestJson:
		chanProcess(response, redisConn2)
	case redisConn3.Request <-&requestJson:
		chanProcess(response, redisConn3)
	case redisConn4.Request <-&requestJson:
		chanProcess(response, redisConn4)
	case redisConn5.Request <-&requestJson:
		chanProcess(response, redisConn5 )
	case redisConn6.Request <-&requestJson:
		chanProcess(response, redisConn6)
	}
}

func chanProcess(response http.ResponseWriter, con redisConnection) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	jsonData, _ := json.Marshal(responseJson{
		Code:      http.StatusOK,
		Status:    http.StatusText(http.StatusOK),
		Data:      con.redisProcess(<-con.Request),
		DebugInfo: []string{con.Host + ":"+con.Port},
	})

	response.Header().Set("Content-Length", fmt.Sprint(len(string(jsonData))))
	fmt.Fprint(response, string(jsonData))
}

func (con redisConnection) redisProcess(requestJson *requestJson) (list []string) {
	var response string
	var err error
	client := redis.NewClient(&redis.Options{
		Addr:     con.Host + ":"+con.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	switch requestJson.Command.Method {
	case "GET":
		response, err = client.Get(requestJson.Command.Key).Result()
	case "SET":
		response, err = client.Set(requestJson.Command.Key, requestJson.Command.Value, requestJson.Command.Time*time.Second).Result()
	}
	fmt.Println(err)
	list = append(list, response)
	return
}

func main() {
	http.HandleFunc("/", rootHandler)
		go func() {
			for i := 0; i < 1000; i++ {

			resp, _ := http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"GET\",\"key\":\"foo"+string(i)+"\"}}"}})
			responseData, _ := ioutil.ReadAll(resp.Body)
			responseString := string(responseData)
			fmt.Println(responseString)

			//resp, _ = http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"SET\",\"key\":\"foo27802\",\"value\":\"супер-креветка\",\"time\":0}}"}})

			//resp, _ = http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"GET\",\"key\":\"foo27802\"}}"}})
			//responseData, _ = ioutil.ReadAll(resp.Body)
			//responseString = string(responseData)
			//fmt.Println(responseString)
			}

		}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
