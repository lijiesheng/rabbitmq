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
 **/
func main() {
	router1 := lib.NewRabbitMQRouting("normal_exchange_1", "normal_key_1")
	for i := 0; i < 200; i++ {
		message := "msg : " + time.Now().Format("2006年1月2日 15:04:05")
		fmt.Println(message)
		time.Sleep(2 * time.Second)
		router1.PublishRouting("Hello kuteng one!  " + strconv.Itoa(i))
	}
}
