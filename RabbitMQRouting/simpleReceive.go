package main

import "daniel-rabbitmq/RabbitMQ"

func main() {
	rabbitMQ := RabbitMQ.NewRabbitMQRouting("daniel_routing","routingKey_1")
	rabbitMQ.ConsumeRouting()
}
