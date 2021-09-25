package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"
    "strings"

	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

const (
	GoogleSMTPServer = "smtp.gmail.com"
	from = "REDACTED"
	pass = ""
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

		to := []string{ee.Email}
		body := fmt.Sprintf("node(%s)", d.Node.Name)
		msg := "From: " + from + "\n" +
			"To: " + strings.Join(to, ",") + "\n" +
			"Subject: AirPost email\n" + body

		err := smtp.SendMail(GoogleSMTPServer + ":587",
			smtp.PlainAuth("", from, pass, GoogleSMTPServer),
			from, to, []byte(msg))

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
	Nid    string 	   `json:"nid"` // drone node id
	Values [][]float64 `json:"values"`
	Tagidx int 		   `json:"tagidx"` // values 내에서 tag가 몇번째 index인지 (0~)
}

type DroneElement struct {
	BaseElement
	Nid      string 	 `json:"nid"`
	Values   [][]float64 `json:"values"`
	Tagidx   int 		 `json:"tagidx"`
	Sent	 bool		 `json:"sent"`
}

func (de *DroneElement) Exec(d *model.LogicData) {
	log.Println("\t\t!!!!in DroneElement.Exec !!!!")
			
	if !de.Sent {
		de.Sent = true
		go func() {
			res := Drone{
				Nid:    de.Nid,
				Values: de.Values,
				Tagidx: de.Tagidx,
			}
			log.Println("\t\tin DroneElement.Exec, res = ", res)
					
			pbytes, _ := json.Marshal(res)
			buff := bytes.NewBuffer(pbytes)
			addr := (*adapter.AddrMap)[d.Node.Sid]
			log.Println("in Drone.Exec, 받는 주소: " + "http://" + addr.Addr + "/drone" + " 전달내용: " + string(pbytes))
			resp, err := http.Post("http://"+addr.Addr+"/drone", "application/json", buff)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
		}()
	}
	de.BaseElement.Exec(d)
}

type AlarmElement struct {
	BaseElement
	Email      string `json:"email"`
	OrderNum   string `json:"ordernum"`
	SrcStation string `json:"src_station"`
	DestTag    string `json:"dest_tag"`
	SrcName    string `json:"src_name"`
	DestName   string `json:"dest_name"`
}

func (ae *AlarmElement) Exec(d *model.LogicData) {
	log.Println("\t\t!!!!in AlarmElement.Exec !!!!")

	to := []string{ae.Email}
	subject := "AirPost 배송 완료 - 송장번호(" + ae.OrderNum + ")"
	body := "송장번호 : " + ae.OrderNum + "\r\n" +
		"출발 스테이션 : " + ae.SrcStation + "\r\n" +
		"도착 태그 : " + ae.DestTag + "\r\n" + "\r\n" +
		ae.DestName + "님, " + ae.SrcName + "님이 발송하신 물품이 배송 완료되었습니다."

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" + body

	err := smtp.SendMail(GoogleSMTPServer + ":587",
		smtp.PlainAuth("", from, pass, GoogleSMTPServer),
		from, to, []byte(msg))

	if err != nil {
		log.Panicln("smtp send error: ", err)
	} else {
		log.Println("smtp send ok")
	}

	ae.BaseElement.Exec(d)
}
