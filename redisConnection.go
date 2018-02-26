package main

import (
    "github.com/go-redis/redis"
    "time"
)

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

