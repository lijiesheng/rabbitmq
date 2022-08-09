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
 * @Date 2022/8/9 2:01 下午
 **/
func main() {
	// 交换器的名称是 pushlish_ljs
	rabbitMQ := lib.NewRabbitMQPubSub("logs","pushlish_ljs_temp")
	for i := 0; i < 100; i++ {
		rabbitMQ.PublishPub("订阅模式 "+ strconv.Itoa(i) + " 条数据")
		fmt.Println("订阅模式 "+ strconv.Itoa(i) + " 条数据")
		time.Sleep(3 * time.Second)
	}
}
