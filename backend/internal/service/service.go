package service

import (
	"context"
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/interface/rest"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	Run(ctx context.Context) error
	Stop()
}

type service struct {
	log  *slog.Logger
	rest *rest.Server
}

func NewService(
	log *slog.Logger,
	rest *rest.Server,
) Service {
	return &service{
		log: log,
	}
}

func (s *service) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.rest.Serve(ctx)
	})

	return g.Wait()
}

func (s *service) Stop() {
}
