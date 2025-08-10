package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TypeEmail  = "task:email"
	TypeReport = "task:report"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type ReportPayload struct {
	Name string `json:"name"`
	At   time.Time `json:"at"`
}

type TaskContainer struct {
	Task *asynq.Task
}

// constructors
func NewEmailTask(to, subject, body string) *TaskContainer {
	p := EmailPayload{To: to, Subject: subject, Body: body}
	b, _ := json.Marshal(p)
	return &TaskContainer{Task: asynq.NewTask(TypeEmail, b)}
}

func NewReportTask() *TaskContainer {
	p := ReportPayload{Name: "daily_report", At: time.Now()}
	b, _ := json.Marshal(p)
	return &TaskContainer{Task: asynq.NewTask(TypeReport, b)}
}

// Handler struct
type TaskHandler struct{}

func NewTaskHandler() *TaskHandler { return &TaskHandler{} }

func (h *TaskHandler) HandleEmail(ctx context.Context, t *asynq.Task) error {
	var p EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Error().Err(err).Msg("invalid payload")
		return err
	}

	// Simulate sending email
	log.Info().Msgf("sending email to=%s subject=%s", p.To, p.Subject)
	if p.To == "fail@example.com" {
		return errors.New("simulated email failure")
	}
	return nil
}

func (h *TaskHandler) HandleReport(ctx context.Context, t *asynq.Task) error {
	var p ReportPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Error().Err(err).Msg("invalid payload")
		return err
	}
	
	log.Info().Msgf("generating report %s at %s", p.Name, p.At.Format(time.RFC3339))
	select {
	case <-time.After(2 * time.Second):
		// succeed
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
