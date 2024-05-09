package controllers

import "log/slog"

type OrganizationsController interface {
}

type organizationsController struct {
	log *slog.Logger
}
