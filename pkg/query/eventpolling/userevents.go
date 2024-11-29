package eventpolling

import (
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/PRYVT/identification/pkg/aggregates"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (ep *EventPolling) ProcessUserEvent(event models.Event) error {
	if event.AggregateType == "user" {
		ua, err := aggregates.NewUserAggregate(uuid.MustParse(event.AggregateId))
		if err != nil {
			return err
		}
		uI := aggregates.GetUserModelFromAggregate(ua)
		err = ep.userRepo.AddOrReplaceUser(uI)
		if err != nil {
			log.Err(err).Msg("Error while processing user event")
			return err
		}
	}
	return nil
}
