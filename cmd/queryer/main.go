package main

import (
	"os"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	tcpClient "github.com/L4B0MB4/EVTSRC/pkg/tcp/client"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/eventhandling"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/httphandler"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/httphandler/controller"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store/repository"
	"github.com/PRYVT/utils/pkg/auth"
	"github.com/PRYVT/utils/pkg/eventpolling"
	utilsRepo "github.com/PRYVT/utils/pkg/store/repository"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db := store.DatabaseConnection{}
	db.SetUp()
	conn, err := db.GetDbConnection()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessfull initalization of db")
		return
	}
	log.Debug().Msg("Db Connection was successful")

	c, err := client.NewEventSourcingHttpClient(client.RetrieveEventSourcingClientUrl())
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of client")
		return
	}
	eventRepo := utilsRepo.NewEventRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	uc := controller.NewUserController(userRepo)
	aut := auth.NewAuthMiddleware()
	h := httphandler.NewHttpHandler(uc, aut)

	userEventHandler := eventhandling.NewUserEventHandler(userRepo)

	eventPolling := eventpolling.NewEventPolling(c, eventRepo, userEventHandler)

	tcpC, err := tcpClient.NewTcpEventClient()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of tcp client")
		return
	}
	channel := make(chan string, 1)
	go tcpC.ListenForEvents(channel)

	eventPolling.PollEventsUntilEmpty()
	go func() {
		for {
			select {
			case event := <-channel:
				log.Info().Msgf("Received event: %s", event)
				eventPolling.PollEventsUntilEmpty()
			}
		}
	}()
	h.Start()
}
