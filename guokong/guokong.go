package guokong

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"whxph.com/send-third-platform/utils"
	"whxph.com/send-third-platform/xphapi"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var (
	guokongToken string
	guokongDevices []xphapi.Device
)

// Start guokong
func Start() {
	logrus.Info("guokong start ------")
	guokongUpdateToken()
	guokongUpdateDevices()
	c := cron.New()
	_ = c.AddFunc("0 0 0/12 * * *", guokongUpdateToken)
	_ = c.AddFunc("0 0 0/1 * * *", guokongUpdateDevices)
	_ = c.AddFunc("0 0/10 * * * *", guokongSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func guokongUpdateToken() {
	guokongToken = xphapi.GetToken("guokong", "123456")
}

func guokongUpdateDevices() {
	guokongDevices = xphapi.GetDevices("guokong", guokongToken)
}

func guokongSendData() {
	for _, item := range guokongDevices {
		resp, err := http.Get("http://115.28.187.9:8005/intfa/queryData/" + strconv.Itoa(item.DeviceID))
		if err != nil {
			logrus.Error("获取数据异常")
			return
		}
		result, _ := ioutil.ReadAll(resp.Body)
		var dataEntity xphapi.DataEntity
		_ = json.Unmarshal(result, &dataEntity)
		if len(dataEntity.Entity) > 0 {
			id := "101" + strconv.Itoa(item.DeviceID)[2:]
			content := "DevID:|:" + id +
				"#|#Time:|:" + dataEntity.Entity[0].Datetime +
				"#|#HUMI:|:-1#|#TEMP:|:-1#|#PRE:|:0#|#WINDD:|:-1" +
				"#|#WINDS:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[0].EValue)) +
				"#|#NOISE:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[1].EValue)) +
				"#|#PM25:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[2].EValue)) +
				"#|#PM10:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[3].EValue)) +
				"#|#TSP:|:0"
			logrus.Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
			utils.Invoke("http://27.50.132.176/sanitate/services/SaveYCJCService", content)
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
