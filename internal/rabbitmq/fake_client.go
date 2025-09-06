package rabbitmq

import (
	"email-service/internal/dto"
	"email-service/internal/logger"
	"strings"
	"time"
)

type fakeClient struct {
}

func (c *fakeClient) sendEmail(input dto.EmailDTO) error {
	logger.ZapLogger.Info("[fake] " + "sending email to " + strings.Join(input.To, "") + " subject " + input.Subject)
	return  nil
}

func (c *fakeClient) consumerEmail() {
	msg := make(chan dto.EmailDTO)
	go func() {
		for email := range msg {
			err := c.sendEmail(email)
			if err != nil {

			}
		}
	}()
	go func() {
		for {
			msg <- dto.EmailDTO{
				To:      []string{"user@example.com"},
				Subject: "test",
				Text:    "test message",
			}
			time.Sleep(time.Second * 2)
			
		} 
	}()
	for range msg {}
}