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

var nanyangToken string
var nanyangDevices []Device

// NanyangStart nanyang
func NanyangStart() {
	logrus.Info("nanyang start ------")
	nanyangUpdateToken()
	nanyangUpdateDevices()
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 0/12 * * *", nanyangUpdateToken)
	c.AddFunc("0 0 0/1 * * *", nanyangUpdateDevices)
	c.AddFunc("0 0/10 * * * *", nanyangSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func nanyangUpdateToken() {
	nanyangToken = GetToken("nanyang", "123456")
}

func nanyangUpdateDevices() {
	nanyangDevices = GetDevices("nanyang", nanyangToken)
}

func nanyangSendData() {
	for _, item := range nanyangDevices {
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
				"#|#HUMI:|:-1#|#TEMP:|:-1#|#PRE:|:0" +
				"#|#WINDD:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[1].EValue)) +
				"#|#WINDS:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[0].EValue)) +
				"#|#NOISE:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[4].EValue)) +
				"#|#PM25:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[2].EValue)) +
				"#|#PM10:|:" + fmt.Sprintf("%.2f", String2float(dataEntity.Entity[3].EValue)) +
				"#|#TSP:|:0"
			logrus.Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
			Invoke("http://111.6.77.46:7013/nysanitate/services/SaveYCJCService", content)
		} else {
			logrus.Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
