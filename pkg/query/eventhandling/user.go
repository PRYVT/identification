package eventhandling

import (
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/PRYVT/identification/pkg/aggregates"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserEventHandler struct {
	userRepo *repository.UserRepository
}

func NewUserEventHandler(userRepo *repository.UserRepository) *UserEventHandler {
	return &UserEventHandler{
		userRepo: userRepo,
	}
}

func (eh *UserEventHandler) HandleEvent(event models.Event) error {
	if event.AggregateType == "user" {
		ua, err := aggregates.NewUserAggregate(uuid.MustParse(event.AggregateId))
		if err != nil {
			return err
		}
		uI := aggregates.GetUserModelFromAggregate(ua)
		err = eh.userRepo.AddOrReplaceUser(uI)
		if err != nil {
			log.Err(err).Msg("Error while processing user event")
			return err
		}
	}
	return nil
}
