package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eunnseo/AirPost/logic-core/adapter"
	"github.com/eunnseo/AirPost/logic-core/domain/model"
	"github.com/eunnseo/AirPost/logic-core/setting"
)

// SMTP settings are read from the environment so no credentials are committed
// in source. Defaults target a local MailHog instance for development; set
// SMTP_HOST/SMTP_PORT and SMTP_FROM/SMTP_PASS for production (e.g. Gmail).
var (
	from     = getenv("SMTP_FROM", "airpost@localhost")
	pass     = os.Getenv("SMTP_PASS")
	smtpHost = getenv("SMTP_HOST", "localhost")
	smtpPort = getenv("SMTP_PORT", "1025")
)

// getenv returns the env var value or a fallback when it is unset.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// sendMail delivers a message via the configured SMTP server. Auth is used only
// when a password is set, so unauthenticated dev servers like MailHog work.
func sendMail(to []string, msg string) error {
	addr := smtpHost + ":" + smtpPort
	var auth smtp.Auth
	if pass != "" {
		auth = smtp.PlainAuth("", from, pass, smtpHost)
	}
	return smtp.SendMail(addr, auth, from, to, []byte(msg))
}

type EmailElement struct {
	BaseElement
	Email    string `json:"text"`
	Interval map[string]bool
	mu       sync.Mutex
}

func (ee *EmailElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in EmailElement.Exec !!!!")
	ee.mu.Lock()
	ok, exist := ee.Interval[d.Node.Name]

	if !exist {
		ee.Interval[d.Node.Name] = true
	}
	if ok {
		ee.Interval[d.Node.Name] = false
		ee.mu.Unlock()

		to := []string{ee.Email}
		body := fmt.Sprintf("node(%s)", d.Node.Name)
		msg := "From: " + from + "\n" +
			"To: " + strings.Join(to, ",") + "\n" +
			"Subject: AirPost email\n" + body

		err := sendMail(to, msg)

		if err != nil {
			fmt.Printf("smtp error: %s\n", err)
		} else {
			fmt.Println("Mail sent successfully")
		}

		tick := time.NewTicker(10 * time.Second)
		go func() {
			<-tick.C
			ee.mu.Lock()
			ee.Interval[d.Node.Name] = true
			ee.mu.Unlock()
		}()
	} else {
		ee.mu.Unlock()
	}
	ee.BaseElement.Exec(d)
}


type ActuatorElement struct {
	BaseElement
	Name   string `json:"name"`
	Values []struct {
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
	Interval map[string]bool
	mu       sync.Mutex
}

type Actuator struct {
	Nid    string `json:"nid"`  // node id
	Name   string `json:"name"` // actuator name
	Values []struct {           // action values
		Value int `json:"value"`
		Sleep int `json:"sleep"`
	} `json:"values"`
}

func (ae *ActuatorElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in ActuatorElement.Exec !!!!")
	ae.mu.Lock()
	ok, exist := ae.Interval[d.Node.Name]
	if !exist {
		ae.Interval[d.Node.Name] = true
	}
	if ok {
		ae.Interval[d.Node.Name] = false
		ae.mu.Unlock()
		go func() {
			
			res := Actuator{
				Nid:    "STA" + strconv.Itoa(d.Node.Nid),
				Name:   ae.Name,
				Values: ae.Values,
			}
					
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
		
		tick := time.NewTicker(1 * time.Second)
		go func() {
			<-tick.C
			ae.mu.Lock()
			ae.Interval[d.Node.Name] = true
			ae.mu.Unlock()
		}()
	} else {
		ae.mu.Unlock()
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
	log.Println("\t!!!!in DroneElement.Exec !!!!")
			
	if !de.Sent {
		de.Sent = true
		go func() {
			res := Drone{
				Nid:    "DRO0",
				Values: de.Values,
				Tagidx: 1,
			}
					
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

type AlarmElement struct { // 도착 알림을 위한 action
	BaseElement
	Email      string `json:"email"`
	OrderNum   string `json:"ordernum"`
	SrcStation string `json:"src_station"`
	DestTag    string `json:"dest_tag"`
	SrcName    string `json:"src_name"`
	DestName   string `json:"dest_name"`
}

func (ae *AlarmElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in AlarmElement.Exec !!!!")

	to := []string{ae.Email}
	subject := "AirPost 배송 완료 - 송장번호(" + ae.OrderNum + ")"
	body := "송장번호 : " + ae.OrderNum + "\r\n" +
		"출발 스테이션 : " + ae.SrcStation + "\r\n" +
		"도착 태그 : " + ae.DestTag + "\r\n" + "\r\n" +
		ae.DestName + "님, " + ae.SrcName + "님이 발송하신 물품이 배송 완료되었습니다."

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: " + subject + "\n\n" + body

	err := sendMail(to, msg)

	if err != nil {
		log.Panicln("smtp send error: ", err)
	} else {
		log.Println("smtp send ok")
	}

	ae.BaseElement.Exec(d)
}

type MovingElement struct { // 데이터베이스에 드론 위치를 동기화시키기 위한 action
	BaseElement
	Nid int `json:"nid"`
}

type Moving struct {
	Nid int `json:"nid"` // drone node id
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
		Alt float64 `json:"alt"`
	} `json:"location"`
}

func (me *MovingElement) Exec(d *model.LogicData) {
	log.Println("\t!!!!in TrackingElement.Exec !!!!")
			
	go func() {
		res := Moving{
			Nid:      me.Nid,
			Location: struct{Lat float64 "json:\"lat\""; Lon float64 "json:\"lon\""; Alt float64 "json:\"alt\""}{
				Lat: d.Values["lat"],
				Lon: d.Values["long"],
				Alt: d.Values["alt"],
			},
		}
				
		pbytes, _ := json.Marshal(res)
		buff := bytes.NewBuffer(pbytes)
		log.Println("in Tracking.Exec, 받는 주소: " + "http://" + setting.Appsetting.Server + "/regist/node/update" + " 전달내용: " + string(pbytes))
		resp, err := http.Post("http://"+setting.Appsetting.Server+"/regist/node/update", "application/json", buff)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}()

	me.BaseElement.Exec(d)
}
