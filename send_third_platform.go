package main

import (
	"github.com/sirupsen/logrus"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	logrus.Info("------ send_third_platform start ------")
}

func main() {

	go StandardStart()

	go GuokongStart()

	go HnhebiStart()

	go KaifengStart()

	go NanyangStart()

	go Nanyang2019Start()

	go PingdingshanStart()

	go ZhatucangStart()

	go ZhumadianStart()

	go GkgriddingStart()

	go HnsjsbStart()

	go RenkeStart()

	select {}

}
