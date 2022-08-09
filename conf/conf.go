package conf

import (
	"os"

	"gopkg.in/yaml.v3"
)

type NotificationTemplate struct {
	Title   string `yaml:"title"`
	Content string `yaml:"content"`
}

var Config struct {
	ServerQueryUrl string `yaml:"serverQueryUrl"`
	ClientField    string `yaml:"clientField"` // depending on campus location
	RoomId         string `yaml:"roomId"`
	RoomName       string `yaml:"roomName"`
	Building       string `yaml:"building"`
	Conditions     struct {
		NotifyNearRunningOut  bool    `yaml:"notifyNearRunningOut"`
		NDaysAveragePastUsage int     `yaml:"nDaysAveragePastUsage"`
		DaysLeftThreshold     int     `yaml:"daysLeftThreshold"`
		NotifyAbnormallyHigh  bool    `yaml:"notifyAbnormallyHigh"`
		HighUsageThreshold    float64 `yaml:"highUsageThreshold"`
	} `yaml:"conditions"`
	Templates struct {
		NearRunningOut NotificationTemplate `yaml:"nearRunningOut"`
		AbnormallyHigh NotificationTemplate `yaml:"abnormallyHigh"`
	} `yaml:"templates"`
	Email struct {
		Sender struct {
			SmtpServer string `yaml:"smtpServer"`
			SmtpPort   int    `yaml:"smtpPort"`
			Username   string `yaml:"username"`
			Password   string `yaml:"password"`
		} `yaml:"sender"`
		Recipients []string `yaml:"recipients"`
		DevEmails  []string `yaml:"devEmails"`
	} `yaml:"email"`
}

func LoadConfig(configfile string) {
	data, err := os.ReadFile(configfile)
	if err != nil {
		panic("failed loading config: " + err.Error())
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		panic("failed unmarshaling config: " + err.Error())
	}
}
