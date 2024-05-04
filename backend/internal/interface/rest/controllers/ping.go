package controllers

import (
	"log/slog"
	"net/http"
)

type PingController interface {
	Ping(w http.ResponseWriter, req *http.Request) error
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

func (c *pingController) Ping(w http.ResponseWriter, req *http.Request) error {
	_, err := w.Write([]byte("pong"))

	return err
}
