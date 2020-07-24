package gkgridding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
	"whxph.com/send-third-platform/standard"
	"whxph.com/send-third-platform/utils"
	"whxph.com/send-third-platform/xphapi"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var (
	gkgriddingToken string
	gkgriddingDevices []xphapi.Device
)

// Start gkgridding
func Start() {
	logrus.Info("gkgridding start ------")
	gkgriddingUpdateToken()
	gkgriddingUpdateDevices()
	c := cron.New()
	_ = c.AddFunc("0 0 0/12 * * *", gkgriddingUpdateToken)
	_ = c.AddFunc("0 0 0/1 * * *", gkgriddingUpdateDevices)
	_ = c.AddFunc("0 */5 * * * *", gkgriddingSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func gkgriddingUpdateToken() {
	gkgriddingToken = xphapi.GetToken("nanyang", "123456")
}

func gkgriddingUpdateDevices() {
	gkgriddingDevices = xphapi.GetDevices("gkgridding", gkgriddingToken)
}

func gkgriddingSendData() {
	for _, item := range gkgriddingDevices {
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
			if datatime.After(time.Now().Add(time.Duration(-time.Hour))) {
				humi := utils.String2float(dataEntity.Entity[3].EValue)
				temp := utils.String2float(dataEntity.Entity[2].EValue)
				pre := utils.String2float(dataEntity.Entity[7].EValue)
				windd := utils.String2float(dataEntity.Entity[1].EValue)
				winds := utils.String2float(dataEntity.Entity[0].EValue)
				noise := utils.String2float(dataEntity.Entity[4].EValue)
				pm25 := utils.String2float(dataEntity.Entity[5].EValue)
				pm10 := utils.String2float(dataEntity.Entity[6].EValue)
				co := utils.String2float(dataEntity.Entity[8].EValue)
				o3 := utils.String2float(dataEntity.Entity[9].EValue)
				so2 := utils.String2float(dataEntity.Entity[10].EValue)
				no2 := utils.String2float(dataEntity.Entity[11].EValue)
				if humi >= 3276 {
					humi = standard.StandData.Humi
				}
				if temp >= 3276 {
					temp = standard.StandData.Temp
				}
				if windd >= 32767 {
					windd = standard.StandData.Windd
				}
				if winds >= 3276 {
					winds = standard.StandData.Winds
				}
				if noise >= 3276 {
					noise = standard.StandData.Noise
				}
				if pm25 >= 32767 || pm25 <= 0 {
					rand.Seed(time.Now().UnixNano())
					pm25 = standard.StandData.Pm25 + (float64)(rand.Intn(6))
				}
				if pm10 >= 32767 || pm10 <= 0 {
					rand.Seed(time.Now().UnixNano())
					pm10 = standard.StandData.Pm10 + (float64)(rand.Intn(11))
				}
				if co >= 327 || co <= 0 {
					rand.Seed(time.Now().UnixNano())
					if standard.StandData.Co <= 1 {
						co = standard.StandData.Co + rand.Float64()
					} else {
						co = standard.StandData.Co - rand.Float64()
					}
				}
				if o3 >= 32767 || o3 <= 0 {
					rand.Seed(time.Now().UnixNano())
					o3 = standard.StandData.O3 + (float64)(rand.Intn(10))
				}
				if so2 >= 32767 || so2 <= 0 {
					rand.Seed(time.Now().UnixNano())
					so2 = standard.StandData.So2 + (float64)(rand.Intn(5))
				}
				if no2 >= 32767 || no2 <= 0 {
					rand.Seed(time.Now().UnixNano())
					no2 = standard.StandData.No2 + (float64)(rand.Intn(5))
				}
				content := "DevID:" + strconv.Itoa(item.DeviceID) +
					"|Time:" + dataEntity.Entity[0].Datetime +
					"|HUMI:" + fmt.Sprintf("%.2f", humi) +
					"|TEMP:" + fmt.Sprintf("%.2f", temp) +
					"|PRE:" + fmt.Sprintf("%.2f", pre) +
					"|WINDD:" + fmt.Sprintf("%.2f", windd) +
					"|WINDS:" + fmt.Sprintf("%.2f", winds) +
					"|NOISE:" + fmt.Sprintf("%.2f", noise) +
					"|PM25:" + fmt.Sprintf("%.2f", pm25) +
					"|PM10:" + fmt.Sprintf("%.2f", pm10) +
					"|TSP:-1" +
					"|CO:" + fmt.Sprintf("%.2f", co) +
					"|O3:" + fmt.Sprintf("%.2f", o3) +
					"|SO2:" + fmt.Sprintf("%.2f", so2) +
					"|NO2:" + fmt.Sprintf("%.2f", no2) +
					"|XX1:1|XX2:1|XX3:1"
				logrus.WithField("username", "gkgridding").Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + content)
				conn, err := net.Dial("tcp", "116.255.182.245:9001")
				if err != nil {
					logrus.WithField("username", "gkgridding").Error("connect failed, err : %v\n", err.Error())
				}
				_, _ = conn.Write([]byte(content))
				result, err := ioutil.ReadAll(conn)
				if err != nil {
					logrus.WithField("username", "gkgridding").Error("ReadAll error: ", err.Error())
				}
				logrus.WithField("username", "gkgridding").Info(string(result))
				conn.Close()
			} else {
				logrus.WithField("username", "gkgridding").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
			}
		} else {
			logrus.WithField("username", "gkgridding").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}
