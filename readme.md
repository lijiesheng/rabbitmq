# 文件运行是 simple 文件夹下的生产者和消费者

配置如下：
        消费者：
        非持久化队列，非自动删除
        q, err := r.channel.QueueDeclare(r.QueueName, false, false, false, false, nil, )
        需要手动应答 【autoAck = false ,需要手动应答，推荐这样】
        msgs, err := r.channel.Consume(q.Name, "", false, false, false, false, nil)
        手动应答
        go func() {
          for d := range msgs { // 这是一个死循环
          //消息逻辑处理，可以自行设计逻辑
          log.Printf("Received a message: %s", d.Body)
          time.Sleep(2 * time.Second)
           d.Ack(false)
        }}()

      生成者代码
        非持久化队列，非自动删除
        _, err := r.channel.QueueDeclare(r.QueueName, false, false, false, false, nil, )


1、启动消费者

2、启动生产者 [发送了 200 个消息]

测试1 : 消费者关闭之后再次启动
消费者仍然可以接着之前没有消费完的消费  【因为有应答机制】

测试2 ：rabbitmq 服务器重启
队列中的数据没有了，消费者自然拿不到消息



改成消息持计划  [可以参考 message-durability 分支]
修改的地方
    队列持久化 [包括生产者和消费者]
        
    消息持久化 
        amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			DeliveryMode : 2,
		})
        
    一个新的队列

1、启动消费者

2、启动生产者 [发送了 200 个消息]

测试1 : 消费者关闭之后再次启动
        消费者仍然可以接着之前没有消费完的消费  【因为有应答机制】

测试2 ：rabbitmq 服务器重启
        队列中的数据存在
