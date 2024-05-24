package domain

type Participant struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Position string `json:"position,omitempty"`

	// if is_user == true, Credentials will be filled with data
	Credentials *UserParticipantCredentials `json:"credentials,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
	DeletedAt int64 `json:"deleted_at,omitempty"`

	IsUser   bool `json:"is_user"`
	IsAdmin  bool `json:"is_admin"`
	IsOwner  bool `json:"is_owner"`
	IsActive bool `json:"is_active"`
}

type UserParticipantCredentials struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Telegram string `json:"telegram,omitempty"`
}
