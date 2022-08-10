package main

import "rabbitmq/lib"

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/8/10 10:49 上午
 **/
func main() {
	kuteng_one := lib.NewRabbitMQRouting("jiaohuanjiName", "router2")
	kuteng_one.RecieveRouting()
}
