package service

import (
	"context"
	"shajaro/actor/domain/kong"
)

type ConsumerService struct {
	Ctx    context.Context
	Konger kong.Konger
}

func NewConsumerService(ctx context.Context, konger kong.Konger) ConsumerService {
	return ConsumerService{
		Ctx:    ctx,
		Konger: konger,
	}
}
