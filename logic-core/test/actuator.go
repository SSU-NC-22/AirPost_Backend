package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// type Drone struct {
// 	Nid    string `json:"nid"`	// node id
// 	Values struct {		// action values
// 		WP0 [][]float64 `json:"wp0"` // weigh point 0 (start station -> dest tag path)
// 		WP1 [][]float64 `json:"wp1"` // weigh point 1 (dest tag -> nearby station path)
// 	} `json:"values"`
// }

type Drone struct {
	Nid    string `json:"nid"`	// node id
	Values [][]float64 `json:"values"`
	Tagidx int `json:"tagidx"` // values 내에서 tag가 몇번째 index인지 (0~)
}

func main() {
	for {
		res := Drone{
			Nid: "DRO0",
			Values: [][]float64{},
			Tagidx: 1,
		}

		srcStation := []float64{37.497365670723944, 126.95591909983503, 1, 22} // lat, lon, alt, 명령(22: 이륙)
		tag := []float64{37.49719755738831, 126.95590032437323, 1, 16} // lat, lon, alt, 명령(16: 이동)
		destStation := []float64{37.4971933013496, 126.95619804955307, 0, 21} // lat, lon, alt, 명령(21: 착륙)

		// res.Values.WP0 = append(res.Values.WP0, srcStation)
		// res.Values.WP0 = append(res.Values.WP0, tag)
		// res.Values.WP1 = append(res.Values.WP1, tag)
		// res.Values.WP1 = append(res.Values.WP1, destStation)

		res.Values = append(res.Values, srcStation)
		res.Values = append(res.Values, tag)
		res.Values = append(res.Values, destStation)

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
