package main

import (
	"context"
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
	switch msg := ctx.Message().(type) {
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

type MyGrpcActor struct{}

func (a *MyGrpcActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *proto.MessageRequest:
		// Handle the gRPC request, e.g., send a message to another actor
		response := &proto.MessageResponse{Reply: "Received: " + msg.Message}
		context.Respond(response)
	}
}

type MyService struct {
	pid *actor.PID
}

func (s *MyService) mustEmbedUnimplementedMyServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *MyService) SendMessage(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MyService) MyMethod(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	response, err := s.pid.RequestFuture(req, 5).Result()
	if err != nil {
		return nil, err
	}
	return response.(*proto.MessageResponse), nil
}

func main() {
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
	proto.RegisterMyServiceServer(grpcServer, &MyService{pid: pid})

	// Start the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
