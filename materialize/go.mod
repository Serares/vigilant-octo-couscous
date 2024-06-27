module github.com/Bitstarz-eng/event-processing-challenge/materialize

go 1.21.6

require (
	github.com/Bitstarz-eng/event-processing-challenge/pubsub v1.0.0
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/Bitstarz-eng/event-processing-challenge/internal v1.0.0
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)

replace github.com/Bitstarz-eng/event-processing-challenge/pubsub => ../pubsub

replace github.com/Bitstarz-eng/event-processing-challenge/internal => ../internal
