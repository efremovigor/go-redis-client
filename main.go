package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"github.com/go-redis/redis"
)

func rootHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	//par := request.URL.Query()["qwe"];
	//data := "{\"response\":{\"code\":200,\"status\":\"OK\",\"par\":"+par[0]+"}}"

	request.ParseForm()
	reddis_string := request.Form.Get("reddis_string");




	data := "{\"response\":{\"code\":200,\"status\":\"OK\",\"par\":"+ExampleNewClient(reddis_string)[0]+"}}"

	response.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(response, string(data))
}

func ExampleNewClient(reddis_string string) []string {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, _ := client.Keys("*").Result()
	return pong
}


func main() {
	http.HandleFunc("/", rootHandler)
	go func() {
		resp, err := http.PostForm("http://127.0.0.1:8080/?qwe=2222", url.Values{"reddis_string": {"GET KEYS *"}})
		if err != nil{
			fmt.Print(err)
		}
		responseData,_ := ioutil.ReadAll(resp.Body)
		responseString := string(responseData)
		fmt.Println(responseString)
	}()


	log.Fatal(http.ListenAndServe(":8080", nil))


}