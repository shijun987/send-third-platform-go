package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var linghuiToken string
var linghuiDevices []Device

// LinghuiData LinghuiData
type LinghuiData struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

// LinghuiStruct LinghuiStruct
type LinghuiStruct struct {
	Code string        `json:"code"`
	Time string        `json:"time"`
	Data []LinghuiData `json:"data"`
}

// LinghuiStart 灵慧识别
func LinghuiStart() {
	logrus.Info("3961422 start ------")
	linghuiUpdateToken()
	linghuiUpdateDevices()
	c := cron.New()
	c.AddFunc("0 0 0/12 * * *", linghuiUpdateToken)
	c.AddFunc("0 0 0/1 * * *", linghuiUpdateDevices)
	c.AddFunc("0 0/2 * * * *", linghuiSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func linghuiUpdateToken() {
	linghuiToken = NewGetToken("3961422", "123456")
}

func linghuiUpdateDevices() {
	linghuiDevices = NewGetDevices("3961422", linghuiToken)
}

func linghuiSendData() {
	for _, item := range linghuiDevices {
		resp, err := http.Get("http://47.105.215.208:8005/intfa/queryData/" + strconv.Itoa(item.DeviceID))
		if err != nil {
			logrus.Error("获取数据异常")
			return

		}
		defer resp.Body.Close()
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		dataEntity := DataEntity{}
		json.Unmarshal(result, &dataEntity)
		if len(dataEntity.Entity) > 0 {

			linghuiStruct := LinghuiStruct{}
			linghuiStruct.Code = strconv.Itoa(dataEntity.DeviceID)
			linghuiStruct.Time = dataEntity.Entity[0].Datetime
			for _, value := range dataEntity.Entity {
				linghuiData := LinghuiData{}
				linghuiData.Code = value.EName
				temp, _ := strconv.ParseFloat(value.EValue, 64)
				linghuiData.Value = temp
				linghuiStruct.Data = append(linghuiStruct.Data, linghuiData)
			}

			content, err := json.Marshal(linghuiStruct)
			if err != nil {
				continue
			}
			logrus.WithField("username", "3961422").Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + string(content))
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Post("https://ai.forestrycat.com/ic/app/weather/upload", "application/json", bytes.NewBuffer(content))
			if err != nil {
				logrus.Error(err)
			}
			defer resp.Body.Close()

			result, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Error(err)
			} else {
				logrus.WithField("username", "3961422").Info(string(result))
			}
		} else {
			logrus.WithField("username", "3961422").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
