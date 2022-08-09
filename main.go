package main

import (
	"flag"
	"fmt"
	"runtime/debug"
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

	defer func() {
		if err := recover(); err != nil {
			stacktrace := string(debug.Stack())
			msg := fmt.Sprintf("Failed to check powerbill or send notification email due to the occurence of following error: \n\n%s\n\n%s", err, stacktrace)
			fmt.Printf("!! PANIC AT %s: %s\n%s\n", time.Now().Format(time.RFC3339), err, stacktrace)
			fmt.Println("devEmails:", conf.Config.Email.DevEmails)
			mailerr := notify.
				MakeDefaultEmailHandler(conf.Config.Email.DevEmails).
				Send("ALERT: failed checking powerbill", msg)
			if mailerr != nil {
				fmt.Printf("An error occurred while sending alert email to devEmails: %s", err)
			} else {
				fmt.Println("An alert email has been sent to devEmails.")
			}
		}
	}()

	now := time.Now()
	bills := bill.GetPowerBill(now.Add(-time.Duration(conf.Config.Conditions.NDaysAveragePastUsage+1)*24*time.Hour), now)
	if len(bills) < conf.Config.Conditions.NDaysAveragePastUsage {
		panic(fmt.Sprintf("expected at least %d bills, got %d.", conf.Config.Conditions.NDaysAveragePastUsage, len(bills)))
	}
	for _, v := range bills {
		fmt.Printf("%+v\n", v)
	}
	notify.RegisterNotificationHandler(notify.PrintHandler{})
	notify.RegisterNotificationHandler(notify.MakeDefaultEmailHandler(conf.Config.Email.Recipients))

	if doSendNotification {
		notify.TryNotify(bills)
	} else {
		fmt.Println("run with -notify to send notification (if needed)")
	}
}
