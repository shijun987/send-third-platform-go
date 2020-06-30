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

var hnhebiToken string
var hnhebiDevices []Device

// HnhebiStart hnhebi
func HnhebiStart() {
	logrus.Info("hnhebi start ------")
	hnhebiUpdateToken()
	hnhebiUpdateDevices()
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 0/12 * * *", hnhebiUpdateToken)
	c.AddFunc("0 0 0/1 * * *", hnhebiUpdateDevices)
	c.AddFunc("0 0/10 * * * *", hnhebiSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func hnhebiUpdateToken() {
	hnhebiToken = GetToken("hnhebi", "123456")
}

func hnhebiUpdateDevices() {
	hnhebiDevices = GetDevices("hnhebi", hnhebiToken)
}

func hnhebiSendData() {
	for _, item := range hnhebiDevices {
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
				"#|#HUMI:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[3].EValue)) +
				"#|#TEMP:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[2].EValue)) +
				"#|#PRE:|:0" +
				"#|#WINDD:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[1].EValue)) +
				"#|#WINDS:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[0].EValue)) +
				"#|#NOISE:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[4].EValue)) +
				"#|#PM25:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[5].EValue)) +
				"#|#PM10:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[6].EValue)) +
				"#|#TSP:|:0"
			logrus.Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
			Invoke("http://222.88.186.189:8686/hbsanitate/services/SaveYCJCService", content)
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
