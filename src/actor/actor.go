package main

import (
	"github.com/IBM/sarama"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/wire"
)

// Message represents a gRPC message.
type Message struct {
	// Define your message structure here
}

// MyActor represents a ProtoActor actor to process gRPC messages and send them to Kafka.
type MyActor struct {
	saramaProducer sarama.AsyncProducer
}

// Receive handles incoming messages for the actor.
func (a *MyActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *Message:
		// Process the gRPC message and send it to Kafka.
		msgData := []byte("Message data") // Replace with your data
		producerMessage := &sarama.ProducerMessage{
			Topic: "your-topic", // Replace with your Kafka topic
			Value: sarama.StringEncoder(string(msgData)),
		}
		a.saramaProducer.Input() <- producerMessage
	}
}

// NewMyActor creates a new instance of MyActor with the provided Sarama producer.
func NewMyActor(producer sarama.AsyncProducer) *MyActor {
	return &MyActor{
		saramaProducer: producer,
	}
}

// Create your wire set to inject dependencies.
var actorSet = wire.NewSet(
	provideSaramaProducer,
	wire.Struct(new(MyActor), "saramaProducer"),
)

// provideSaramaProducer is a function to create a Sarama producer.
func provideSaramaProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for range producer.Successes() {
			// Handle successful Kafka deliveries here
		}
	}()

	return producer, nil
}
