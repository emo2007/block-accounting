package domain

type Transaction struct {
	Id             string  `json:"id"`
	Description    string  `json:"description"`
	OrganizationId string  `json:"organization_id"`
	CreatedBy      string  `json:"created_by"`
	Amount         float64 `json:"amount"`
	ToAddr         []byte  `json:"to"`
	MaxFeeAllowed  float64 `json:"max_fee_allowed"`
	Deadline       int64   `json:"deadline,omitempty"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
	ConfirmedAt    int64   `json:"confirmed_at,omitempty"`
	CancelledAt    int64   `json:"cancelled_at,omitempty"`
	CommitedAt     int64   `json:"commited_at,omitempty"`
}
