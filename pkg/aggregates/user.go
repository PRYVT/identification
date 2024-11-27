package aggregates

import (
	"fmt"
	"strings"
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/PRYVT/identification/pkg/events"
	m "github.com/L4B0MB4/PRYVT/identification/pkg/models/command"
	"github.com/PRYVT/utils/pkg/hash"
	"github.com/google/uuid"
)

type UserAggregate struct {
	DisplayName   string
	Name          string
	PasswordHash  string
	Email         string
	ChangeDate    time.Time
	Events        []models.ChangeTrackedEvent
	aggregateType string
	AggregateId   uuid.UUID
	client        *client.EventSourcingHttpClient
}

func NewUserAggregate(id uuid.UUID) (*UserAggregate, error) {

	c, err := client.NewEventSourcingHttpClient(client.RetrieveEventSourcingClientUrl())
	if err != nil {
		panic(err)
	}
	iter, err := c.GetEventsOrdered(id.String())
	if err != nil {
		return nil, fmt.Errorf("COULDN'T RETRIEVE EVENTS ")
	}
	ua := &UserAggregate{
		client:        c,
		Events:        []models.ChangeTrackedEvent{},
		aggregateType: "user",
		AggregateId:   id,
		ChangeDate:    time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
	}

	for {
		ev, ok := iter.Next()
		if !ok {
			break
		}
		changeTrackedEv := models.ChangeTrackedEvent{
			Event: *ev,
			IsNew: false,
		}
		ua.addEvent(&changeTrackedEv)
	}
	return ua, nil
}

func (ua *UserAggregate) apply_DisplayNameChangedEvent(e *events.DisplayNameChangedEvent) {
	ua.DisplayName = e.DisplayName
	ua.ChangeDate = e.ChangeDate

}
func (ua *UserAggregate) apply_UserCreatedEvent(e *events.UserCreatedEvent) {
	ua.Name = e.Name
	ua.DisplayName = e.Name
	ua.ChangeDate = e.CreationDate
	ua.PasswordHash = e.PasswordHash
	ua.Email = e.Email
}

func (ua *UserAggregate) addEvent(ev *models.ChangeTrackedEvent) {
	switch ev.Name {
	case "NameChangeEvent":
		e := events.UnsafeDeserializeAny[events.DisplayNameChangedEvent](ev.Data)
		ua.apply_DisplayNameChangedEvent(e)
	case "UserCreatedEvent":
		e := events.UnsafeDeserializeAny[events.UserCreatedEvent](ev.Data)
		ua.apply_UserCreatedEvent(e)
	default:
		panic(fmt.Errorf("NO KNOWN EVENT %v", ev))
	}
	if ev.Version == 0 {
		ev.IsNew = true
	}
	v := len(ua.Events) + 1 //for validation we need to start at 1
	ev.Version = int64(v)
	ev.AggregateType = ua.aggregateType
	ev.AggregateId = ua.AggregateId.String()
	ua.Events = append(ua.Events, *ev)
}

func (ua *UserAggregate) saveChanges() error {
	return ua.client.AddEvents(ua.AggregateId.String(), ua.Events)
}
func (ua *UserAggregate) ChangeDisplayName(name string) error {
	if len(ua.Events) == 0 {
		return fmt.Errorf("user does not yet exist")
	}

	if ua.DisplayName != name && len(name) <= 50 && time.Since(ua.ChangeDate).Seconds() > 10 {
		ua.addEvent(events.NewNameChangedEvent(name))
		err := ua.saveChanges()
		if err != nil {
			return fmt.Errorf("error saving changes")
		}
		return nil
	}
	return fmt.Errorf("validating username failed")
}

func (ua *UserAggregate) CreateUser(userCreate m.UserCreate) error {

	if len(ua.Events) != 0 {
		return fmt.Errorf("user already exists")
	}

	if !strings.Contains(userCreate.Email, "@") {
		return fmt.Errorf("email does not contain @ sign")

	}
	if !(len(userCreate.Name) > 5 && len(userCreate.Name) < 50) {
		return fmt.Errorf("username not between 5 and 50 characters")
	}
	if !(len(userCreate.Password) >= 8 && len(userCreate.Password) < 50) {

		return fmt.Errorf("password not between 8 and 50 characters")

	}
	hashedPw := hash.HashPassword(userCreate.Password)
	ua.addEvent(events.NewUserCreateEvent(userCreate, hashedPw))
	err := ua.saveChanges()
	if err != nil {
		return fmt.Errorf("ERROR ")
	}
	return nil
}
