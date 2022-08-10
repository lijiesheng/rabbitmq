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
 * @Date 2022/8/9 5:42 下午
 **/
func main() {
	routingOne := lib.NewRabbitMQRouting("jiaohuanjiName", "router1")
	routingTwo := lib.NewRabbitMQRouting("jiaohuanjiName", "router2")
	for i := 0; i <= 100; i++ {
		if i%2 == 0 {
			routingOne.PublishRouting("Hello kuteng one!  " + strconv.Itoa(i))
		} else {
			routingTwo.PublishRouting("Hello kuteng two!  " + strconv.Itoa(i))
		}
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
