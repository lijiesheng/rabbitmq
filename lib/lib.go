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



func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// 订阅模式创建 rabbitmq 实例
func NewRabbitMQPubSub(queueName,exchangeName string) *RabbitMQ {
	// 创建RabbitMQ 实例
	rabbitmq := NewRabbitMQ(queueName, exchangeName, "")
	var err error
	// 获取 connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to conect rabbitmq!")
	// 获取 channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}



// 发送 fanout
func (r *RabbitMQ) PublishPub(message string) {
	// 1、尝试创建交换机
	fmt.Println("r.Exchange==>", r.Exchange)
	err := r.channel.ExchangeDeclare(
		r.Exchange,   // 使用命名的交换器
		"fanout",     // 交换器类型
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
		)
	r.failOnErr(err, "Failed to declare ")
	// 2、发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",   // routing key
		false, //
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		},
	)
}

// 消费端代码
func (r *RabbitMQ) RecieveSub() {
	fmt.Println("r.Exchange===>",r.Exchange)
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,   // 使用命名的交换器
		"fanout",     // 交换器类型
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	r.failOnErr(err, "failed to declare an exchange")
	// 试探性创建队列，注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"",   // 队列名称，随机产生队列名称
		false,
		false,
		true,
		false,
		nil,
	)

	// 交换机绑定到一个队列
	fmt.Println("r.Exchange ==>", r.Exchange, "q.Name==>", q.Name)
	err = r.channel.QueueBind(
		q.Name,   // 队列的名称, 如果不指定会生成一个随机的队列
		"",
		r.Exchange,   // 交换器名称
		false,
		nil,
		)

	r.failOnErr(err, "failed to declare a queue")
	err = r.channel.QueueBind(q.Name, "", r.Exchange, false, nil)
	// 消费模型
	messages, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message : %s", d.Body)
		}
	}()
	fmt.Println("请按出 ctrl + c\n")
	<-forever
}
