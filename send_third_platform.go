package main

import (
	"github.com/sirupsen/logrus"
	"whxph.com/send-third-platform/czzhujian"
	"whxph.com/send-third-platform/gkgridding"
	"whxph.com/send-third-platform/guokong"
	"whxph.com/send-third-platform/hnhebi"
	"whxph.com/send-third-platform/hnsjsb"
	"whxph.com/send-third-platform/kaifeng"
	"whxph.com/send-third-platform/linghui"
	"whxph.com/send-third-platform/nanyang"
	"whxph.com/send-third-platform/nanyang2019"
	"whxph.com/send-third-platform/pingdingshan"
	"whxph.com/send-third-platform/renke"
	"whxph.com/send-third-platform/standard"
	"whxph.com/send-third-platform/zhatucang"
	"whxph.com/send-third-platform/zhumadian"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	logrus.Info("------ send_third_platform start ------")
}

func main() {

	// go standard.Start()

	// go guokong.Start()

	// go hnhebi.Start()

	// go kaifeng.Start()

	// go nanyang.Start()

	// go nanyang2019.Start()

	// go pingdingshan.Start()

	// go zhatucang.Start()

	// go zhumadian.Start()

	// go gkgridding.Start()

	// go hnsjsb.Start()

	// go renke.Start()

	// go linghui.Start()

	go czzhujian.Start()

	select {}

}
