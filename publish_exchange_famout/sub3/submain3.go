package main

import "rabbitmq/lib"

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/8/9 2:07 下午
 **/
func main() {
	rabbitMQ := lib.NewRabbitMQPubSub("logs","pushlish_ljs1")
	rabbitMQ.RecieveSub()
}
