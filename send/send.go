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
 * 先打开 recieve1 和 dead 消费者
 * 打开 send 生产者
 * 关闭 recieve1 生产者 可以看到 dead 消费者在 10s 之后开始消费
 **/
func main() {
	router1 := lib.NewRabbitMQRouting("normal_exchange_1", "normal_key_1")
	for i := 0; i < 20; i++ {
		message := "msg : " + time.Now().Format("2006年1月2日 15:04:05")
		fmt.Println(message)
		time.Sleep(2 * time.Second)
		router1.PublishRouting("Hello kuteng one!  " + strconv.Itoa(i))
	}
}
