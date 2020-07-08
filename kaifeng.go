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

var kaifengToken string
var kaifengDevices []Device

// KaifengStart kaifeng
func KaifengStart() {
	logrus.Info("kaifeng start ------")
	kaifengUpdateToken()
	kaifengUpdateDevices()
	c := cron.New()
	c.AddFunc("0 0 0/12 * * *", kaifengUpdateToken)
	c.AddFunc("0 0 0/1 * * *", kaifengUpdateDevices)
	c.AddFunc("0 0/10 * * * *", kaifengSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func kaifengUpdateToken() {
	kaifengToken = GetToken("kaifeng", "123456")
}

func kaifengUpdateDevices() {
	kaifengDevices = GetDevices("kaifeng", kaifengToken)
}

func kaifengSendData() {
	for _, item := range kaifengDevices {
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
			datatime, _ := time.Parse("2006-01-02 15:04:05", dataEntity.Entity[0].Datetime)
			if datatime.After(time.Now().Add(time.Duration(-time.Hour))) {
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
				Invoke("http://221.176.156.141:8686/kfsanitate/services/SaveYCJCService", content)
			} else {
				logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
			}
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
