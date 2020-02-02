package main

import (
	"daniel-rabbitmq/RabbitMQ"
	"fmt"
	"strconv"
)

func main() {
	rabbitMQ := RabbitMQ.NewRabbitMQSimple("daniel")
	for i := 0; i < 100; i++ {
		rabbitMQ.PublicSimple("Hello Daniel ! " + strconv.Itoa(i))
	}
	fmt.Println("sent out !!!")
}
