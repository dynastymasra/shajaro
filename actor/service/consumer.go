package service

import (
	"context"
	"shajaro/actor/domain/kong"
)

type ConsumerService struct {
	Ctx      context.Context
	Consumer kong.Consumer
}

func NewConsumerService(ctx context.Context, consumer kong.Consumer) ConsumerService {
	return ConsumerService{
		Ctx:      ctx,
		Consumer: consumer,
	}
}
