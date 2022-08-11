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

func (r * RabbitMQ) RecieveRoutingList(str_list[]string) {
	// 1、试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, // 交换器名字
		"direct",   // 交换器类型
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil)        // arguments
	r.failOnErr(err, "Failed to declare an exchange")

	// 2、试探性创建队列，队列不要写名字
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	failOnError(err, "Failed to declare a queue")

	str_list = append(str_list, r.Key)

	for _, v := range str_list {
		// 绑定多个
		err = r.channel.QueueBind(
			q.Name,
			v,
			r.Exchange,
			false,
			nil)
	}

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

func (r *RabbitMQ) RecieveDeadRouting() {
	// 1、试探性创建交换机
	err := r.channel.ExchangeDeclare(r.Exchange,"direct", true, false, false, false, nil)
	r.failOnErr(err, "创建死信交换机失败")
	// 2、创建死信队列
	q, err := r.channel.QueueDeclare("", true, false, false, false, nil)   // todo 这里可能要变
	// 3、队列绑定（将队列、routing-key、交换机三者绑定到一起）
	err = r.channel.QueueBind(q.Name, r.Key, r.Exchange, false, nil)
	failOnError(err, "绑定失败")
	//消费消息
	messges, err := r.channel.Consume(q.Name, "", true, false, false, false, nil )
	forever := make(chan bool)
	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}

func (r *RabbitMQ) RecieveRouting() {
	// 1、试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange, // 交换器名字
		"direct",   // 交换器类型
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil)        // arguments
	r.failOnErr(err, "Failed to declare an exchange")

	// 2、试探性创建队列，队列不要写名字
	q, err := r.channel.QueueDeclare(
		"normal_queue_1",
		true,
		false,
		false,
		false,
		amqp.Table{
			//"x-message-ttl": 10000, // 消息过期时间（队列级别）,毫秒 , 这里也可以不写，在生成消息的时候指定   Expiration: "10000", // 消息过期时间,毫秒
			"x-dead-letter-exchange":"dead-exchange_1", // 指定死信交换机
			"x-dead-letter-routing-key": "dead-key_1", // 指定死信routing-key
		})
	failOnError(err, "Failed to declare a queue")

	// 队列 交换器 routing-keying 三者绑定在一起
	fmt.Println(q.Name, r.Key, r.Exchange)
	err = r.channel.QueueBind(
		q.Name,
		//需要绑定key
		r.Key,
		r.Exchange,
		false,
		nil)

	fmt.Println("q.Name = ", q.Name , ", r.Key = ", r.Key, " , r.Exchange = ", r.Exchange)
	//消费消息
	messges, err := r.channel.Consume(q.Name, "", false, false, false, false, nil, )
	forever := make(chan bool)
	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}

func (r *RabbitMQ) RecieveRoutingDeadQueue() {
	// 1、创建交换器
	err := r.channel.ExchangeDeclare(r.Exchange, "direct", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")
	// 2、试探性创建队列，队列不要写名字
	q, err := r.channel.QueueDeclare("dead-queue_1", true, false, true, false, nil)
	// 3. 交换机 队列 routing-key 绑定
	err = r.channel.QueueBind("dead-queue_1", r.Key, r.Exchange, false, nil)
	// 4、消费
	//消费消息
	fmt.Println("死信消费，", q.Name, " ,", r.Key, ", ", r.Exchange)
	messges, err := r.channel.Consume("dead-queue_1", "", false, false, false, false, nil, )
	forever := make(chan bool)
	go func() {
		for d := range messges {
			fmt.Println("死信消费 11111111")
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}


// 路由模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	fmt.Println("dsa==>",r.Exchange, r.Key)
	// 1、尝试创建交换机
	err := r.channel.ExchangeDeclare(r.Exchange, "direct", true, false, false, false, nil)
	r.failOnErr(err, "Failed to declare an exchange")
	// 2、发送消息
	r.channel.Publish(r.Exchange, r.Key,     false, false,
		amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
		Expiration: "10000", // 消息过期时间,毫秒
	})
}
