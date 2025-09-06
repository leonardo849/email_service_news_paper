package rabbitmq

import (
	"email-service/internal/dto"
	"email-service/internal/logger"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"go.uber.org/zap"
)

const smtpHost string = "smtp.gmail.com"
const smtpPort string = "587"

func (c *client) sendEmail(input dto.EmailDTO) error {
	from := os.Getenv("SERVICE_EMAIL")
	password := os.Getenv("SERVICE_PASSWORD")
	if from == "" || password == "" {
		err := fmt.Errorf("from email or password is empty")
		logger.ZapLogger.Fatal(err.Error(), zap.String("function", "c.sendEmail"))
		return err
	}
	toHeader := strings.Join(input.To, ",")


	message := []byte(
		"From: " + from + "\r\n" +
		"To: " + toHeader + "\r\n" +
		"Subject: " + input.Subject + "\r\n" + "\r\n" + input.Text + "\r\n",
	)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, input.To, message); err != nil {
		logger.ZapLogger.Error(err.Error(), zap.String("function", "c.sendEmail"))
		return err
	}

	return nil
}

func (c *client) consumerEmail() {
	const exchangeName = "email_direct"
	const queueName = "email_service_queue"

	if err := c.ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		logger.ZapLogger.Fatal(err.Error(), zap.String("function", "client.ConsumerEmail"))
	}

	q, err := c.ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.ZapLogger.Fatal(err.Error(), zap.String("function", "client.ConsumerEmail"))
	}

	err = c.ch.QueueBind(
		q.Name,
		key,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		logger.ZapLogger.Fatal(err.Error(), zap.String("function", "client.ConsumerEmail"))
	}
	msgs, err := c.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.ZapLogger.Fatal(err.Error(), zap.String("function", "client.ConsumerEmail"))
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var email dto.EmailDTO
			logger.ZapLogger.Info("one more email")
			err := json.Unmarshal(d.Body, &email)
			if err != nil {
				logger.ZapLogger.Error("error in json.unmarshal email", zap.Error(err))
			}
			if err :=c.sendEmail(email); err != nil {
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
			
		}
	}()
	<-forever
}
