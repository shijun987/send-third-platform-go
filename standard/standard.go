package standard

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"whxph.com/send-third-platform/utils"
	"whxph.com/send-third-platform/xphapi"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var StandData Data

// Data 参考值
type Data struct {
	Humi  float64
	Temp  float64
	Pre   float64
	Windd float64
	Winds float64
	Noise float64
	Pm25  float64
	Pm10  float64
	Tsp   float64
	Co    float64
	O3    float64
	So2   float64
	No2   float64
}

// Start 获取标准值
func Start() {
	logrus.Info("standard start ------")
	updateStandard()
	c := cron.New()
	_ = c.AddFunc("0 0/5 * * * *", updateStandard)
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
	var dataEntity xphapi.DataEntity
	_ = json.Unmarshal(result, &dataEntity)
	StandData.Humi = utils.String2float(dataEntity.Entity[3].EValue)
	StandData.Temp = utils.String2float(dataEntity.Entity[2].EValue)
	StandData.Pre = utils.String2float(dataEntity.Entity[7].EValue)
	StandData.Windd = utils.String2float(dataEntity.Entity[1].EValue)
	StandData.Winds = utils.String2float(dataEntity.Entity[0].EValue)
	StandData.Noise = utils.String2float(dataEntity.Entity[4].EValue)
	StandData.Pm25 = utils.String2float(dataEntity.Entity[5].EValue)
	StandData.Pm10 = utils.String2float(dataEntity.Entity[6].EValue)
	StandData.Co = utils.String2float(dataEntity.Entity[8].EValue)
	StandData.O3 = utils.String2float(dataEntity.Entity[9].EValue)
	StandData.So2 = utils.String2float(dataEntity.Entity[10].EValue)
	StandData.No2 = utils.String2float(dataEntity.Entity[11].EValue)
	if StandData.Humi >= 3276 {
		StandData.Humi = 0
	}
	if StandData.Temp >= 3276 {
		StandData.Temp = 0
	}
	if StandData.Pre >= 3276 {
		StandData.Pre = 0
	}
	if StandData.Windd >= 32767 {
		StandData.Windd = 0
	}
	if StandData.Winds >= 3276 {
		StandData.Winds = 0
	}
	if StandData.Noise >= 3276 {
		StandData.Noise = 0
	}
	if StandData.Pm25 >= 32767 {
		StandData.Pm25 = 0
	}
	if StandData.Pm10 >= 32767 {
		StandData.Pm10 = 0
	}
	if StandData.Co >= 327 {
		StandData.Co = 0
	}
	if StandData.O3 >= 32767 {
		StandData.O3 = 0
	}
	if StandData.So2 >= 32767 {
		StandData.So2 = 0
	}
	if StandData.No2 >= 32767 {
		StandData.No2 = 0
	}
}
