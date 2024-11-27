package main

import (
	"os"

	"github.com/L4B0MB4/PRYVT/identification/pkg/command/httphandler"
	"github.com/L4B0MB4/PRYVT/identification/pkg/command/httphandler/controller"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	uc := controller.NewUserController()
	h := httphandler.NewHttpHandler(uc)

	h.Start()
}
