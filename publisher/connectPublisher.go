package publisher

import (
	"encoding/json"
	"github.com/Wilddogmoto/RebbitTest/config"
	"github.com/Wilddogmoto/RebbitTest/logmq"
	"github.com/streadway/amqp"
	"math/rand"
)

func ConnectPublisher() {

	var (
		err      error
		body     []byte
		l        = logmq.InitializationLogger()
		connect  *amqp.Connection
		amqpChan *amqp.Channel
		queue    amqp.Queue
	)

	if connect, err = amqp.Dial(config.Config.AMQPConnectionURL); err != nil {
		l.Errorf("bad connection:%v", err)
	}
	defer connect.Close()

	if amqpChan, err = connect.Channel(); err != nil {
		l.Errorf("bad create channel:%v", err)
	}

	defer amqpChan.Close()

	if queue, err = amqpChan.QueueDeclare("add",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		l.Errorf("error queue:%v", err)
	}
	//rand.Seed(time.Now().UnixNano())

	addTask := config.AddTask{Number1: rand.Intn(999), Number2: rand.Intn(999)}
	if body, err = json.Marshal(addTask); err != nil {
		l.Errorf("error encod json:%v", err)
	}

	if err = amqpChan.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "texst/plain",
		Body:         body,
	}); err != nil {
		l.Errorf("error publishing message: %v", err)
	}
	l.Infof("addTask:%d + %d", addTask.Number1, addTask.Number2)
}
