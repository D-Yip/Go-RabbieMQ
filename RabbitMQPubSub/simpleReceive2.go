package main

import "daniel-rabbitmq/RabbitMQ"

func main() {
	rabbitMQ := RabbitMQ.NewRabbitMQPubSub("daniel_exchange")
	rabbitMQ.ConsumeSub()
}
