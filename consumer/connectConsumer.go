package consumer

import (
	"encoding/json"
	"github.com/Wilddogmoto/RebbitTest/config"
	"github.com/Wilddogmoto/RebbitTest/logmq"
	"github.com/streadway/amqp"
	"os"
)

func ConnectConsumer() {
	var (
		err         error
		l           = logmq.InitializationLogger()
		connect     *amqp.Connection
		amqpChan    *amqp.Channel
		queue       amqp.Queue
		messageChan <-chan amqp.Delivery
	)
	connect, err = amqp.Dial(config.Config.AMQPConnectionURL)
	logmq.ResponseError(err, "bad connection")
	defer connect.Close()

	amqpChan, err = connect.Channel()
	logmq.ResponseError(err, "bad create channel")
	defer amqpChan.Close()

	queue, err = amqpChan.QueueDeclare("add",
		true,
		false,
		false,
		false,
		nil,
	)
	logmq.ResponseError(err, "error quere")

	err = amqpChan.Qos(1, 0, false)
	logmq.ResponseError(err, "bad config Qos")

	messageChan, err = amqpChan.Consume(queue.Name, "",
		false,
		false,
		false,
		false,
		nil,
	)
	logmq.ResponseError(err, "bad register consume")

	stopChan := make(chan bool)

	go func() {
		l.Infof("Consumer rdy, PID: %d", os.Getpid())
		for d := range messageChan {
			l.Infof("Received a message: %s", d.Body)
			addTask := &config.AddTask{}

			err = json.Unmarshal(d.Body, addTask)
			logmq.ResponseError(err, "error decod json")

			l.Infof("result of %d + %d is: %d", addTask.Number1, addTask.Number2, addTask.Number1+addTask.Number2)

			err = d.Ack(false)
			logmq.ResponseError(err, "Error acknowledging message")

			l.Info("Acknowledged message")
		}
	}()
	<-stopChan
}
