package main

import (
	"daniel-rabbitmq/RabbitMQ"
	"fmt"
	"strconv"
)

func main() {
	rabbitMQ := RabbitMQ.NewRabbitMQRouting("daniel_routing","routingKey_1")
	rabbitMQTwo := RabbitMQ.NewRabbitMQRouting("daniel_routing","routingKey_2")
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			rabbitMQ.PublicRouting("Hello Daniel ! " + strconv.Itoa(i))
		} else {
			rabbitMQTwo.PublicRouting("Hello Daniel ! " + strconv.Itoa(i))
		}
	}
	fmt.Println("sent out !!!")
}
