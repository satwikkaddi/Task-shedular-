package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"taskscheduler/internal/config"
	"taskscheduler/internal/tasks"
	"taskscheduler/pkg/asynqclient"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.LoadConfig()
	r := asynqclient.RedisConnOpt(cfg.RedisAddr)

	srv := asynq.NewServer(
		r,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	handler := tasks.NewTaskHandler()

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmail, handler.HandleEmail)
	mux.HandleFunc(tasks.TypeReport, handler.HandleReport)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatal().Err(err).Msg("asynq server failed")
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	srv.Shutdown()
	log.Info().Msg("worker stopped")
}
