package ctxmeta

import (
	"context"
	"fmt"

	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/google/uuid"
)

type ContextKey string

var (
	UserContextKey                    = ContextKey("user")
	OrganizationIdContextKey          = ContextKey("org-id")
	OrganizationParticipantContextKey = ContextKey("org-participant")
)

func UserContext(parent context.Context, user *models.User) context.Context {
	return context.WithValue(parent, UserContextKey, user)
}

func User(ctx context.Context) (*models.User, error) {
	if user, ok := ctx.Value(UserContextKey).(*models.User); ok {
		return user, nil
	}

	return nil, fmt.Errorf("error user not passed in context")
}

func OrganizationParticipantContext(
	parent context.Context,
	participant models.OrganizationParticipant,
) context.Context {
	return context.WithValue(parent, OrganizationParticipantContextKey, participant)
}

func OrganizationIdContext(parent context.Context, id uuid.UUID) context.Context {
	return context.WithValue(parent, OrganizationIdContextKey, id)
}

func OrganizationId(ctx context.Context) (uuid.UUID, error) {
	if id, ok := ctx.Value(OrganizationIdContextKey).(uuid.UUID); ok {
		return id, nil
	}

	return uuid.Nil, fmt.Errorf("error organization id not passed in context")
}
