// TCP Client
package main

import (
	"log"
	"net"
)

func main() {
	log.Println("TCP Client start")
	conn, err := net.Dial("tcp", "192.168.0.18:8085")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}
	// defer conn.Close()

	log.Println("success to connect to server")
	send := `{"sid": 1, "state": [{"nid": 3, "state": false, "battery": 70, "location": [37.49304548532277, 126.96038837990257, 10]}]}`

	_, err = conn.Write([]byte(send))
	if err != nil {
		log.Println("failed to write data : ", err)
	} else {
		log.Println("success to write data \nmsg : ", send)
	}

	// 	time.Sleep(time.Duration(1) * time.Second)
}
