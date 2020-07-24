package czzhujian

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
	"whxph.com/send-third-platform/standard"
	"whxph.com/send-third-platform/utils"
	"whxph.com/send-third-platform/xphapi"

	"github.com/sirupsen/logrus"
)

var (
	token        string
	devices      []xphapi.Device
	sockets      []net.Conn
	serialNumber []uint8
)

// Start czzhujian
func Start() {
	logrus.Info("czzhujian start ------")
	updateToken()
	updateDevices()
	c := cron.New()
	_ = c.AddFunc("0 0 0/12 * * *", updateToken)
	_ = c.AddFunc("0 0 0/1 * * *", updateDevices)
	_ = c.AddFunc("5 */4 * * * *", sendNoise)
	_ = c.AddFunc("15 */4 * * * *", sendDust)
	c.Start()
	defer c.Stop()
	select {}
}

func updateToken() {
	token = xphapi.GetToken("czzhujian", "123456")
}

func updateDevices() {
	devices = xphapi.GetDevices("czzhujian", token)
	sockets = make([]net.Conn, len(devices))
	serialNumber = make([]uint8, len(devices))
	//for index, item := range devices {
	//	conn, err := net.Dial("tcp", "183.203.96.67:10012")
	//	//conn, err := net.Dial("tcp", "127.0.0.1:8899")
	//	if err != nil {
	//		logrus.WithField("username", "czzhujian").Error("connect failed, err : %v\n", err.Error())
	//	}
	//	serialNumber = append(serialNumber, 0)
	//	login(index, item, conn)
	//	sockets = append(sockets, conn)
	//}
}

func sendNoise() {
	var buffer bytes.Buffer
	var dataBuf bytes.Buffer
	for index, item := range devices {
		conn, err := net.Dial("tcp", "183.203.96.67:10012")
		if err != nil {
			logrus.WithField("username", "czzhujian").Error("connect failed, err : %v\n", err.Error())
		}
		serialNumber[index] = 0
		login(index, item, conn)
		sockets[index] = conn
		resp, err := http.Get("http://115.28.187.9:8005/intfa/queryData/" + strconv.Itoa(item.DeviceID))
		if err != nil {
			logrus.Error(err)
			return
		}

		result, _ := ioutil.ReadAll(resp.Body)
		var dataEntity xphapi.DataEntity
		_ = json.Unmarshal(result, &dataEntity)
		if len(dataEntity.Entity) > 0 {
			datatime, _ := time.Parse("2006-01-02 15:04:05", dataEntity.Entity[0].Datetime)
			recDate := dataEntity.Entity[0].Datetime
			if datatime.After(time.Now().Add(-time.Hour)) {
				noise := utils.String2float(dataEntity.Entity[1].EValue)
				if noise >= 3276 {
					noise = standard.StandData.Noise
				}
				dataBuf.Reset()
				dataBuf.WriteString(`{"DataType":"Min",`)
				dataBuf.WriteString(`"DeviceId":"` + strconv.Itoa(item.DeviceID) + `",`)
				dataBuf.WriteString(`"DB":"` + fmt.Sprintf("%0.2f", noise) + `",`)
				dataBuf.WriteString(`"RecDate":"` + recDate + `"}`)
				buffer.Reset()
				buffer.WriteByte(0x15)
				buffer.WriteByte(0x03)
				buffer.WriteByte(serialNumber[index])
				serialNumber[index]++
				buffer.WriteByte(0x00)
				buffer.WriteByte(0x00)
				buffer.WriteByte(0x00)
				buffer.WriteByte(0x00)
				dataLen := uint16(dataBuf.Len())
				buffer.WriteByte(byte(dataLen))
				buffer.WriteByte(byte(dataLen >> 8))
				buffer.Write(dataBuf.Bytes())
				buffer.WriteByte(CheckSum(buffer.Bytes()[1:]))
				buffer.WriteByte(0x02)
				logrus.WithField("username", "czzhujian").Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + dataBuf.String())
				_, _ = sockets[index].Write(buffer.Bytes())
				//time.Sleep(500 * time.Millisecond)
				//reader := bufio.NewReader(sockets[index])
				//msg, err := reader.ReadBytes(0x02)
				//if err != nil {
				//	logrus.WithField("username", "czzhujian").Error(err)
				//}
				//msgHexStr := ""
				//for _, item := range msg {
				//	msgHexStr += fmt.Sprintf("%02x", item)
				//}
				//logrus.WithField("username", "czzhujian").Info(msgHexStr)
			} else {
				logrus.WithField("username", "czzhujian").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
			}
		} else {
			logrus.WithField("username", "czzhujian").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func sendDust() {
	var buffer bytes.Buffer
	var dataBuf bytes.Buffer
	for index, item := range devices {
		resp, err := http.Get("http://115.28.187.9:8005/intfa/queryData/" + strconv.Itoa(item.DeviceID))
		if err != nil {
			logrus.Error(err)
			return
		}

		result, _ := ioutil.ReadAll(resp.Body)
		var dataEntity xphapi.DataEntity
		_ = json.Unmarshal(result, &dataEntity)
		if len(dataEntity.Entity) > 0 {
			datatime, _ := time.Parse("2006-01-02 15:04:05", dataEntity.Entity[0].Datetime)
			recDate := dataEntity.Entity[0].Datetime
			if datatime.After(time.Now().Add(-time.Hour)) {
				pm25 := utils.String2float(dataEntity.Entity[4].EValue)
				pm10 := utils.String2float(dataEntity.Entity[5].EValue)
				if pm25 >= 32767 || pm25 <= 0 {
					rand.Seed(time.Now().UnixNano())
					pm25 = standard.StandData.Pm25 + (float64)(rand.Intn(6))
				}
				if pm10 >= 32767 || pm10 <= 0 {
					rand.Seed(time.Now().UnixNano())
					pm10 = standard.StandData.Pm10 + (float64)(rand.Intn(11))
				}
				dataBuf.Reset()
				dataBuf.WriteString(`{"DataType":"Min",`)
				dataBuf.WriteString(`"DeviceId":"` + strconv.Itoa(item.DeviceID) + `",`)
				dataBuf.WriteString(`"PM10":"` + fmt.Sprintf("%0.2f", pm10) + `",`)
				dataBuf.WriteString(`"PM2.5":"` + fmt.Sprintf("%0.2f", pm25) + `",`)
				dataBuf.WriteString(`"RecDate":"` + recDate + `"}`)
				buffer.Reset()
				buffer.WriteByte(0x15)
				buffer.WriteByte(0x05)
				buffer.WriteByte(serialNumber[index])
				serialNumber[index]++
				buffer.WriteByte(0x00)
				buffer.WriteByte(0x00)
				buffer.WriteByte(0x00)
				buffer.WriteByte(0x00)
				dataLen := uint16(dataBuf.Len())
				buffer.WriteByte(byte(dataLen))
				buffer.WriteByte(byte(dataLen >> 8))
				buffer.Write(dataBuf.Bytes())
				buffer.WriteByte(CheckSum(buffer.Bytes()[1:]))
				buffer.WriteByte(0x02)
				logrus.WithField("username", "czzhujian").Info("[" + strconv.Itoa(dataEntity.DeviceID) + "]: " + dataBuf.String())
				_, _ = sockets[index].Write(buffer.Bytes())
				time.Sleep(500 * time.Millisecond)
				reader := bufio.NewReader(sockets[index])
				msg, err := reader.ReadBytes(0x02)
				if err != nil {
					logrus.WithField("username", "czzhujian").Error(err)
				}
				msgHexStr := ""
				for _, item := range msg {
					msgHexStr += fmt.Sprintf("%02x", item)
				}
				logrus.WithField("username", "czzhujian").Info(msgHexStr)
				sockets[index].Close()
			} else {
				logrus.WithField("username", "czzhujian").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
			}
		} else {
			logrus.WithField("username", "czzhujian").Warn("[" + strconv.Itoa(dataEntity.DeviceID) + "]: 暂无数据")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func login(index int, device xphapi.Device, conn net.Conn) {
	var buffer bytes.Buffer
	var dataBuf bytes.Buffer
	buffer.WriteByte(0x15)
	buffer.WriteByte(0x01)
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)
	buffer.WriteByte(0x00)
	dataBuf.WriteString(`{"DeviceId":"` + strconv.Itoa(device.DeviceID) + `"}`)
	dataLen := uint16(dataBuf.Len())
	buffer.WriteByte(byte(dataLen))
	buffer.WriteByte(byte(dataLen >> 8))
	buffer.Write(dataBuf.Bytes())
	buffer.WriteByte(CheckSum(buffer.Bytes()[1:]))
	buffer.WriteByte(0x02)
	_, _ = conn.Write(buffer.Bytes())
	serialNumber[index]++
	//time.Sleep(500 * time.Millisecond)
	//reader := bufio.NewReader(conn)
	//msg, err := reader.ReadBytes(0x02)
	//if err != nil {
	//	logrus.WithField("username", "czzhujian").Error(err)
	//}
	//msgHexStr := ""
	//for _, item := range msg {
	//	msgHexStr += fmt.Sprintf("%02x", item)
	//}
	//logrus.WithField("username", "czzhujian").Info(msgHexStr)
}

func CheckSum(data []byte) uint8 {
	var sum uint8
	for _, item := range data {
		sum += item
	}
	return sum
}
