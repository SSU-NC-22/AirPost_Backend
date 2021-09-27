package healthCheckUC

import (
	"fmt"
	"time"
	"log"

	"github.com/eunnseo/AirPost/health-check/setting"
	"github.com/go-resty/resty/v2"
)

var (
	appClient *resty.Client
	url       string
)

func init() {
	log.Println("init in healthCheckHelper")
	appClient = resty.New()
	appClient.SetRetryCount(2).SetRetryWaitTime(100 * time.Millisecond).SetTimeout(500 * time.Millisecond)
	url = fmt.Sprintf("http://%s%s", setting.Appsetting.Server, setting.Appsetting.RequestPath)
	log.Println("url : ", url)
}
