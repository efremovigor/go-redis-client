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
var redisList = map[string]string{
	"10.20.29.11:6379" : "95.213.208.34:6380",
		"10.20.29.143:6379" : "95.213.208.34:6381",
		"10.20.30.11:6379" : "95.213.208.34:6382",
		"10.20.28.10:6379" : "95.213.208.34:6383",
		"10.20.28.11:6379" : "95.213.208.34:6384",
		"10.20.28.12:6379" : "95.213.208.34:6385",
}

type responseJson struct {
	Code   int      `json:"code"`
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

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

func rootHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	//par := request.URL.Query()["qwe"];
	//data := "{\"response\":{\"code\":200,\"status\":\"OK\",\"par\":"+par[0]+"}}"

	request.ParseForm()
	var requestJson = createRequestJson()
	json.Unmarshal([]byte(request.Form.Get("json")), &requestJson)
	json.NewDecoder(request.Body).Decode(&requestJson)
	jsonData, _ := json.Marshal(responseJson{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Data:   redisProcess(&requestJson),
	})

	response.Header().Set("Content-Length", fmt.Sprint(len(string(jsonData))))
	fmt.Fprint(response, string(jsonData))
}

func redisProcess(requestJson *requestJson) (list []string) {
	var response string
	client := redis.NewClient(&redis.Options{
		Addr:     "95.213.208.34:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	switch requestJson.Command.Method {
	case "GET":
		response, _ = client.Get(requestJson.Command.Key).Result()
	case "SET":
		response, _ = client.Set(requestJson.Command.Key, requestJson.Command.Value, requestJson.Command.Time*time.Second).Result()
	}
	list = append(list, response)
	return
}

func main() {
	http.HandleFunc("/", rootHandler)
	go func() {

		resp, _ := http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"GET\",\"key\":\"foo27802\"}}"}})
		responseData, _ := ioutil.ReadAll(resp.Body)
		responseString := string(responseData)
		fmt.Println(responseString)

		resp, _ = http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"SET\",\"key\":\"foo27802\",\"value\":\"супер-креветка\",\"time\":1}}"}})

		resp, _ = http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"GET\",\"key\":\"foo27802\"}}"}})
		responseData, _ = ioutil.ReadAll(resp.Body)
		responseString = string(responseData)
		fmt.Println(responseString)

		time.Sleep(1 * time.Second)

		resp, _ = http.PostForm("http://127.0.0.1:8080", url.Values{"json": {"{\"command\":{\"method\":\"GET\",\"key\":\"foo27802\"}}"}})
		responseData, _ = ioutil.ReadAll(resp.Body)
		responseString = string(responseData)
		fmt.Println(responseString)
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
