package services

import (
	"context"
	"fmt"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

type TelegramService struct {
	service *telegram.Telegram
	receivers []int64
}

func NewTelegram(apiKey string, receivers []int64) (*TelegramService, error) {
	tg, err := telegram.New(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram service: %w", err)
	}

	service := &TelegramService{
		service:   tg,
		receivers: receivers,
	}

	for _, receiver := range receivers {
		tg.AddReceivers(receiver)
	}

	return service, nil
}

func (t *TelegramService) Send(ctx context.Context, title, message string) error {
	notifier := notify.New()
	notifier.UseServices(t.service)
	return notifier.Send(ctx, title, message)
}

func (t *TelegramService) Name() string {
	return "telegram"
}
