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
 * @Date 2022/8/10 12:45 下午
 **/
func main() {
	kutengOne := lib.NewRabbitMQTopic("exKutengTopic", "a.red.b")
	kutengTwo := lib.NewRabbitMQTopic("exKutengTopic", "b.c.yellow")
	for i := 0; i <= 100; i++ {
		kutengOne.PublicTopic("a.red.b" + strconv.Itoa(i))
		kutengTwo.PublicTopic("b.c.yellow" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
