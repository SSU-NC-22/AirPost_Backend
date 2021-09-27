// TCP Client
package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.Println("TCP Client start")

	s1 := `{"sid": 1, "state": [{"nid": 3, "state": false, "battery": 16, "location": [37.49304548532277, 126.96038837990257, 10]}]}`
	s2 := `{"sid": 1, "state": [{"nid": 3, "state": true, "battery": 16, "location": [37.49304548532277, 126.96038837990257, 10]}]}`
	s3 := `{"sid": 1, "state": [{"nid": 3, "state": true, "battery": 16, "location": [37.49304548532277, 126.96038837990257, 10]}]}`
	s4 := `{"sid": 1, "state": [{"nid": 3, "state": true, "battery": 15, "location": [37.49304548532277, 126.96038837990257, 10]}]}`
	s5 := `{"sid": 1, "state": [{"nid": 3, "state": false, "battery": 15, "location": [37.49304548532277, 126.96038837990257, 10]}]}`
	s6 := `{"sid": 1, "state": [{"nid": 3, "state": false, "battery": 15, "location": [37.49304548532277, 126.96038837990257, 10]}]}`
	sl := []string{s1, s2, s3, s4, s5, s6}

	for _, send := range sl {
		conn, err := net.Dial("tcp", "192.168.0.18:8085")
		if nil != err {
			log.Fatalf("failed to connect to server")
		}
		// defer conn.Close()

		log.Println("success to connect to server")

		_, err = conn.Write([]byte(send))
		if err != nil {
			log.Println("failed to write data : ", err)
		} else {
			log.Println("success to write data \nmsg : ", send)
		}

		time.Sleep(time.Duration(2) * time.Second)
	}
}
