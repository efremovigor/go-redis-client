package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"strings"
)

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
