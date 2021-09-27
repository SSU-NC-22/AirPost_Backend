// TCP Client
package main

import (
	"log"
	"net"
)

func main() {
	/*** sid: 1 ***/
	conn, err := net.Dial("tcp", "192.168.0.18:8085")
	if nil != err {
		log.Fatalf("failed to connect to server")
	}
	// defer conn.Close()

	log.Println("success to connect to server 1")
	send := `{"sid": 1, "state": [{"nid": 3, "state": true, "battery": 15.5, "location": [37.49304548532277, 126.96038837990257, 10]}]}`

	_, err = conn.Write([]byte(send))
	if err != nil {
		log.Println("failed to write data : ", err)
	} else {
		log.Println("success to write data \nmsg : ", send)
	}

	// /*** sid: 2 ***/
	// conn, err = net.Dial("tcp", "192.168.0.18:8085")
	// if nil != err {
	// 	log.Fatalf("failed to connect to server")
	// }
	// // defer conn.Close()


	// log.Println("success to connect to server 2")
	// send = `{"sid": 2, "state": [{"nid": 1, "state": true, "battery": 16.1, "location": [37.49304548532277, 126.96038837990257, 10]},{"nid": 2, "state": true, "battery": 16.2, "location": [37.49304548532277, 126.96038837990257, 10]}]}`

	// _, err = conn.Write([]byte(send))
	// if err != nil {
	// 	log.Println("failed to write data : ", err)
	// } else {
	// 	log.Println("success to write data \nmsg : ", send)
	// }

	// /*** sid: 3 ***/
	// conn, err = net.Dial("tcp", "192.168.0.18:8085")
	// if nil != err {
	// 	log.Fatalf("failed to connect to server")
	// }
	// // defer conn.Close()

	// log.Println("success to connect to server 3")
	// send = `{"sid": 3, "state": [{"nid": 4, "state": true, "battery": 14.9, "location": [37.49304548532277, 126.96038837990257, 10]}]}`

	// _, err = conn.Write([]byte(send))
	// if err != nil {
	// 	log.Println("failed to write data : ", err)
	// } else {
	// 	log.Println("success to write data \nmsg : ", send)
	// }

	// 	time.Sleep(time.Duration(1) * time.Second)
}
