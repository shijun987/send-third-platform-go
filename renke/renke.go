package renke

import (
	"encoding/hex"
	"encoding/json"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"whxph.com/send-third-platform/utils"
)

// Renke 消息体
type Renke struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

// Data 数据
type Data struct {
	GroupID        string         `json:"groupID"`
	DeviceKey      string         `json:"deviceKey"`
	DeviceAddr     int            `json:"deviceAddr"`
	NodeID         int            `json:"nodeID"`
	NodeType       int            `json:"nodeType"`
	DeviceDisabled bool           `json:"deviceDisabled"`
	DeviceName     string         `json:"deviceName"`
	Lng            float64        `json:"lng"`
	Lat            float64        `json:"lat"`
	DeviceStatus   int            `json:"deviceStatus"`
	RealTimeData   []RealTimeData `json:"realTimeData"`
}

// RealTimeData 实时数据
type RealTimeData struct {
	DataName  string `json:"dataName"`
	DataValue string `json:"dataValue"`
	IsAlarm   bool   `json:"isAlarm"`
	AlarmMsg  string `json:"alarmMsg"`
	Alarm     bool   `json:"alarm"`
}

// Start 建大仁科温湿度计
func Start() {
	logrus.Info("renke start ------")
	renkeSendData()
	c := cron.New()
	_ = c.AddFunc("0 0/1 * * * *", renkeSendData)
	c.Start()
	defer c.Stop()
	select {}
}

func renkeSendData() {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "http://www.0531yun.cn/wsjc/app/GetDeviceData?groupId=103001", nil)
	if err != nil || req == nil {
		logrus.Error("ReadAll error: ", err.Error())
	}
	req.Header.Set("userId", "9998")
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("ReadAll error: ", err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var renkeData Renke
	_ = json.Unmarshal(body, &renkeData)

	if renkeData.Data[0].RealTimeData[0].DataValue != "" {
		temperature, _ := strconv.ParseFloat(renkeData.Data[0].RealTimeData[0].DataValue, 64)
		temperature = temperature * 10
		humidity, _ := strconv.ParseFloat(renkeData.Data[0].RealTimeData[1].DataValue, 64)
		humidity = humidity * 10
		sendData := make([]byte, 72)
		sendData[0] = 0x15
		sendData[1] = 0x01
		sendData[2] = 0x5C
		sendData[3] = 0x58
		sendData[4] = 0xA2
		sendData[5] = 0x40
		sendData[6] = byte(int16(temperature) >> 8)
		sendData[7] = byte(int16(temperature) & 0xFF)
		sendData[8] = byte(int16(humidity) >> 8)
		sendData[9] = byte(int16(humidity) & 0xFF)
		for i := 0; i < 14; i++ {
			sendData[10+i*2] = 0x7F
			sendData[11+i*2] = 0xFF
		}
		for i := 0; i < 32; i++ {
			sendData[38+i] = 0x00
		}
		crc := utils.Crc16(sendData, 70)
		sendData[70] = byte(crc & 0xFF)
		sendData[71] = byte(crc >> 8)
		logrus.Info(strings.ToUpper(hex.EncodeToString(sendData)))
		conn, err := net.Dial("tcp", "47.105.215.208:8880")
		if err != nil {
			logrus.Error("connect failed, err : %v\n", err.Error())
		}
		_, _ = conn.Write(sendData)
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			logrus.Error("ReadAll error: ", err.Error())
		}
		logrus.Info(strings.ToUpper(hex.EncodeToString(result)))
		conn.Close()
	}
}
