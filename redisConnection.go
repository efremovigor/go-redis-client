package main

import (
    "github.com/go-redis/redis"
    "time"
)

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
        Addr:     con.Host + ":" + con.Port,
        Password: con.Password,
        DB:       con.DB,
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

