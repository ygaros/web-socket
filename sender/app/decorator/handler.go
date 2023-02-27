package decorator

import (
	"context"
	"fmt"
)

type CommandHandler[T any] interface {
	Handle(ctx context.Context, command T) error
}
type QueryHandler[Q, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

func NewCommandHadlerWithDefaultDecorators[T any](handler CommandHandler[T]) CommandHandler[T] {
	return loggingCommandHandler[T]{
		c: timeMetricsCommandHandler[T]{
			c: handler,
		},
	}
}
func NewQueryHandlerWithDefaultDecorators[Q, R any](handler QueryHandler[Q, R]) QueryHandler[Q, R] {
	return loggingQueryHandler[Q, R]{
		q: timeMetricsQueryHandler[Q, R]{
			q: handler,
		},
	}
}
func getName(action any) string {
	return fmt.Sprintf("%#T", action)
}
