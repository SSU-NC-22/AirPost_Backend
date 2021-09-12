package main

import (
	"log"

	"github.com/go-resty/resty/v2"
)


func main() {
	startStation := "126.95591909983503,37.497365670723944"
	tag := "126.95590032437323,37.49719755738831"
	destStation := "126.95619804955307,37.4971933013496"

	client := resty.New()

	resp1, err := client.R().
		SetQueryString("start=" + startStation + "&goal=" + tag).
		SetHeader("X-NCP-APIGW-API-KEY-ID", "6a14n8xual").
		SetHeader("X-NCP-APIGW-API-KEY", "vej8eUozJVRvtrdCZcTlV4ea9ljEriJUxdEa7j42").
		Get("https://naveropenapi.apigw.ntruss.com/map-direction/v1/driving")
	if err != nil || !resp1.IsSuccess() {
		log.Println("resp 1")
		panic(err)
	}

	resp1.Body()

	resp2, err := client.R().
		SetQueryString("start=" + tag + "&goal=" + destStation).
		SetHeader("X-NCP-APIGW-API-KEY-ID", "6a14n8xual").
		SetHeader("X-NCP-APIGW-API-KEY", "vej8eUozJVRvtrdCZcTlV4ea9ljEriJUxdEa7j42").
		Get("https://naveropenapi.apigw.ntruss.com/map-direction/v1/driving")
	if err != nil || !resp2.IsSuccess() {
		log.Println("resp 2")
		panic(err)
	}
}
