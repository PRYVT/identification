package eventpolling

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store/repository"
	"github.com/rs/zerolog/log"
)

type EventPolling struct {
	client    *client.EventSourcingHttpClient
	eventRepo *repository.EventRepository
	userRepo  *repository.UserRepository
}

func NewEventPolling(client *client.EventSourcingHttpClient, eventRepo *repository.EventRepository, userRepo *repository.UserRepository) *EventPolling {
	if client == nil || eventRepo == nil || userRepo == nil {
		return nil
	}
	return &EventPolling{client: client, eventRepo: eventRepo, userRepo: userRepo}
}

func (ep *EventPolling) PollEvents(callback func(event models.Event) error) {

	hadMoreThenZeroEvents := true
	for {
		if hadMoreThenZeroEvents {
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(500 * time.Millisecond)
		}
		eId, err := ep.eventRepo.GetLastEvent()
		if err != nil {
			log.Err(err).Msg("Error while getting last events")
			continue
		}
		events, err := ep.client.GetEventsSince(eId, 2)
		if err != nil {
			log.Err(err).Msg("Error while polling events")
			continue
		}

		for _, event := range events {

			err := callback(event)
			if err != nil {
				log.Err(err).Msg("Error while processing event")
				break
			}
		}
		if len(events) == 0 {
			hadMoreThenZeroEvents = true
			continue
		}
		hadMoreThenZeroEvents = false
		//will this break the db consistency if there are going to be multiple instances of this service?
		// probably but if we dont a volume (that both instances use as a db file) this should be fine
		err = ep.eventRepo.ReplaceEvent(events[len(events)-1].Id)
		if err != nil {
			log.Err(err).Msg("Error while replacing event")
			break
		}
	}

}
