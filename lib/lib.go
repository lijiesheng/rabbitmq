package lib

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// 连接信息
//url := "amqp://账号:密码@host:port/vhost
const MQURL = "amqp://kuteng:ljs024816@127.0.0.1:5672/kuteng"

// rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	// bind Key 名称
	Key string
	// 连接信息
	Mqurl string
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
}

// 使用完rabbitMQ 后， 断开 channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数，mysql redis mongo 都建议自己定义一种错误日志
// 其他错误逻辑都在这个里面用
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Printf("%s:%s", message, err)
	}
}

// 路由创建
// 创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	// 1、创建 RabbitMQ实例
	rabbitMQ := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 2、获取 connection
	rabbitMQ.conn, err = amqp.Dial(rabbitMQ.Mqurl)
	rabbitMQ.failOnErr(err, "failed to connect rabbitmq!")
	// 3、获取 channel
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnErr(err, "failed to open a channel")
	return rabbitMQ
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

//话题模式
//创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	// 获取 connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 话题模式发送消息
func (r *RabbitMQ) PublicTopic(message string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		//要改成topic
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")
	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		//要设置
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

//话题模式接受消息
//要注意key,规则
//其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 kuteng.* 表示匹配 kuteng.hello, kuteng.hello.one需要用kuteng.#才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	//1.试探性创建交换机
	r.channel.ExchangeDeclare(r.Exchange, "topic", true, false, false, false, nil)
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")
	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		r.Key,
		r.Exchange,
		false,
		nil)
	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}

