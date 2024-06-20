package queue

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/ctxmeta"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/google/uuid"
)

func TestJobMarshal(t *testing.T) {
	ctx := ctxmeta.UserContext(context.Background(), &models.User{
		ID:   uuid.New(),
		Name: "kjdsfhkjfg",
		Credentails: &models.UserCredentials{
			Email: "jkdfhgls",
		},
		PK:        []byte("1234567890qwertyuiop"),
		Bip39Seed: []byte("poiuytrewq0987654321"),
		Mnemonic:  "mnemonic mnemonic mnemonicccc",
		Activated: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	ctx = ctxmeta.OrganizationIdContext(ctx, uuid.New())

	job := &Job{
		ID:             "123",
		IdempotencyKey: "123",
		Context:        ctx,
		Payload:        &JobDeployMultisig{OwnersPubKeys: []string{"sdfdf", "sdfsd"}, Confirmations: 2},
		CreatedAt:      time.Now().UnixMilli(),
	}

	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	t.Log(string(data))

	var job2 *Job = new(Job)

	if err := json.Unmarshal(data, job2); err != nil {
		t.Fatalf("err: %s", err.Error())
	}

	t.Logf("%+v", job2)
}
