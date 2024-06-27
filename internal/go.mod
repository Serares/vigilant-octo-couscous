module github.com/Bitstarz-eng/event-processing-challenge/internal

go 1.21.6

require (
	github.com/Bitstarz-eng/event-processing-challenge/pubsub v0.0.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/net v0.0.0-20220526153639-5463443f8c37
)

require (
	github.com/Bitstarz-eng/event-processing-challenge/exchanger v1.0.0
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)

replace github.com/Bitstarz-eng/event-processing-challenge/pubsub => ../pubsub

replace github.com/Bitstarz-eng/event-processing-challenge/exchanger => ../exchanger
