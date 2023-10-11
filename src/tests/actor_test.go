package tests

import (
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
)

func TestMyActor_ProcessMessage(t *testing.T) {
	system := actor.NewActorSystem()
	myActor, err := system.ActorOf(actor.PropsFromProducer(func() actor.Actor { return &MyActor{} }))
	assert.NoError(t, err)

	// Define your test cases here
	// Send messages to the actor and assert the results
}
