package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
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
	Context        context.Context
	Payload        any
	CreatedAt      int64
}

type job struct {
	ID             string      `json:"id"`
	IdempotencyKey string      `json:"idempotency_key"`
	Context        *JobContext `json:"context"`
	Type           string      `json:"_type"`
	Payload        []byte      `json:"payload"`
	CreatedAt      int64       `json:"created_at"`
}

func (j *Job) MarshalJSON() ([]byte, error) {
	payload, err := json.Marshal(j.Payload)
	if err != nil {
		return nil, fmt.Errorf("error marshal job payload. %w", err)
	}

	ja := &job{
		ID:             j.ID,
		IdempotencyKey: j.IdempotencyKey,
		Context:        newOutgoingCoutext(j.Context),
		Type:           jobType(j.Payload),
		Payload:        payload,
		CreatedAt:      j.CreatedAt,
	}

	return json.Marshal(ja)
}

// TODO: fix this memory overhead
func (j *Job) UnmarshalJSON(data []byte) error {
	ja := &job{}

	err := json.Unmarshal(data, ja)
	if err != nil {
		return err
	}

	j.Payload, err = payloadByType(ja.Type, ja.Payload)
	if err != nil {
		return err
	}

	j.ID = ja.ID
	j.IdempotencyKey = ja.IdempotencyKey
	j.Context = ja.Context
	j.CreatedAt = ja.CreatedAt

	return nil
}

func payloadByType(t string, data []byte) (any, error) {
	switch t {
	case "job_deploy_multisig":
		var dm JobDeployMultisig

		if err := json.Unmarshal(data, &dm); err != nil {
			return nil, err
		}

		return &dm, nil
	default:
		return nil, fmt.Errorf("error unknown job type")
	}
}

func jobType(job any) string {
	switch job.(type) {
	case *JobDeployMultisig:
		return "job_deploy_multisig"
	default:
		return ""
	}
}

type JobContext struct {
	Parent *JobContext `json:"_parent"`
	Key    any         `json:"_key"`
	Val    any         `json:"_value"`
}

func (c *JobContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (c *JobContext) Done() <-chan struct{} {
	return nil
}

func (c *JobContext) Err() error {
	return nil
}

func (c *JobContext) Value(key any) any {
	if c.Key == key {
		return c.Val
	}

	return c.Parent.Value(key)
}

func newOutgoingCoutext(ctx context.Context) *JobContext {
	var jobCtx *JobContext = new(JobContext)

	lastFrame := jobCtx

	if user, err := ctxmeta.User(ctx); err == nil {
		lastFrame.Key = ctxmeta.UserContextKey
		lastFrame.Val = user
		lastFrame.Parent = new(JobContext)

		lastFrame = lastFrame.Parent
	}

	if orgId, err := ctxmeta.OrganizationId(ctx); err == nil {
		lastFrame.Key = ctxmeta.OrganizationIdContextKey
		lastFrame.Val = orgId
		// lastFrame.Parent = new(JobContext)

		// lastFrame = lastFrame.Parent
	}

	return jobCtx
}
