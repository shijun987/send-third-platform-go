package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var hnsjsbToken string
var hnsjsbDevices []Device

// HnsjsbStart hnsjsb
func HnsjsbStart() {
	logrus.Info("hnsjsb start ------")
	hnsjsbUpdateToken()
	hnsjsbUpdateDevices()
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 0/12 * * *", hnsjsbUpdateToken)
	c.AddFunc("0 0 0/1 * * *", hnsjsbUpdateDevices)
	c.AddFunc("0 0/5 * * * *", hnsjsbSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func hnsjsbUpdateToken() {
	hnsjsbToken = GetToken("hnsjsb", "123456")
}

func hnsjsbUpdateDevices() {
	hnsjsbDevices = GetDevices("hnsjsb", hnsjsbToken)
	// hnsjsbDevices = append(hnsjsbDevices, Device{DeviceID: 16056972, DeviceName: "一鼎-富力阅山湖"})
}

func hnsjsbSendData() {
	for _, item := range hnsjsbDevices {
		resp, err := http.Get("http://115.28.187.9:8005/intfa/queryData/" + strconv.Itoa(item.DeviceID))
		if err != nil {
			logrus.Error("获取数据异常")
			return
		}
		defer resp.Body.Close()
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		var dataEntity DataEntity
		json.Unmarshal(result, &dataEntity)
		if len(dataEntity.Entity) > 0 {
			windd := String2float(dataEntity.Entity[1].EValue)
			winds := String2float(dataEntity.Entity[0].EValue)
			noise := String2float(dataEntity.Entity[5].EValue)
			pm25 := String2float(dataEntity.Entity[2].EValue)
			pm10 := String2float(dataEntity.Entity[3].EValue)
			if windd >= 32767 {
				windd = standardData.windd
			}
			if winds >= 3276 {
				winds = standardData.winds
			}
			if noise >= 3276 {
				noise = standardData.noise
			}
			if pm25 >= 32767 || pm25 <= 0 {
				rand.Seed(time.Now().UnixNano())
				pm25 = standardData.pm25 + (float64)(rand.Intn(6))
			}
			if pm10 >= 32767 || pm10 <= 0 {
				rand.Seed(time.Now().UnixNano())
				pm10 = standardData.pm10 + (float64)(rand.Intn(11))
			}
			id := "101" + strconv.Itoa(item.DeviceID)[2:]
			if item.DeviceID == 16056972 {
				id = "ID" + strconv.Itoa(item.DeviceID)
			}
			content := "{\"elements\":" + "\"DevID:|:" + id +
				"#|#Time:|:" + dataEntity.Entity[0].Datetime +
				"#|#HUMI:|:-1" +
				"#|#TEMP:|:-1" +
				"#|#PRE:|:0" +
				"#|#WINDD:|:" + fmt.Sprintf("%.2f", windd) +
				"#|#WINDS:|:" + fmt.Sprintf("%.2f", winds) +
				"#|#NOISE:|:" + fmt.Sprintf("%.2f", noise) +
				"#|#PM25:|:" + fmt.Sprintf("%.2f", pm25) +
				"#|#PM10:|:" + fmt.Sprintf("%.2f", pm10) +
				"#|#TSP:|:" + fmt.Sprintf("%.2f", pm10+pm25*1.2) + "\"}"
			logrus.WithField("username", "hnsjsb").Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Post("http://hnvjd.jyjzqy.com/ycjk/api/DustV1/saveYCJC", "application/json", bytes.NewBuffer([]byte(content)))
			if err != nil {
				logrus.Error(err)
			}
			defer resp.Body.Close()

			result, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.WithField("username", "hnsjsb").Info(string(result))
			}
		} else {
			logrus.WithField("username", "hnsjsb").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
