package eventpolling

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/L4B0MB4/PRYVT/identification/pkg/aggregates"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store/repository"
	"github.com/google/uuid"
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

func (ep *EventPolling) PollEvents() {

	for {
		time.Sleep(50 * time.Millisecond)
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
			if event.AggregateType == "user" {
				ua, err := aggregates.NewUserAggregate(uuid.MustParse(event.AggregateId))
				if err != nil {
					log.Err(err).Msg("Error while creating user aggregate")
					break
				}
				uI := aggregates.GetUserModelFromAggregate(ua)
				err = ep.userRepo.AddOrReplaceUser(uI)
				if err != nil {
					log.Err(err).Msg("Error while adding or replacing user")
					break
				}
			}
		}
		if len(events) == 0 {
			continue
		}
		//will this break the db consistency if there are going to be multiple instances of this service?
		// probably but if we dont a volume (that both instances use as a db file) this should be fine
		err = ep.eventRepo.ReplaceEvent(events[len(events)-1].Id)
		if err != nil {
			log.Err(err).Msg("Error while replacing event")
			break
		}
	}

}
