package queue

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	User              string
	Password          string
	Host              string
	Port              string
	Vhost             string
	ConsumerQueueName string
	ConsumerName      string
	AutoAck           bool
	Args              amqp.Table
	Channel           *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {

	rabbitMQArgs := amqp.Table{}
	rabbitMQArgs["x-dead-letter-exchange"] = os.Getenv("RABBITMQ_DLX")

	rabbitMQ := RabbitMQ{
		User:              os.Getenv("RABBITMQ_DEFAULT_USER"),
		Password:          os.Getenv("RABBITMQ_DEFAULT_PASS"),
		Host:              os.Getenv("RABBITMQ_DEFAULT_HOST"),
		Port:              os.Getenv("RABBITMQ_DEFAULT_PORT"),
		Vhost:             os.Getenv("RABBITMQ_DEFAULT_VHOST"),
		ConsumerQueueName: os.Getenv("RABBITMQ_CONSUMER_QUEUE_NAME"),
		ConsumerName:      os.Getenv("RABBITMQ_CONSUMER_NAME"),
		AutoAck:           false,
		Args:              rabbitMQArgs,
	}

	return &rabbitMQ
}

func (r *RabbitMQ) Connect() *amqp.Channel {
	dns := "amqp://" + r.User + ":" + r.Password + "@" + r.Host + ":" + r.Port + r.Vhost

	conn, err := amqp.Dial(dns)
	failOnError(err, "Failed to connect to RAbbitMQ")

	r.Channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel with RAbbitMQ")

	return r.Channel
}

func (r *RabbitMQ) Consume(messageChannel chan amqp.Delivery) {

	q, err := r.Channel.QueueDeclare(
		r.ConsumerQueueName, //name
		true,                //durable
		false,               //delete when used
		false,               //exclusive
		false,               //no-wait
		r.Args,              //arguments
	)
	failOnError(err, "Failed to declare a queue")

	incomingMessage, err := r.Channel.Consume(
		q.Name,         //queue
		r.ConsumerName, //consumer
		r.AutoAck,      //auto-ack
		false,          //exclusive
		false,          //no-local
		false,          //no-wait
		nil,            //arguments
	)
	failOnError(err, "Failed do register a consumer")

	go func() {
		for message := range incomingMessage {
			log.Println("Incoming new message")
			messageChannel <- message
		}
		log.Println("RabbitMQ channel closed")
		close(messageChannel)
	}()
}

func (r *RabbitMQ) Notify(message string, contentType string, exchange string, routingKey string) error {
	err := r.Channel.Publish(
		exchange,   //exchange
		routingKey, //routing key
		false,      //mandatory
		false,      //immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		},
	)

	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("#{msg}: #{err}")
	}
}
