package factory

import (
	"database/sql"
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/pkg/config"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/cache"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/organizations"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/users"
	"github.com/redis/go-redis/v9"
)

func provideUsersRepository(db *sql.DB) users.Repository {
	return users.NewRepository(db)
}

func provideOrganizationsRepository(db *sql.DB) organizations.Repository {
	return organizations.NewRepository(db)
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
