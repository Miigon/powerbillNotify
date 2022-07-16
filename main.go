package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/miigon/powerbillNotify/bill"
	"github.com/miigon/powerbillNotify/conf"
	"github.com/miigon/powerbillNotify/notify"
)

var configFilePath string
var doSendNotification bool

func main() {
	flag.StringVar(&configFilePath, "config", "./config.yaml", "config file to use")
	flag.BoolVar(&doSendNotification, "notify", false, "whether send a notification or not when conditions are met")
	flag.Parse()

	fmt.Printf("Using config %s\n", configFilePath)
	conf.LoadConfig(configFilePath)

	fmt.Println("Fetching bills...")

	now := time.Now()
	bills := bill.GetPowerBill(now.Add(-time.Duration(conf.Config.Conditions.NDaysAveragePastUsage+1)*24*time.Hour), now)
	for _, v := range bills {
		fmt.Printf("%+v\n", v)
	}
	notify.RegisterNotificationHandler(notify.PrintHandler{})
	notify.RegisterNotificationHandler(notify.EmailHandler{
		SmtpServer: conf.Config.Email.Sender.SmtpServer,
		SmtpPort:   conf.Config.Email.Sender.SmtpPort,
		Username:   conf.Config.Email.Sender.Username,
		Password:   conf.Config.Email.Sender.Password,
		Recipients: conf.Config.Email.Recipients,
	})

	if doSendNotification {
		notify.TryNotify(bills)
	} else {
		fmt.Println("run with -notify to send notification (if needed)")
	}
}
