package controllers

import (
	"log/slog"
	"net/http"
)

type PingController interface {
	Ping(w http.ResponseWriter, req *http.Request) ([]byte, error)
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

func (c *pingController) Ping(w http.ResponseWriter, req *http.Request) ([]byte, error) {
	return []byte("pong"), nil
}
