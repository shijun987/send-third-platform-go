package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var guokongToken string
var guokongDevices []Device

// GuokongStart guokong
func GuokongStart() {
	logrus.Info("guokong start ------")
	guokongUpdateToken()
	guokongUpdateDevices()
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 0/12 * * *", guokongUpdateToken)
	c.AddFunc("0 0 0/1 * * *", guokongUpdateDevices)
	c.AddFunc("0 0/10 * * * *", guokongSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func guokongUpdateToken() {
	guokongToken = GetToken("guokong", "123456")
}

func guokongUpdateDevices() {
	guokongDevices = GetDevices("guokong", guokongToken)
}

func guokongSendData() {
	for _, item := range guokongDevices {
		resp, err := http.Get("http://115.28.187.9:8005/intfa/queryData/" + strconv.Itoa(item.DeviceID))
		if err != nil {
			logrus.Error("获取数据异常")
			return
		}
		defer resp.Body.Close()
		result, _ := ioutil.ReadAll(resp.Body)
		var dataEntity DataEntity
		json.Unmarshal(result, &dataEntity)
		if len(dataEntity.Entity) > 0 {
			id := "101" + strconv.Itoa(item.DeviceID)[2:]
			content := "DevID:|:" + id +
				"#|#Time:|:" + dataEntity.Entity[0].Datetime +
				"#|#HUMI:|:-1#|#TEMP:|:-1#|#PRE:|:0#|#WINDD:|:-1" +
				"#|#WINDS:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[0].EValue)) +
				"#|#NOISE:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[1].EValue)) +
				"#|#PM25:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[2].EValue)) +
				"#|#PM10:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[3].EValue)) +
				"#|#TSP:|:0"
			logrus.Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
			Invoke("http://27.50.132.176/sanitate/services/SaveYCJCService", content)
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
