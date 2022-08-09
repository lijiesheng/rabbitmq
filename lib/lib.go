package lib

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
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

// 创建简单模式下的 RabbitMQ 实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	// 创建 RabbitMQ 实例
	rabbitmq := NewRabbitMQ(queueName, "", "") // 交换机和 key 的值都默认是空
	var err error
	// 获取 connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	// 获取 channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 简单模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	// 申请队列，如果队列不存在，就制动创建，存在则跳过
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化 【如果服务器重启，队列中的数据就没有了】
		true,     // true 是持久化
		//是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否具有阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 调用 channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
	amqp.Publishing{
		DeliveryMode: amqp.Persistent,   // 持久
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// simple 模式下消费者
// 带应答信号，只有消费了，队列才会删除
func (r *RabbitMQ) ConsumeSimple() {
	//1.申请队列，如果队列不存在会自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		true,   // true 是持久化
		//是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否具有阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 接收消息
	// 返回的 msg 是一个 chanel, 是不停的接收来自客户端的消息
	msgs, err := r.channel.Consume(q.Name, // queue
		//用来区分多个消费者
		"", // consumer
		//是否自动应答
		//true, // 不要手动应答
		false, // 需要手动应该   这样
		//是否独有
		false, // exclusive
		//设置为true，表示 不能将同一个Conenction中生产者发送的消息传递给这个Connection中 的消费者
		false, // no-local
		//列是否阻塞
		false, // no-wait
		nil)   // args)

	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs { // 这是一个死循环
			//消息逻辑处理，可以自行设计逻辑
			log.Printf("Received a message: %s", d.Body)
			time.Sleep(2 * time.Second)
			fmt.Println("完成任务")
			d.Ack(false)
		}
		fmt.Println("消息处理完成")
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}