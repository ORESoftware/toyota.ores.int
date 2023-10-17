package main

import (
	"github.com/IBM/sarama"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"log"
	"net"
	proto "toyota.ores.int/src/proto/src"
)

// Message represents a gRPC message.
type Message struct {
	// Define your message structure here
}

// Define your ProtoActor actor to process gRPC messages.
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

func InitializeMyActor() (*MyActor, error) {
	wire.Build(actorSet)
	return nil, nil
}

func InitializeSaramaProducer() (sarama.AsyncProducer, error) {
	wire.Build(provideSaramaProducer)
	return nil, nil
}

func main() {

	log.Println("he here")
	// Create a ProtoActor system and start the actor.
	system := actor.NewActorSystem()
	context := system.Root
	props := actor.PropsFromProducer(func() actor.Actor { return &MyActor{} })
	pid, _ := context.SpawnNamed(props, "myActor")

	// Send a message to the actor
	context.RequestFuture(pid, "Hello, Actor!", 1)

	// Create a gRPC server
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register your gRPC service
	proto.RegisterMyServiceServer(grpcServer, &proto.MyService{Actor: context, Pid: pid})

	if _, err := InitializeMyActor(); err != nil {
		log.Fatalf("Failed to init actor wire: %v", err)
	}

	if _, err := InitializeSaramaProducer(); err != nil {
		log.Fatalf("Failed to init saram: %v", err)
	}

	// Start the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
