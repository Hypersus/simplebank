package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskVerifyEmail(
		ctx context.Context,
		payload *PayloadVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(opts asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(opts)
	return &RedisTaskDistributor{
		client: client,
	}
}
