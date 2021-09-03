// AirPost
package setting

import (
	"os"
	"strconv"
)

func GetenvInt(target *int, init int, env string) {
	var err error

	temp := os.Getenv(env)
	if temp == "" {
		*target = init
	} else {
		if *target, err = strconv.Atoi(temp); err != nil {
			*target = init
		}
	}
}

/**************************************************************/
/* Health setting                                             */
/**************************************************************/
type Health struct {
	Server string
	Listen string
}

func GetenvStr(target *string, init, env string) {
	*target = os.Getenv(env)
	if *target == "" {
		*target = init
	}
}

func (hs *Health) Getenv() {
	// GetenvStr(&hs.Server, "221.140.150.7:8083", "HEALTH_SERVER")
	// GetenvStr(&hs.Listen, "221.140.150.7:8085", "HEALTH_LISTEN")
	GetenvStr(&hs.Server, "192.168.0.18:8083", "HEALTH_SERVER")
	GetenvStr(&hs.Listen, "192.168.0.18:8085", "HEALTH_LISTEN")
}

var Healthsetting = &Health{}

/**************************************************************/
/* App setting                                                */
/**************************************************************/
type App struct {
	Server      string
	RequestPath string
}

func (as *App) Getenv() {
	as.Server = os.Getenv("APP_SERVER")
	if as.Server == "" {
		as.Server = "localhost:8081"
	}
	as.RequestPath = "/regist/sink"
}

var Appsetting = &App{}

/**************************************************************/
/* Status setting                                             */
/**************************************************************/
type Status struct {
	Count int
	Tick  int
	Drop  int
}

func (ss *Status) Getenv() {
	GetenvInt(&ss.Count, 5, "STATUS_COUNT")
	GetenvInt(&ss.Tick, 60, "STATUS_TICK")
	GetenvInt(&ss.Drop, 1, "STATUS_DROP")
}

var StatusSetting = &Status{}

func init() {
	Healthsetting.Getenv()
	Appsetting.Getenv()
	StatusSetting.Getenv()

	//fmt.Printf("Health : &v\nApp : %v\nStatus : %v\n\n", Healthsetting, Appsetting, StatusSetting)
}
