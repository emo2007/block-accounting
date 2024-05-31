package domain

type Organization struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	CreatedAt uint64 `json:"created_at"`
	UpdatedAt uint64 `json:"updated_at"`
}
