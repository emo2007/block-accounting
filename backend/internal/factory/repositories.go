package factory

import (
	"database/sql"
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/auth"
	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/cache"
	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/organizations"
	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/transactions"
	"github.com/emochka2007/block-accounting/internal/infrastructure/repository/users"
	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/redis/go-redis/v9"
)

func provideUsersRepository(db *sql.DB) users.Repository {
	return users.NewRepository(db)
}

func provideOrganizationsRepository(
	db *sql.DB,
	uRepo users.Repository,
) organizations.Repository {
	return organizations.NewRepository(db, uRepo)
}

func provideTxRepository(db *sql.DB, or organizations.Repository) transactions.Repository {
	return transactions.NewRepository(db, or)
}

func provideAuthRepository(db *sql.DB) auth.Repository {
	return auth.NewRepository(db)
}

func provideRedisConnection(c config.Config) (*redis.Client, func()) {
	r := redis.NewClient(&redis.Options{
		Addr:     c.DB.CacheHost,
		Username: c.DB.CacheUser,
		Password: c.DB.CacheSecret,
	})

	return r, func() { r.Close() }
}

func provideRedisCache(c *redis.Client, log *slog.Logger) cache.Cache {
	return cache.NewRedisCache(
		log.WithGroup("redis-cache"),
		c,
	)
}
