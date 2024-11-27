package main

import (
	"os"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/eventpolling"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/httphandler"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/httphandler/controller"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store"
	"github.com/L4B0MB4/PRYVT/identification/pkg/query/store/repository"
	"github.com/PRYVT/utils/pkg/auth"
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
	tokenManager, err := auth.NewTokenManager()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of token manager")
		return
	}
	eventRepo := repository.NewEventRepository(conn)
	userRepo := repository.NewUserRepository(conn)
	uc := controller.NewUserController(userRepo, tokenManager)
	aut := auth.NewAuthMiddleware(tokenManager)
	h := httphandler.NewHttpHandler(uc, aut)

	eventPolling := eventpolling.NewEventPolling(c, eventRepo, userRepo)
	go eventPolling.PollEvents()

	h.Start()
}
