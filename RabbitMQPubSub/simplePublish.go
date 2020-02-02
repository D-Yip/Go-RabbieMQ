package main

import (
	"daniel-rabbitmq/RabbitMQ"
	"fmt"
	"strconv"
)

func main() {
	rabbitMQ := RabbitMQ.NewRabbitMQPubSub("daniel_exchange")
	for i := 0; i < 100; i++ {
		rabbitMQ.PublicPub("Hello Daniel ! " + strconv.Itoa(i))
	}
	fmt.Println("sent out !!!")
}
