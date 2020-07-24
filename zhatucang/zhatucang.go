package zhatucang

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
	zhatucangToken string
	zhatucangDevices []xphapi.Device
)

// Start zhatucang
func Start() {
	logrus.Info("zhatucang start ------")
	zhatucangUpdateToken()
	zhatucangUpdateDevices()
	c := cron.New()
	_ = c.AddFunc("0 0 0/12 * * *", zhatucangUpdateToken)
	_ = c.AddFunc("0 0 0/1 * * *", zhatucangUpdateDevices)
	_ = c.AddFunc("0 0/10 * * * *", zhatucangSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func zhatucangUpdateToken() {
	zhatucangToken = xphapi.GetToken("zhatucang", "123456")
}

func zhatucangUpdateDevices() {
	zhatucangDevices = xphapi.GetDevices("zhatucang", zhatucangToken)
}

func zhatucangSendData() {
	for _, item := range zhatucangDevices {
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
			if datatime.After(time.Now().Add(-time.Hour)) {
				id := "101" + strconv.Itoa(item.DeviceID)[2:]
				content := "DevID:|:" + id +
					"#|#Time:|:" + dataEntity.Entity[0].Datetime +
					"#|#HUMI:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[3].EValue)) +
					"#|#TEMP:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[2].EValue)) +
					"#|#PRE:|:0" +
					"#|#WINDD:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[1].EValue)) +
					"#|#WINDS:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[0].EValue)) +
					"#|#NOISE:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[4].EValue)) +
					"#|#PM25:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[5].EValue)) +
					"#|#PM10:|:" + fmt.Sprintf("%.2f", utils.String2float(dataEntity.Entity[6].EValue)) +
					"#|#TSP:|:0"
				logrus.Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
				utils.Invoke("http://42.236.61.105:8686/xncjk/services/SaveYCJCService", content)
			} else {
				logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
			}
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		resp.Body.Close()
		time.Sleep(500 * time.Millisecond)
	}
}
