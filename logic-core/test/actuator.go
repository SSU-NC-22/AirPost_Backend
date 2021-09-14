package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Actuator struct {
	Nid    string `json:"nid"`	// node id
	Aid    string `json:"aid"` // actuator id
	Values []struct {		// action values
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
}

type Drone struct {
	Nid    string `json:"nid"`	// node id
	Values struct {		// action values
		WP0 [][]float64 `json:"wp0"` // weigh point 0 (start station -> dest tag path)
		WP1 [][]float64 `json:"wp1"` // weigh point 1 (dest tag -> nearby station path)
	} `json:"values"`
}

func main() {
	for {
		// res := Actuator{
		// 	Nid: 1,
		// 	Aid: 1,
		// 	Values: []struct{Value int "json:\"value\""; Sleep int "json:\"sleep\""}{},
		// }
		// res.Values = append(res.Values, struct{Value int "json:\"value\""; Sleep int "json:\"sleep\""}{1, 1})

		res := Drone{
			Nid: "DRO0",
			Values: struct{WP0 [][]float64 "json:\"wp0\""; WP1 [][]float64 "json:\"wp1\""}{},
		}

		srcStation := []float64{37.497365670723944, 126.95591909983503}
		tag := []float64{37.49719755738831, 126.95590032437323}
		destStation := []float64{37.4971933013496, 126.95619804955307}

		res.Values.WP0 = append(res.Values.WP0, srcStation)
		res.Values.WP0 = append(res.Values.WP0, tag)
		res.Values.WP1 = append(res.Values.WP1, tag)
		res.Values.WP1 = append(res.Values.WP1, destStation)

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
