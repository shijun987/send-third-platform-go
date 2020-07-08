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

var nanyang2019Token string
var nanyang2019Devices []Device

// Nanyang2019Start nanyang2019
func Nanyang2019Start() {
	logrus.Info("nanyang2019 start ------")
	nanyang2019UpdateToken()
	nanyang2019UpdateDevices()
	c := cron.New()
	c.AddFunc("0 0 0/12 * * *", nanyang2019UpdateToken)
	c.AddFunc("0 0 0/1 * * *", nanyang2019UpdateDevices)
	c.AddFunc("0 0/10 * * * *", nanyang2019SendData)
	c.Start()
	defer c.Stop()
	select {}
}

func nanyang2019UpdateToken() {
	nanyang2019Token = GetToken("nanyang2019", "123456")
}

func nanyang2019UpdateDevices() {
	nanyang2019Devices = GetDevices("nanyang2019", nanyang2019Token)
}

func nanyang2019SendData() {
	for _, item := range nanyang2019Devices {
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
				Invoke("http://111.6.77.46:7013/nysanitate/services/SaveYCJCService", content)
			} else {
				logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
			}
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
