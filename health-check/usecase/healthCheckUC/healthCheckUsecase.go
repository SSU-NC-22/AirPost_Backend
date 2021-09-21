package healthCheckUC
import (
	"encoding/json"
	"io"
	"log"
	"fmt"
	"net"
	"time"

	"github.com/eunnseo/AirPost/health-check/adapter"
	"github.com/eunnseo/AirPost/health-check/domain/repository"
	"github.com/eunnseo/AirPost/health-check/setting"
)

type healthCheckUsecase struct {
	sr    repository.StatusRepo
	event chan interface{}
}

func NewHealthCheckUsecase(sr repository.StatusRepo, e chan interface{}) *healthCheckUsecase {
	hu := &healthCheckUsecase{
		sr:    sr,
		event: e,
	}
	l, err := net.Listen("tcp", setting.Healthsetting.Listen)
	if nil != err {
		log.Fatalf("fail to bind address to Listen; err: %v", err)
	}
	//defer l.Close()

	go func() {
		for {
			conn, err := l.Accept()
			if nil != err {
				log.Printf("fail to accept; err: %v", err)
				continue
			}
			go hu.healthCheck(conn)
		}
	}()
	return hu
}

func (hu *healthCheckUsecase) healthCheck(conn net.Conn) {
	log.Println("===== healthCheck start =====")

	// for {
	recvBuf := make([]byte, 4096)
	n, err := conn.Read(recvBuf)	
	if nil != err {
		if io.EOF == err {
			log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
			return
		}
		log.Printf("fail to receive data; err: %v", err)
		return
	}
	if n > 0 {
		log.Println("health Info : ", string(recvBuf))
		var healthInfo adapter.HealthInfo
		var states adapter.States
		recvBuf = ClearPadding(recvBuf)
		json.Unmarshal(recvBuf, &healthInfo)

		states.State = healthInfo
		states.Timestamp = fmt.Sprint(time.Now().Unix())
		log.Println("convert to json, healthInfo :", healthInfo)
		tmphealth := hu.sr.UpdateTable(states) // 변화가 생긴 것들만 뭘로 변했는지 알려줌 ex : {1 [{1 1} {2 1} {8 0}]}
		log.Println("change occurred, healthInfo.state :", tmphealth)

		hu.event <- tmphealth // go to wu.event in NewWebsocketUsecase websocketUsecase.go
	}
	// }
}

func ClearPadding(buf []byte) []byte {
	var res []byte
	isCleared := false
	for i := 1; i < 4096; i++ {
		if (buf[i-1] == 125) && (buf[i] == 0) {
			res = buf[:i]
			isCleared = true
			break
		}
	}
	if !isCleared {
		return buf
	}
	return res
}
