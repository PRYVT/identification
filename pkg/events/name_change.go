package events

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
)

type DisplayNameChangedEvent struct {
	DisplayName string
	ChangeDate  time.Time
}

func NewNameChangedEvent(name string) *models.ChangeTrackedEvent {
	b := UnsafeSerializeAny(DisplayNameChangedEvent{DisplayName: name, ChangeDate: time.Now()})
	return &models.ChangeTrackedEvent{
		Event: models.Event{

			Name: "NameChangeEvent",
			Data: b,
		},
	}
}
