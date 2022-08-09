package main

import "rabbitmq/lib"
// 公平的分发
func main() {
	rabbitMQ := lib.NewRabbitMQSimple("task_queue_333333")
	rabbitMQ.ConsumeSimple()
}
