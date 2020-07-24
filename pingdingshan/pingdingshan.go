package pingdingshan

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
	pingdingshanToken string
	pingdingshanDevices []xphapi.Device
)

// Start pingdingshan
func Start() {
	logrus.Info("pingdingshan start ------")
	pingdingshanUpdateToken()
	pingdingshanUpdateDevices()
	c := cron.New()
	_ = c.AddFunc("0 0 0/12 * * *", pingdingshanUpdateToken)
	_ = c.AddFunc("0 0 0/1 * * *", pingdingshanUpdateDevices)
	_ = c.AddFunc("0 0/10 * * * *", pingdingshanSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func pingdingshanUpdateToken() {
	pingdingshanToken = xphapi.GetToken("pingdingshan", "123456")
}

func pingdingshanUpdateDevices() {
	pingdingshanDevices = xphapi.GetDevices("pingdingshan", pingdingshanToken)
}

func pingdingshanSendData() {
	for _, item := range pingdingshanDevices {
		resp, err := http.Get("http://115.28.187.9:8005/intfa/queryData/" + strconv.Itoa(item.DeviceID))
		if err != nil {
			logrus.Error("获取数据异常")
			return
		}
		result, _ := ioutil.ReadAll(resp.Body)
		var dataEntity xphapi.DataEntity
		_ = json.Unmarshal(result, &dataEntity)
		if len(dataEntity.Entity) > 0 {
			datatime, _ := time.Parse("2006-01-02 15:04:05", dataEntity.Entity[0].Datetime)
			if datatime.After(time.Now().Add(time.Duration(-time.Hour))) {
				id := "101" + strconv.Itoa(item.DeviceID)[2:]
				content := "DevID:|:" + id +
					"#|#Time:|:" + dataEntity.Entity[0].Datetime +
					"#|#HUMI:|:-1#|#TEMP:|:-1#|#PRE:|:0#|#WINDD:|:-1" +
					"#|#WINDS:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[0].EValue)) +
					"#|#NOISE:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[4].EValue)) +
					"#|#PM25:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[5].EValue)) +
					"#|#PM10:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[6].EValue)) +
					"#|#TSP:|:0"
				logrus.Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
				utils.Invoke("http://123.163.55.113:8686/pdssanitate/services/SaveYCJCService", content)
			} else {
				logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
			}
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
