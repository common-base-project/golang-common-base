package connection

import (
	"github.com/streadway/amqp"
	"golang-common-base/pkg/logger"
)

//实现异步通信
//系统解耦
//流量削峰

//引入消息队列带来的延迟问题
//增加了系统的复杂度
//可能产生数据不一致的问题

var ConnMq *amqp.Connection
var ChannelMq *amqp.Channel

// 连接 rabbitMQ
func OpenMq() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ConnMq = conn

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()
	ChannelMq = channel
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Fatalf("%s: %s", msg, err)
	}
}
