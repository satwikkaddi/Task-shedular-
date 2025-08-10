package asynqclient

import (
	"github.com/hibiken/asynq"
)

func RedisConnOpt(addr string) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{Addr: addr}
}
