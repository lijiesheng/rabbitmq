package main

import "rabbitmq/lib"

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/1/24 5:49 下午
 **/
func main() {
	kutengOne := lib.NewRabbitMQTopic("exKutengTopic", "*.red.*")
	kutengOne.RecieveTopic()
}
