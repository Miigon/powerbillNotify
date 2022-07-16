package notify

import "fmt"

type NotificationHandler interface {
	Send(title string, content string) error
	String() string
}

var handlers []NotificationHandler

func RegisterNotificationHandler(handler NotificationHandler) {
	handlers = append(handlers, handler)
}

func SendNotification(title string, content string) {
	fmt.Printf("\n<<<< Sending notification: %s\n", title)
	for i, h := range handlers {
		fmt.Printf("[%d]: %s\n", i, h.String())
		err := h.Send(title, content)
		if err != nil {
			fmt.Println("error occurred: " + err.Error())
		}
	}
	fmt.Printf("====\n\n")
}
