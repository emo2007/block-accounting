package domain

import (
	"github.com/emochka2007/block-accounting/internal/interface/rest/domain/hal"
)

type Organization struct {
	*hal.Resource
	Id        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
}
