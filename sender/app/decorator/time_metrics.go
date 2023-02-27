package decorator

import (
	"context"
	"log"
	"time"
)

type timeMetricsCommandHandler[C any] struct {
	c CommandHandler[C]
}
type timeMetricsQueryHandler[Q any, R any] struct {
	q QueryHandler[Q, R]
}

func (handler timeMetricsCommandHandler[C]) Handle(ctx context.Context, cmd C) (err error) {
	start := time.Now()

	defer func() {
		diff := time.Now().Sub(start)
		if err != nil {
			log.Println(err)
			log.Printf("Failure after %v\n", diff)
		} else {
			log.Printf("Success after %v\n", diff)
		}
	}()
	return handler.c.Handle(ctx, cmd)
}
func (handler timeMetricsQueryHandler[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	start := time.Now()

	defer func() {
		diff := time.Now().Sub(start)
		if err != nil {
			log.Println(err)
			log.Printf("Failure after %v\n", diff)
		} else {
			log.Printf("Success after %v\n", diff)
		}
	}()
	return handler.q.Handle(ctx, query)
}
