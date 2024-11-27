package query

import (
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID           uuid.UUID `json:"id"`
	DisplayName  string    `json:"display_name"`
	Name         string    `json:"name,omitempty"`
	PasswordHash string    `json:"password_hash,omitempty"`
	Email        string    `json:"email,omitempty"`
	ChangeDate   time.Time `json:"change_date,omitempty"`
}
