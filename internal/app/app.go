package app

import (
	"context"
	"log"
	"time"

	consistenthash "github.com/kamkalis/object-storage/internal/domain/consistent-hash"
	"github.com/kamkalis/object-storage/internal/domain/docker"
	"github.com/kamkalis/object-storage/internal/domain/manager"

	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/service"

	"github.com/kamkalis/object-storage/internal/config"
	"github.com/kamkalis/object-storage/internal/server"
)

type app struct {
	config         *config.Config
	server         *server.Server
	storageService domain.StorageService
	manager        domain.NodeManager
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
	a.initManager()
	a.initStorageService()
	a.initHTTPServer()
}

func (a *app) initHTTPServer() {
	s, err := server.New(a.config, a.storageService)
	if err != nil {
		log.Fatalf("cannot init server: %v\n", err)
	}
	a.server = s
}

func (a *app) start() {
	ctx := context.Background()
	err := a.manager.RefreshNodes(ctx)
	if err != nil {
		log.Fatalf("failed initial node refresh: %v", err)
	}
	go a.refreshJob(ctx, time.NewTicker(a.config.Discovery.RefreshDuration))
	a.server.Start()
}

func (a *app) initStorageService() {
	a.storageService = service.NewStorage(a.manager)
}

func (a *app) initManager() {
	lb := consistenthash.NewRingLoadBalancer()
	discoverer, err := docker.NewNodeDiscoverer(
		a.config.StorageCluster.NodeIdentifier,
		a.config.StorageCluster.NodeAPIPort,
		a.config.StorageCluster.NetworkIdentifier,
	)
	if err != nil {
		log.Fatalf("cannot init node discoverer: %v", err)
	}
	a.manager = manager.NewStorageManager(lb, discoverer)
}

func (a *app) refreshJob(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			if err := a.manager.RefreshNodes(ctx); err != nil {
				log.Printf("failed to refresh nodes in background: %v\n", err)
			}
		case <-ctx.Done():
			log.Printf("ctx done: %v", ctx.Err())
			return
		}
	}
}

func Run() {
	a := app{}
	a.initApp()
	a.start()
}
