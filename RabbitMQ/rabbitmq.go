package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//	url	格式	amqp://username:password@ip:port/vhost
const MQURL="amqp://daniel:694975@192.168.92.112:5672/daniel"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	//	队列名称
	QueueName string
	//	交换机
	Exchange string
	//	key
	Key string
	//	连接信息
	MqUrl string
}

//	创建RabbitMQ结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitMQ := &RabbitMQ{
		QueueName:queueName,
		Exchange:exchange,
		Key:key,
		MqUrl:MQURL}
	var err error
	//	创建rabbitMQ连接
	rabbitMQ.conn,err = amqp.Dial(rabbitMQ.MqUrl)
	rabbitMQ.failOnError(err,"创建连接错误")
	rabbitMQ.channel, err = rabbitMQ.conn.Channel()
	rabbitMQ.failOnError(err,"获取channel失败")
	return rabbitMQ
}

//	断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

//	错误处理函数
func (r *RabbitMQ) failOnError(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s",message,err)
		panic(fmt.Sprintf("%s:%s",message,err))
	}
}

//	创建简单模式下rabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName,"","")
}

//	订阅模式创建RabbitMQ实例
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	return NewRabbitMQ("",exchange,"")
}

//	路由模式
//	创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	return NewRabbitMQ("",exchangeName,routingKey)
}

//	简单模式下发布消息
func (r *RabbitMQ) PublicSimple(message string) {
	//	申请队列,如果队列不存在会自动创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,	//	是否持久化
		false,	//	是否自动删除
		false,	//	是否具有排他性
		false,	//	是否阻塞
		nil,	//	额外参数
	)
	if err != nil {
		fmt.Println(err)
	}
	//发送消息到队列
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//	如果为true，根据exchange类型和route key规则，如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//	如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		},
	)
}

//	简单模式下消费消息
func (r *RabbitMQ) ConsumeSimple() {
	//	申请队列,如果队列不存在会自动创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,	//	是否持久化
		false,	//	是否自动删除
		false,	//	是否具有排他性
		false,	//	是否阻塞
		nil,	//	额外参数
	)
	if err != nil {
		fmt.Println(err)
	}
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",	//	用来区分多个消费者
		true,	//	是否自动应答
		false,	//	是否具有排他性
		false,	//	如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,	//	是否阻塞
		nil,	//	额外参数
	)
	if err != nil {
		fmt.Println(err)
	}
	forever := make(chan bool)
	//	启用协程处理消息
	go func() {
		for d :=range msgs {
			log.Printf("Received a message: %s",d.Body)
		}
	}()

	log.Printf("[*] waiting for messages, to exit press ctrl + c")
	<-forever
}

//	订阅模式生产
func (r *RabbitMQ) PublicPub(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err,"Failed to declare an exchange")
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(message),
		})
}

//	订阅模式消费
func (r *RabbitMQ) ConsumeSub() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err,"Failed to declare an exchange")
	//	试探性创建队列，队列名称不写
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(err,"Failed to declare a queue")

	//	绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)

	//	消费消息
	msgs, err := r.channel.Consume(
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
		for d := range msgs{
			log.Printf("Received a message : %s",d.Body)
		}
	}()

	fmt.Println("退出请按 Ctrl + c")
	<-forever
}

//	路由模式发送消息
func (r *RabbitMQ) PublicRouting(message string) {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err,"Failed to declare an exchange")
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(message),
		})
}

//	路由模式接收消息
func (r *RabbitMQ) ConsumeRouting() {
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnError(err,"Failed to declare an exchange")
	//	试探性创建队列，队列名称不写
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnError(err,"Failed to declare a queue")

	//	绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	//	消费消息
	msgs, err := r.channel.Consume(
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
		for d := range msgs{
			log.Printf("Received a message : %s",d.Body)
		}
	}()

	fmt.Println("退出请按 Ctrl + c")
	<-forever
}