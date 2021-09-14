package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

const (
	from    = "airpost@gmail.com"
	pass    = "REDACTED"
	bodyFmt = "node(%s)"
	msgFmt  = "From: %s\nTo: %s\nSubject: AirPost email\n\n%s"
)

type EmailElement struct {
	BaseElement
	Email    string `json:"text"`
	Interval map[string]bool
}

func (ee *EmailElement) Exec(d *model.LogicData) {
	log.Println("\t\t!!!!in EmailElement.Exec !!!!")
	ok, exist := ee.Interval[d.Node.Name]

	if !exist {
		ee.Interval[d.Node.Name] = true
	}
	if ok {
		ee.Interval[d.Node.Name] = false

		body := fmt.Sprintf(bodyFmt, d.Node.Name)
		msg := fmt.Sprintf(msgFmt, from, ee.Email, body)

		err := smtp.SendMail("smtp.gmail.com:587",
			smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
			from, []string{ee.Email}, []byte(msg))
		if err != nil {
			fmt.Printf("smtp error: %s\n", err)
		} else {
			fmt.Println("Mail sent successfully")
		}

		tick := time.NewTicker(10 * time.Second)
		go func() {
			<-tick.C
			ee.Interval[d.Node.Name] = true
		}()
	}
	ee.BaseElement.Exec(d)
}


type ActuatorElement struct {
	BaseElement
	Aid    int `json:"aid"`
	Values []struct {
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
	Interval map[string]bool
}

type Actuator struct {
	Nid    int `json:"nid"`	// node id
	Aid    int `json:"aid"` // actuator id
	Values []struct {		// action values
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
}

func (ae *ActuatorElement) Exec(d *model.LogicData) {
	log.Println("!!!!in ActuatorElement.Exec !!!!")
	ok, exist := ae.Interval[d.Node.Name]
	if !exist {
		ae.Interval[d.Node.Name] = true
	}
	if ok {
		ae.Interval[d.Node.Name] = false
		go func() {
			
			res := Actuator{
				Nid:    d.Node.Nid,
				Aid:    ae.Aid,
				Values: ae.Values,
			}
			log.Println("in ActuatorElement.Exec, res = ", res)
					
			pbytes, _ := json.Marshal(res)
			buff := bytes.NewBuffer(pbytes)
			addr := (*adapter.AddrMap)[d.Node.Sid] // sink address
			log.Println("in Act.Exec, 받는 주소: " + "http://" + addr.Addr + "/actuator" + " 전달내용: " + string(pbytes))
			resp, err := http.Post("http://"+addr.Addr+"/actuator", "application/json", buff)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
		}()
		
		tick := time.NewTicker(30 * time.Second)
		go func() {
			<-tick.C
			ae.Interval[d.Node.Name] = true
		}()
	}
	ae.BaseElement.Exec(d)
}

type Drone struct {
	DroneID int `json:"drone_id"`
	Values struct {
		WP0 [][]float64 `json:"wp0"` // weigh point 0 (start station -> dest tag path)
		WP1 [][]float64 `json:"wp1"` // weigh point 1 (dest tag -> nearby station path)
	} `json:"values"`
}

type DroneElement struct {
	BaseElement
	DroneID int `json:"drone_id"`
	Value struct {
		WP0 [][]float64 `json:"wp0"`
		WP1 [][]float64 `json:"wp1"`
	} `json:"value"`
	Interval map[string]bool
}

func (de *DroneElement) Exec(d *model.LogicData) {
	log.Println("\t\t!!!!in DroneElement.Exec !!!!")
	ok, exist := de.Interval[d.Node.Name]
	if !exist {
		de.Interval[d.Node.Name] = true
	}
	if ok {
		de.Interval[d.Node.Name] = false
			
		res := Drone{
			DroneID: 1,
			Values: struct{WP0 [][]float64 "json:\"wp0\""; WP1 [][]float64 "json:\"wp1\""}{},
		}

		srcStation := []float64{37.497365670723944, 126.95591909983503} // 위도(lat), 경도(lon)
		tag := []float64{37.49719755738831, 126.95590032437323}
		destStation := []float64{37.4971933013496, 126.95619804955307}

		res.Values.WP0 = append(res.Values.WP0, srcStation)
		res.Values.WP0 = append(res.Values.WP0, tag)
		res.Values.WP1 = append(res.Values.WP1, tag)
		res.Values.WP1 = append(res.Values.WP1, destStation)

		log.Println("\t\tin DroneElement.Exec, res = ", res)
				
		pbytes, _ := json.Marshal(res)
		buff := bytes.NewBuffer(pbytes)
		addr := (*adapter.AddrMap)[d.Node.Sid]
		log.Println("in Drone.Exec, 받는 주소: " + "http://" + addr.Addr + "/drone" + " 전달내용: " + string(pbytes))
		resp, err := http.Post("http://"+addr.Addr+"/drone", "application/json", buff)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		tick := time.NewTicker(10 * time.Second)
		go func() {
			<-tick.C
			de.Interval[d.Node.Name] = true
		}()
	}
	de.BaseElement.Exec(d)
}
