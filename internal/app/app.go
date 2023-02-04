package app

import (
	"log"

	"github.com/kamkalis/object-storage/internal/config"
	"github.com/kamkalis/object-storage/internal/server"
)

type app struct {
	config *config.Config
	server *server.Server
}

func (a *app) initConfig() {
	c, err := config.Load()
	if err != nil {
		log.Fatalf("cannot initialize config for app: %v\n", err)
	}
	a.config = c
}

func (a *app) initApp() {
	a.initConfig()
	a.initHTTPServer()
}

func (a *app) initHTTPServer() {
	s, err := server.New(a.config)
	if err != nil {
		log.Fatalf("cannot init server: %v\n", err)
	}
	a.server = s
}

func (a *app) start() {
	a.server.Start()
}

func Run() {
	a := app{}
	a.initApp()
	a.start()
}
