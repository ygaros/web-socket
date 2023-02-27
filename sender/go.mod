module sender

go 1.18

require (
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/chi/v5 v5.0.8
	github.com/mattn/go-sqlite3 v1.14.16 //it is used silently to provide sqlite3 driver
	go.uber.org/zap v1.24.0
)

require github.com/google/uuid v1.3.0

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)
