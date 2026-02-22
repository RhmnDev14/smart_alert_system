package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// TaskProducer defines the interface for publishing background tasks globally
type TaskProducer interface {
	Publish(taskType string, payload []byte, retries int) error
	Close() error
}

// TaskConsumer defines the interface for acting as a worker processing background tasks globally
type TaskConsumer interface {
	RegisterHandler(taskType string, handler func(context.Context, []byte) error)
	Start() error
	Stop()
}

// asynqProducer implementation
type asynqProducer struct {
	client *asynq.Client
}

// NewTaskProducer creates a new global Asynq task producer (publisher)
func NewTaskProducer(host, port, password string, db int) TaskProducer {
	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	}
	client := asynq.NewClient(redisOpt)
	return &asynqProducer{client: client}
}

func (p *asynqProducer) Publish(taskType string, payload []byte, retries int) error {
	task := asynq.NewTask(taskType, payload, asynq.MaxRetry(retries))
	info, err := p.client.Enqueue(task)
	if err != nil {
		return err
	}
	log.Printf("queue: Enqueued task %s (ID: %s)", taskType, info.ID)
	return nil
}

func (p *asynqProducer) Close() error {
	return p.client.Close()
}

// asynqConsumer implementation
type asynqConsumer struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

// NewTaskConsumer creates a new global Asynq worker server (consumer)
func NewTaskConsumer(host, port, password string, db int, concurrency int) TaskConsumer {
	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	}
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"default": 10,
			},
		},
	)
	mux := asynq.NewServeMux()

	return &asynqConsumer{
		server: srv,
		mux:    mux,
	}
}

func (c *asynqConsumer) RegisterHandler(taskType string, handler func(context.Context, []byte) error) {
	// Wrap generic handler into Asynq's HandlerFunc
	asynqHandler := func(ctx context.Context, t *asynq.Task) error {
		return handler(ctx, t.Payload())
	}
	c.mux.HandleFunc(taskType, asynqHandler)
}

func (c *asynqConsumer) Start() error {
	log.Println("queue: Starting Task Consumer worker...")
	return c.server.Start(c.mux)
}

func (c *asynqConsumer) Stop() {
	log.Println("queue: Stopping Task Consumer worker...")
	c.server.Shutdown()
}
