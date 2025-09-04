package main

import (
	"email-service/config"
	"email-service/internal/logger"
	"email-service/internal/rabbitmq"
	"log"

	"go.uber.org/zap"
)



func main() {
	if err := config.SetupEnvVar(); err != nil {
		log.Fatal(err.Error())
	}
	if err := logger.StartLogger(); err != nil {
		logger.ZapLogger.Fatal(err.Error())
	}
	if err := rabbitmq.ConnectToRabbitMQ(); err != nil {
		logger.ZapLogger.Error(err.Error(), zap.String("function", "rabbitmq.connecttorabbitmq"))
	}


}