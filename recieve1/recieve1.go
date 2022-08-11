package main

import "rabbitmq/lib"

// 创建队列，交换机
// 正常消费
func main() {
	dead := lib.NewRabbitMQRouting("normal_exchange_1", "normal_key_1")
	dead.RecieveRouting()
}
