package rabbitservice

import (
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

type PubSubService struct {
	ConnectionString string
	Logger           *slog.Logger
	Channel          *amqp091.Channel
}

func NewPubService(connStr string, log *slog.Logger) *PubSubService {
	return &PubSubService{
		ConnectionString: connStr,
		Logger:           log.WithGroup("PubService"),
	}
}

func (s *PubSubService) Init() error {
	conn, err := amqp091.Dial(s.ConnectionString)
	if err != nil {
		s.Logger.Error("Failed to connect to RabbitMQ", slog.Any("error", err))
		return err
	}

	ch, err := conn.Channel()
	// ch.Qos(
	// 	1,     // prefetch count
	// 	0,     // prefetch size
	// 	false, // global
	// )
	if err != nil {
		s.Logger.Error("Failed to open a channel", slog.Any("error", err))
		return err
	}
	err = ch.ExchangeDeclare(
		"casino", // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		s.Logger.Error("Failed to declare an exchange", slog.Any("error", err))
		return err
	}
	s.Channel = ch
	s.Logger.Info("Initiated pub service success")
	return nil
}

func (s *PubSubService) Publish(msg interface{}) error {
	if s.Channel == nil {
		err := s.Init()
		if err != nil {
			return err
		}
	}

	toJson, err := json.Marshal(msg)
	if err != nil {
		s.Logger.Error("Failed to marshal message", slog.Any("error", err))
		return err
	}

	err = s.Channel.Publish(
		"casino", // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        toJson,
		},
	)

	if err != nil {
		s.Logger.Error("Failed to publish a message", slog.Any("error", err))
		return err
	}
	s.Logger.Info("Published event: ", "message:", msg)
	return nil
}

func (s *PubSubService) Close() {
	s.Channel.Close()
}

func (s *PubSubService) Subscribe() (<-chan amqp091.Delivery, error) {
	if s.Channel == nil {
		err := s.Init()
		if err != nil {
			return nil, err
		}
	}

	q, err := s.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		s.Logger.Error("Failed to declare a queue", slog.Any("error", err))
		return nil, err
	}

	err = s.Channel.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		"casino", // exchange
		false,
		nil,
	)

	if err != nil {
		s.Logger.Error("Failed to bind a queue", slog.Any("error", err))
		return nil, err
	}

	msgs, err := s.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		s.Logger.Error("Failed to register a consumer", slog.Any("error", err))
		return nil, err
	}

	return msgs, nil
}
