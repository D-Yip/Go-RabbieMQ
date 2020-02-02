package main

import "daniel-rabbitmq/RabbitMQ"

func main() {
	rabbitMQ := RabbitMQ.NewRabbitMQSimple("daniel")
	rabbitMQ.ConsumeSimple()
}
