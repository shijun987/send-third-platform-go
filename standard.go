package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var standardData StandardData

// StandardData 参考值
type StandardData struct {
	humi  float64
	temp  float64
	pre   float64
	windd float64
	winds float64
	noise float64
	pm25  float64
	pm10  float64
	tsp   float64
	co    float64
	o3    float64
	so2   float64
	no2   float64
}

// StandardStart 获取标准值
func StandardStart() {
	logrus.Info("standard start ------")
	updateStandard()
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0/5 * * * *", updateStandard)
	c.Start()
	defer c.Stop()
	select {}
}

func updateStandard() {
	resp, err := http.Get("http://115.28.187.9:8005/intfa/queryData/16054169")
	if err != nil {
		logrus.Error("获取数据异常")
		return
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	var dataEntity DataEntity
	json.Unmarshal(result, &dataEntity)
	standardData.humi = String2float(dataEntity.Entity[3].EValue)
	standardData.temp = String2float(dataEntity.Entity[2].EValue)
	standardData.pre = String2float(dataEntity.Entity[7].EValue)
	standardData.windd = String2float(dataEntity.Entity[1].EValue)
	standardData.winds = String2float(dataEntity.Entity[0].EValue)
	standardData.noise = String2float(dataEntity.Entity[4].EValue)
	standardData.pm25 = String2float(dataEntity.Entity[5].EValue)
	standardData.pm10 = String2float(dataEntity.Entity[6].EValue)
	standardData.co = String2float(dataEntity.Entity[8].EValue)
	standardData.o3 = String2float(dataEntity.Entity[9].EValue)
	standardData.so2 = String2float(dataEntity.Entity[10].EValue)
	standardData.no2 = String2float(dataEntity.Entity[11].EValue)
	if standardData.humi >= 3276 {
		standardData.humi = 0
	}
	if standardData.temp >= 3276 {
		standardData.temp = 0
	}
	if standardData.pre >= 3276 {
		standardData.pre = 0
	}
	if standardData.windd >= 32767 {
		standardData.windd = 0
	}
	if standardData.winds >= 3276 {
		standardData.winds = 0
	}
	if standardData.noise >= 3276 {
		standardData.noise = 0
	}
	if standardData.pm25 >= 32767 {
		standardData.pm25 = 0
	}
	if standardData.pm10 >= 32767 {
		standardData.pm10 = 0
	}
	if standardData.co >= 327 {
		standardData.co = 0
	}
	if standardData.o3 >= 32767 {
		standardData.o3 = 0
	}
	if standardData.so2 >= 32767 {
		standardData.so2 = 0
	}
	if standardData.no2 >= 32767 {
		standardData.no2 = 0
	}
}
