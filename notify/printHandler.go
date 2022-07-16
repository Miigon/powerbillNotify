package notify

import "fmt"

type PrintHandler struct{}

func (h PrintHandler) Send(title string, content string) error {
	fmt.Printf("title: %s\ncontent: %s\n", title, content)
	return nil
}

func (h PrintHandler) String() string {
	return "printHandler"
}
