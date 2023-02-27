package decorator

import (
	"context"
	"log"
)

type loggingCommandHandler[C any] struct {
	c CommandHandler[C]
}
type loggingQueryHandler[Q any, R any] struct {
	q QueryHandler[Q, R]
}

func (handler loggingCommandHandler[C]) Handle(ctx context.Context, cmd C) (err error) {
	name := getName(cmd)
	log.Printf("Executing %s, body: %#v", name, cmd)
	defer func() {
		if err != nil {
			log.Println(err)
			log.Println("Command execution failed")
		} else {
			log.Println("Command executed successfully")
		}
	}()
	return handler.c.Handle(ctx, cmd)
}
func (handler loggingQueryHandler[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	name := getName(query)
	log.Printf("Executing %s, body: %#v", name, query)
	defer func() {
		if err != nil {
			log.Println(err)
			log.Println("Query execution failed")
		} else {
			log.Println("Query executed successfully")
		}
	}()
	return handler.q.Handle(ctx, query)
}
