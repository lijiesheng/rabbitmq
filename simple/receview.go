package main

import "rabbitmq/lib"

func main() {
	rabbitMQ := lib.NewRabbitMQSimple("task_queue_333333")
	rabbitMQ.ConsumeSimple()
}
