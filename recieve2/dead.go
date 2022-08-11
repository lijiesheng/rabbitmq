package main

import "rabbitmq/lib"

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/8/10 6:51 下午
 **/
func main() {
	dead := lib.NewRabbitMQRouting("dead-exchange_1", "dead-key_1")
	dead.RecieveRoutingDeadQueue()
}
