package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailDelivery = "email: deliver"
)

type EmailDeliveryPayload struct {
	UserID     uint
	TemplateID string
}

func NewEmailDeliveryTask(userID uint, tmplID string) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{UserID: userID, TemplateID: tmplID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%d, template_id=%s", payload.UserID, payload.TemplateID)
	// Email delivery code ...
	return nil
}
