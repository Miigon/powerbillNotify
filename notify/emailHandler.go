package notify

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type EmailHandler struct {
	SmtpServer string
	SmtpPort   int
	Username   string
	Password   string
	Recipients []string
}

func (h EmailHandler) Send(title string, content string) error {
	return smtp.SendMail(
		h.SmtpServer+":"+strconv.Itoa(h.SmtpPort),
		smtp.PlainAuth("", h.Username, h.Password, h.SmtpServer),
		h.Username,
		h.Recipients,
		[]byte(fmt.Sprintf("From: powerbillNotify <%s>\nSubject: %s\n%s", h.Username, title, content)))
}

func (h EmailHandler) String() string {
	return fmt.Sprintf("emailHandler<%s> to %v", h.Username, h.Recipients)
}
