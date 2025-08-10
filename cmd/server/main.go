package main

import (
	"context"
	"fmt"
	"time"

	"taskscheduler/internal/config"
	"taskscheduler/internal/tasks"
	"taskscheduler/pkg/asynqclient"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Create Redis connection for client and scheduler
	r := asynqclient.RedisConnOpt(cfg.RedisAddr)

	// Simple demo enqueuer
	client := asynq.NewClient(r)
	defer client.Close()

	// Enqueue a one-off task delayed by 30s
	t := tasks.NewEmailTask("welcome@example.com", "Welcome!", "Thanks for joining.")
	info, err := client.Enqueue(t.Task, asynq.ProcessIn(30*time.Second), asynq.MaxRetry(5))
	if err != nil {
		log.Fatal().Err(err).Msg("could not enqueue task")
	}
	log.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// Start an Asynq Scheduler for a recurring job (every minute)
	scheduler := asynq.NewScheduler(r, &asynq.SchedulerOpts{
		Location: time.UTC,
	})
	_, err = scheduler.Register("@every 1m", tasks.NewReportTask().Task)
	if err != nil {
		log.Fatal().Err(err).Msg("could not register scheduled task")
	}

	ctx := context.Background()
	if err := scheduler.Run(); err != nil {
		log.Fatal().Err(err).Msg("scheduler stopped")
	}

	<-ctx.Done()
	fmt.Println("server exiting")
}
