package domain

// Generic

type Collection[T any] struct {
	Items      []T        `json:"items,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	NextCursor string `json:"next_cursor,omitempty"`
	TotalItems uint32 `json:"total_items,omitempty"`
}

// Auth related DTO's

type JoinRequest struct {
	Name       string `json:"name,omitempty"`
	Credentals struct {
		Email    string `json:"email,omitempty"`
		Phone    string `json:"phone,omitempty"`
		Telegram string `json:"telegram,omitempty"`
	} `json:"credentals,omitempty"`

	Mnemonic string `json:"mnemonic"`
}

type JoinResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Mnemonic string `json:"mnemonic"`
}

type RefreshRequest struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	ExpiredAt    int64  `json:"token_expired_at"`
	RefreshToken string `json:"refresh_token"`
	RTExpiredAt  int64  `json:"refresh_token_expired_at"`
}

// Organizations

type NewOrganizationRequest struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	WalletMnemonic string `json:"wallet_mnemonic,omitempty"`
}

type ListOrganizationsRequest struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  uint8  `json:"limit,omitempty"` // Default: 50, Max: 50
}

// Transactions

type NewTransactionRequest struct {
	Description   string  `json:"description,omitempty"`
	Amount        float64 `json:"amount,omitempty"`
	ToAddr        string  `json:"to,omitempty"`
	MaxFeeAllowed float64 `json:"max_fee_allowed,omitempty"`
	Deadline      int64   `json:"deadline,omitempty"`
}

type ListTransactionsRequest struct {
	IDs       []string `json:"ids,omitempty"`
	CreatedBy string   `json:"created_by,omitempty"`
	To        string   `json:"to,omitempty"`

	Cancelled bool `json:"cancelled,omitempty"`
	Confirmed bool `json:"confirmed,omitempty"`
	Commited  bool `json:"commited,omitempty"`
	Expired   bool `json:"expired,omitempty"`
	Pending   bool `json:"pending,omitempty"`

	Cursor string `json:"cursor,omitempty"`
	Limit  uint8  `json:"limit,omitempty"` // Default: 50, Max: 50
}

type UpdateTransactionStatusRequest struct {
	Cancel  bool `json:"cancel,omitempty"`
	Confirm bool `json:"confirm,omitempty"`
}

// Participants

type ListParticipantsRequest struct {
	IDs    []string `json:"ids,omitempty"`
	Cursor string   `json:"cursor,omitempty"` // not implemented
	Limit  uint8    `json:"limit,omitempty"`  // not implemented
}

type AddEmployeeRequest struct {
	Name          string `json:"name"`
	Position      string `json:"position"`
	WalletAddress string `json:"wallet_address"`
}

type NewMultisigRequest struct {
	Title  string `json:"title"`
	Owners []struct {
		PublicKey string `json:"public_key"`
	} `json:"owners"`
	Confirmations int `json:"confirmations"`
}

type NewInviteLinkRequest struct {
	ExpirationDate int `json:"expiration_date"`
}
