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
		log.Println("\t\tin EmailElement.Exec, ok")
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
	log.Println("\t\tin EmailElement.Exec, before BaseElement.Exec(d)")
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
	log.Println("\t\t!!!!in ActuatorElement.Exec !!!!")
	ok, exist := ae.Interval[d.Node.Name]
	if !exist {
		ae.Interval[d.Node.Name] = true
	}
	if ok {
		log.Println("\t\tin ActuatorElement.Exec, ok")
		ae.Interval[d.Node.Name] = false
		go func() {
			log.Println("\t\tin ActuatorElement.Exec, go routine 1")
			
			res := Actuator{
				Nid:    d.Node.Nid,
				Aid:    ae.Aid,
				Values: ae.Values,
			}
			log.Println("\t\tin ActuatorElement.Exec, res = ", res)
					
			pbytes, _ := json.Marshal(res)
			buff := bytes.NewBuffer(pbytes)
			addr := (*adapter.AddrMap)[d.Node.Sid] // sink address
			log.Println("\t\tin ActuatorElement.Exec, addr = ", addr)
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
			log.Println("\t\tin ActuatorElement.Exec, go routine 2")
		}()
	}
	log.Println("\t\tin ActuatorElement.Exec, before BaseElement.Exec(d)")
	ae.BaseElement.Exec(d)
}

