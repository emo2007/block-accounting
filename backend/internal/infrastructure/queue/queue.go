package queue

import (
	"context"
	"fmt"
)

type QueueDriver interface {
	Put(ctx context.Context, job any) error
	Pop(ctx context.Context) (any, error)
}

type Queue[T any] struct {
	driver QueueDriver
}

func NewWithDriver[T any](
	driver QueueDriver,
) *Queue[T] {
	return &Queue[T]{
		driver: driver,
	}
}

func (q *Queue[T]) Put(
	ctx context.Context,
	job T,
) error {
	return q.driver.Put(ctx, job)
}

func (q *Queue[T]) Pop(ctx context.Context) (*T, error) {
	job, err := q.driver.Pop(ctx)
	if err != nil {
		return nil, fmt.Errorf("queue: error pop a job from the queue. %w", err)
	}

	if t, ok := job.(T); ok {
		return &t, nil
	}

	return nil, fmt.Errorf("queue: error unexpected job type")
}

type Job struct {
	ID             string
	IdempotencyKey string
	Payload        any
	CreatedAt      int64
}
