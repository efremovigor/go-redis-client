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
	"strings"
)

const DEV  = true

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


func createRedisConnection(host,port,password string) redisConnection {
	redisConnection := redisConnection{}
	redisConnection.Host = host
	redisConnection.Port = port
	redisConnection.Request = make(chan *requestJson, 1)
	redisConnection.DB = 0
	redisConnection.Password = password
	return redisConnection
}


func (con *redisConnection) redisProcess(requestJson *requestJson) (list []string,err error) {
	var response string
	client := redis.NewClient(&redis.Options{
		Addr:     con.Host + ":"+con.Port,
		Password: "",
		DB:       0,
	})
	switch requestJson.Command.Method {
	case "GET":
		response, err = client.Get(requestJson.Command.Key).Result()
	case "SET":
		response, err = client.Set(requestJson.Command.Key, requestJson.Command.Value, requestJson.Command.Time*time.Second).Result()
	}
	list = append(list, response)
	return list, err
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
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	cache := <- con.Request
	redisOutput, err := con.redisProcess(cache)

	if err != nil {
		output := strings.Split(err.Error(), " ")
		if output[0] == "MOVED" {
			redirectConn := redisConnMap[redisMappingDev[output[2]]]
			redirectConn = createRedisConnection(redirectConn.Host, redirectConn.Port, redirectConn.Password)
			redisOutput, err = redirectConn.redisProcess(cache)
		}
	}

	jsonData, _ := json.Marshal(responseJson{
		Code:      http.StatusOK,
		Status:    http.StatusText(http.StatusOK),
		Data:      redisOutput,
		DebugInfo: []string{con.Host + ":"+con.Port},
	})

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
