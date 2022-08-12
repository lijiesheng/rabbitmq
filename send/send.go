package main

import (
	"fmt"
	"rabbitmq/lib"
	"strconv"
	"time"
)

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/8/10 5:31 下午
 * 生产者在交换机 normal_exchange 发送
 * 先打开 recieve1 和 dead
 * 在打开 send
 * 关闭 revieve1
 * 队列中存在的数据 + 死信消费的数据 + 正常消费的数据 = 生产的数据
 **/
func main() {
	router1 := lib.NewRabbitMQRouting("normal_exchange_1", "normal_key_1")
	for i := 0; i < 15000; i++ {
		message := "msg : " + time.Now().Format("2006年1月2日 15:04:05")
		fmt.Println(message)
		time.Sleep(2 * time.Second)
		router1.PublishRouting("Hello kuteng one!  " + strconv.Itoa(i))
	}
}
