package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitClient struct {
	conn *amqp.Connection
	channel *amqp.Channel
	url string
}

func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}
	
	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %s", err)
	}
	
	return conn, channel, nil
}


func NewRabbitClient(connectionUrl string) (*RabbitClient, error) {
	conn, channel, err := newConnection(connectionUrl)
	if err != nil {
		return nil, err
	}
	
	return &RabbitClient{
		conn: conn,
		channel: channel,
		url: connectionUrl,
	}, nil
}

func (client *RabbitClient) ConsumeMessages(exchange, routingKey,queueName string) (<-chan amqp.Delivery, error) {

	err := client.channel.ExchangeDeclare(
		exchange, // name
		"direct",
		true, // durable
		true, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %s", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName, // name
		true, // durable
		true, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %s", err)
	}

	err = client.channel.QueueBind(
		queue.Name, // queue name
		routingKey, // routing key
		exchange, // exchange
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %s", err)
	}

	msgs, err := client.channel.Consume(
		queueName, // queue
		"goapp", // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume messages from queue: %s", err)
	}
	
	return msgs, nil
}

func (client *RabbitClient) PublishMessage(exchange, routingKey string, queueName string, message []byte) error {
	
	err := client.channel.ExchangeDeclare(
		exchange, // name
		"direct",
		true, // durable
		true, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %s", err)
	}

	queue, err := client.channel.QueueDeclare(
		queueName, // name
		true, // durable
		true, // delete when unused
		false, // exclusive
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}

	err = client.channel.QueueBind(
		queue.Name, // queue name
		routingKey, // routing key
		exchange, // exchange
		false, // no-wait
		nil, // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %s", err)
	}

	err = client.channel.Publish(
		exchange, // exchange
		routingKey, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %s", err)
	}

	return nil
}

func (client *RabbitClient) Close() {
	client.channel.Close()
	client.conn.Close()
}