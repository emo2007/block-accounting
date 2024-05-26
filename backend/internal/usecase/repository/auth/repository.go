package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/emochka2007/block-accounting/internal/pkg/models"
	sqltools "github.com/emochka2007/block-accounting/internal/pkg/sqlutils"
	"github.com/google/uuid"
)

type AddTokenParams struct {
	UserId uuid.UUID

	Token          string
	TokenExpiredAt time.Time

	RefreshToken          string
	RefreshTokenExpiredAt time.Time

	CreatedAt time.Time

	RemoteAddr string
}

type GetTokenParams struct {
	UserId       uuid.UUID
	Token        string
	RefreshToken string
}

type RefreshTokenParams struct {
	UserId uuid.UUID

	OldToken       string
	Token          string
	TokenExpiredAt time.Time

	OldRefreshToken       string
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
}

type AccessToken struct {
	UserId uuid.UUID

	Token          string
	TokenExpiredAt time.Time

	RefreshToken          string
	RefreshTokenExpiredAt time.Time

	CreatedAt time.Time
}

type Repository interface {
	AddToken(ctx context.Context, params AddTokenParams) error
	GetTokens(ctx context.Context, params GetTokenParams) (*AccessToken, error)
	RefreshToken(ctx context.Context, params RefreshTokenParams) error

	AddInvite(ctx context.Context, params AddInviteParams) error
	MarkAsUsedLink(ctx context.Context, linkHash string, usedAt time.Time) error
}

type repositorySQL struct {
	db *sql.DB
}

func (r *repositorySQL) AddToken(ctx context.Context, params AddTokenParams) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Insert("access_tokens").
			Columns(
				"user_id",
				"token",
				"refresh_token",
				"token_expired_at",
				"refresh_token_expired_at",
				"remote_addr",
			).
			Values(
				params.UserId,
				params.Token,
				params.RefreshToken,
				params.TokenExpiredAt,
				params.RefreshTokenExpiredAt,
				params.RemoteAddr,
			).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error add tokens. %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *repositorySQL) RefreshToken(ctx context.Context, params RefreshTokenParams) error {
	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		updateQuery := sq.Update("access_tokens").
			SetMap(sq.Eq{
				"token":                    params.Token,
				"refresh_token":            params.RefreshToken,
				"token_expired_at":         params.TokenExpiredAt,
				"refresh_token_expired_at": params.RefreshTokenExpiredAt,
			}).
			Where(sq.Eq{
				"user_id":       params.UserId,
				"token":         params.OldToken,
				"refresh_token": params.OldRefreshToken,
			}).PlaceholderFormat(sq.Dollar)

		if _, err := updateQuery.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error update tokens. %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *repositorySQL) GetTokens(ctx context.Context, params GetTokenParams) (*AccessToken, error) {
	var token *AccessToken = new(AccessToken)

	if err := sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Select(
			"user_id",
			"token",
			"token_expired_at",
			"refresh_token",
			"refresh_token_expired_at",
			"created_at",
		).From("access_tokens").
			Where(sq.Eq{
				"token":   params.Token,
				"user_id": params.UserId,
			}).PlaceholderFormat(sq.Dollar)

		if params.RefreshToken != "" {
			query = query.Where(sq.Eq{
				"refresh_token": params.RefreshToken,
			})
		}

		rows, err := query.RunWith(r.Conn(ctx)).QueryContext(ctx)
		if err != nil {
			return fmt.Errorf("error fetch token from database. %w", err)
		}

		defer func() {
			if cErr := rows.Close(); cErr != nil {
				err = errors.Join(fmt.Errorf("error close database rows. %w", cErr), err)
			}
		}()

		for rows.Next() {
			if err := rows.Scan(
				&token.UserId,
				&token.Token,
				&token.TokenExpiredAt,
				&token.RefreshToken,
				&token.RefreshTokenExpiredAt,
				&token.CreatedAt,
			); err != nil {
				return fmt.Errorf("error scan row. %w", err)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return token, nil
}

type AddInviteParams struct {
	LinkHash       string
	OrganizationID uuid.UUID
	CreatedBy      models.User
	CreatedAt      time.Time
	ExpiredAt      time.Time
}

func (r *repositorySQL) AddInvite(
	ctx context.Context,
	params AddInviteParams,
) error {
	return sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Insert("invites").Columns(
			"link_hash",
			"organization_id",
			"created_by",
			"created_at",
			"expired_at",
		).Values(
			params.LinkHash,
			params.OrganizationID,
			params.CreatedBy.Id(),
			params.CreatedAt,
			params.ExpiredAt,
		).PlaceholderFormat(sq.Dollar)

		if _, err := query.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error add invite link. %w", err)
		}

		return nil
	})
}

func (r *repositorySQL) MarkAsUsedLink(
	ctx context.Context,
	linkHash string,
	usedAt time.Time,
) error {
	return sqltools.Transaction(ctx, r.db, func(ctx context.Context) error {
		query := sq.Select("expired_at").From("invites").Where(sq.Eq{
			"link_hash": linkHash,
		}).Limit(1).PlaceholderFormat(sq.Dollar)

		var expAt time.Time

		if err := query.RunWith(r.Conn(ctx)).QueryRowContext(ctx).Scan(&expAt); err != nil {
			return fmt.Errorf("error fetch expiration date from database. %w", err)
		}

		if expAt.After(time.Now()) {
			return ErrorInviteLinkExpired
		}

		updateQuery := sq.Update("invites").SetMap(sq.Eq{
			"used_at": usedAt,
		})

		if _, err := updateQuery.RunWith(r.Conn(ctx)).ExecContext(ctx); err != nil {
			return fmt.Errorf("error add invite link. %w", err)
		}

		return nil
	})
}

func NewRepository(db *sql.DB) Repository {
	return &repositorySQL{
		db: db,
	}
}

func (s *repositorySQL) Conn(ctx context.Context) sqltools.DBTX {
	if tx, ok := ctx.Value(sqltools.TxCtxKey).(*sql.Tx); ok {
		return tx
	}

	return s.db
}
