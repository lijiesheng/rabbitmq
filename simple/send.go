package main

import (
	"fmt"
	"rabbitmq/lib"
	"strconv"
)

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/8/8 6:46 下午
 **/
func main() {
	rabbitMQ := lib.NewRabbitMQSimple("test1234")
	for i := 0; i < 20; i++ {
		str := strconv.Itoa(i*i) + ",  hello ljs"
		rabbitMQ.PublishSimple(str) // 发送字符串
	}
	fmt.Println("发送成功")
}
