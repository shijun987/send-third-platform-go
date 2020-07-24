package xphapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Token token
type Token struct {
	Token      string `json:"token"`
	Expiration int    `json:"expiration"`
	Message    string `json:"message"`
	UserID     int    `json:"userID"`
}

// User 用户信息
type User struct {
	Username int      `json:"username"`
	UserType string   `json:"userType"`
	Devices  []Device `json:"devices"`
}

// Device 设备信息
type Device struct {
	DeviceID   int    `json:"facId"`
	DeviceName string `json:"facName"`
}

// DataEntity 数据
type DataEntity struct {
	DeviceID int      `json:"deviceId"`
	Entity   []Entity `json:"entity"`
}

// Entity 实体
type Entity struct {
	Datetime string `json:"datetime"`
	EUnit    string `json:"eUnit"`
	EValue   string `json:"eValue"`
	EKey     string `json:"eKey"`
	EName    string `json:"eName"`
	ENum     string `json:"eNum"`
}

// GetToken 获取token
func GetToken(username, password string) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	loginParam := map[string]string{"username": username, "password": password}
	jsonStr, _ := json.Marshal(loginParam)
	resp, err := client.Post("http://115.28.187.9:8005/login", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	var token Token
	_ = json.Unmarshal(result, &token)
	return token.Token
}

// GetDevices 获取设备ID
func GetDevices(username, token string) []Device {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "http://115.28.187.9:8005/user/"+username, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("token", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var user User
	_ = json.Unmarshal(body, &user)
	return user.Devices
}
