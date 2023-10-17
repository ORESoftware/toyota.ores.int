package proto

import (
	"context"
	"github.com/asynkron/protoactor-go/actor"
)

type MyGrpcActor struct{}

func (a *MyGrpcActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *MessageRequest:
		// Handle the gRPC request, e.g., send a message to another Actor
		response := &MessageResponse{Reply: "Received: " + msg.MessageData}
		context.Respond(response)
	}
}

type MyService struct {
	Actor *actor.RootContext
	Pid   *actor.PID
}

func (s *MyService) SendMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	//TODO implement me
	s.Actor.Send(s.Pid, req)
	return nil, nil
}

func (s *MyService) mustEmbedUnimplementedMyServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *MyService) MyMethod(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	response, err := s.Actor.RequestFuture(s.Pid, req, 5).Result()
	if err != nil {
		return nil, err
	}
	return response.(*MessageResponse), nil
}
