package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/emochka2007/block-accounting/internal/interface/rest"
)

type Service interface {
	Run(ctx context.Context) error
	Stop()
}

type ServiceImpl struct {
	log  *slog.Logger
	rest *rest.Server
}

func NewService(
	log *slog.Logger,
	rest *rest.Server,
) Service {
	return &ServiceImpl{
		log:  log,
		rest: rest,
	}
}

func (s *ServiceImpl) Run(ctx context.Context) error {
	errch := make(chan error)

	defer s.rest.Close()

	go func() {
		defer func() {
			close(errch)
		}()

		errch <- s.rest.Serve(ctx)
	}()

	select {
	case <-ctx.Done():
		s.log.Info("shutting down service")

		return nil
	case err := <-errch:
		return fmt.Errorf("error at service runtime. %w", err)
	}
}

func (s *ServiceImpl) Stop() {

}
