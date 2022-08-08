package main

import "rabbitmq/lib"

func main() {
	rabbitMQ := lib.NewRabbitMQSimple("test1234")
	rabbitMQ.ConsumeSimple()
}
