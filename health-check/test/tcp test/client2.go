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
	send := `{"sid": 2, "state": [{"nid": 1, "state": true, "battery": 16.1, "location": [37.49304548532277, 126.96038837990257, 10]},{"nid": 2, "state": true, "battery": 16.2, "location": [37.49304548532277, 126.96038837990257, 10]}]}`

	_, err = conn.Write([]byte(send))
	if err != nil {
		log.Println("failed to write data : ", err)
	} else {
		log.Println("success to write data \nmsg : ", send)
	}

	// 	time.Sleep(time.Duration(1) * time.Second)

	conn, err = net.Dial("tcp", "192.168.0.18:8085")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}
	// defer conn.Close()

	log.Println("success to connect to server")
	send = `{"sid": 3, "state": [{"nid": 4, "state": true, "battery": 14.9, "location": [37.49304548532277, 126.96038837990257, 10]}]}`

	_, err = conn.Write([]byte(send))
	if err != nil {
		log.Println("failed to write data : ", err)
	} else {
		log.Println("success to write data \nmsg : ", send)
	}
}
