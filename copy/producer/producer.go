package main



import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

var (
	exchange_name string = "j_exch_head"
	routing_key string = "jkey"
	etype = amqp.ExchangeHeaders
)

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/8/17 10:50 上午
 **/
func main() {
	conn, _ := amqp.Dial("amqp://kuteng:ljs024816@127.0.0.1:5672/kuteng")
	ch, _:= conn.Channel()
	body := "Hello World!! " + time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(body)
	ch.Confirm(false)    // broker 会返回一个 confirm.ok

	// 声明交换器
	ch.ExchangeDeclare(exchange_name, etype, true, false, false, true, nil)

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	// 处理方法
	defer confirmOne(confirms)
	ch.Publish(exchange_name, routing_key, false, false,amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body ),
		Headers:     amqp.Table{"x-match": "any", "mail": "470047253@qq.com", "author": "Jhonny"},
	},)
	defer conn.Close()
	defer ch.Close()
}

// 消息确认
func confirmOne(confirms <- chan amqp.Confirmation) {
	if confirmed := <-confirms; confirmed.Ack {
		fmt.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		fmt.Printf("confirmed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}