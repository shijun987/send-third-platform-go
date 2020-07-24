package utils

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// SoapResult wsdl解析结果
type SoapResult struct {
	XMLName xml.Name
	Body    struct {
		XMLName          xml.Name
		SaveYCJCResponse struct {
			XMLName xml.Name
			Out     []string `xml:"out"`
		} `xml:"saveYCJCResponse"`
	}
}

// Invoke 调用wsdl
func Invoke(url, value string) {

	client := &http.Client{}
	var payload []byte
	if url == "http://42.236.61.105:8686/xncjk/services/SaveYCJCService" {
		payload = []byte(strings.TrimSpace(`
		<?xml version="1.0" encoding="utf-16"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
			<soap:Body>
				<saveYCJC xmlns="http://api.controller.myself.com/">
					<elements xmlns="">` + value + `</elements>
				</saveYCJC>
			</soap:Body>
		</soap:Envelope>
	`))
	} else {
		payload = []byte(strings.TrimSpace(`
		<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
			<soap:Body>
				<saveYCJC xmlns="http://service.tblycjc.webservice.client.dekn.com.cn">
					<in0>` + value + `</in0>
				</saveYCJC>
			</soap:Body>
		</soap:Envelope>
	`))
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result := new(SoapResult)
	err = xml.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return
	}
	logrus.Info(result.Body.SaveYCJCResponse.Out)
}
