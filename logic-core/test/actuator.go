package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Actuator struct {
	Nid    int `json:"nid"`	// node id
	Aid    int `json:"aid"` // actuator id
	Values []struct {		// action values
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
}

func main() {
	for {
		res := Actuator{
			Nid: 1,
			Aid: 1,
			Values: []struct{Value int "json:\"value\""; Sleep int "json:\"sleep\""}{},
		}
		res.Values = append(res.Values, struct{Value int "json:\"value\""; Sleep int "json:\"sleep\""}{1, 1})

		pbytes, _ := json.Marshal(res)
		buff := bytes.NewBuffer(pbytes)
		addr := "192.168.0.18:5000"
		log.Println("in Act.Exec, 받는 주소: " + "http://" + addr + "/actuator" + " 전달내용: " + string(pbytes))
		_, err := http.Post("http://"+addr+"/actuator", "application/json", buff)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}
}
