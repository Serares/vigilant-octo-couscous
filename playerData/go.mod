module github.com/Bitstarz-eng/event-processing-challenge/playerData

go 1.21.6

require (
	github.com/Bitstarz-eng/event-processing-challenge/internal v0.0.0-00010101000000-000000000000
	github.com/Bitstarz-eng/event-processing-challenge/pubsub v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect

replace github.com/Bitstarz-eng/event-processing-challenge/pubsub => ../pubsub

replace github.com/Bitstarz-eng/event-processing-challenge/internal => ../internal
