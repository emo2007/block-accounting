package controllers

import (
	"context"
	"log/slog"
	"net/http"
)

type PingController interface {
	HandlePing(ctx context.Context, req *http.Request, w http.ResponseWriter) error
}

type pingController struct {
	log *slog.Logger
}

func NewPingController(
	log *slog.Logger,
) PingController {
	return &pingController{
		log: log,
	}
}

func (c *pingController) HandlePing(ctx context.Context, req *http.Request, w http.ResponseWriter) error {
	_, err := w.Write([]byte("pong"))

	return err
}
