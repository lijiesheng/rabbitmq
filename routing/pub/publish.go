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
	routingThree := lib.NewRabbitMQRouting("jiaohuanjiName", "router3")
	for i := 0; i <= 100; i++ {
		if i%3 == 0 {
			routingOne.PublishRouting("Hello kuteng one!  " + strconv.Itoa(i))
			fmt.Println("Hello kuteng one!  " + strconv.Itoa(i))
		} else if i %3 == 1 {
			routingTwo.PublishRouting("Hello kuteng two!  " + strconv.Itoa(i))
			fmt.Println("Hello kuteng two!  " + strconv.Itoa(i))
		} else if i %3 == 2 {
			routingThree.PublishRouting("Hello kuteng three!  " + strconv.Itoa(i))
			fmt.Println("Hello kuteng three!  " + strconv.Itoa(i))
		}
		time.Sleep(1 * time.Second)
	}
}
