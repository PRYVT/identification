package events

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	m "github.com/L4B0MB4/PRYVT/identification/pkg/models/command"
)

type UserCreatedEvent struct {
	Name         string
	DisplayName  string
	PasswordHash string
	Email        string
	CreationDate time.Time
}

func NewUserCreateEvent(uc m.UserCreate, hashedPw string) *models.ChangeTrackedEvent {

	b := UnsafeSerializeAny(UserCreatedEvent{
		Name:         uc.Name,
		DisplayName:  uc.Name,
		PasswordHash: hashedPw,
		Email:        uc.Email,
		CreationDate: time.Now(),
	})
	return &models.ChangeTrackedEvent{
		Event: models.Event{
			Name: "UserCreatedEvent",
			Data: b,
		},
	}
}
