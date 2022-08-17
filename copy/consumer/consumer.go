package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var (
	exchange_name string = "j_exch_head"
	routing_key string = "jkey"
	etype = amqp.ExchangeHeaders
	queue_name = "j_queue"
)

/**
 * @Description
 * @Author lijiesheng
 * @Date 2022/8/17 4:32 下午
 **/
func main() {
	conn, _ := amqp.Dial("amqp://kuteng:ljs024816@127.0.0.1:5672/kuteng")
	ch, _:= conn.Channel()
	ch.QueueDeclare(queue_name, true, false, true, false, nil)
	ch.QueueBind(queue_name, routing_key, exchange_name, false, amqp.Table{
		"mail" : "470047253@qq.com",
	})
	ch.ExchangeDeclare(exchange_name, etype, true, false, false, false,nil)
	// 监听
	msgs, _ := ch.Consume(queue_name, "", true, false, false, false, nil)
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			//println("tset")
			log.Printf(" [x] %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
