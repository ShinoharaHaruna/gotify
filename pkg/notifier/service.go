package notifier

import "context"

type Service interface {
	Send(ctx context.Context, title, message string) error
	Name() string
}
