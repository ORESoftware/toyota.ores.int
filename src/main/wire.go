// wire.go

//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Shopify/sarama"
	"github.com/google/wire"
)

func InitializeMyActor() (*MyActor, error) {
	wire.Build(actorSet)
	return nil, nil
}

func InitializeSaramaProducer() (sarama.AsyncProducer, error) {
	wire.Build(provideSaramaProducer)
	return nil, nil
}
