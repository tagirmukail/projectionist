package apps

import (
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"

	"projectionist/config"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

type HealthCheck struct {
	notifier *Notifier
	syncChan chan string // chan with service name
	sync.Mutex
	dones      map[int]chan bool
	httpClient http.Client
	cfg        *config.Config
	dbProvider provider.IDBProvider
}

func NewHealthCkeck(cfg *config.Config, db *sql.DB, syncChan chan string) *HealthCheck {
	return &HealthCheck{
		notifier: NewNotifier(cfg),
		syncChan: syncChan,
		Mutex:    sync.Mutex{},
		dones:    make(map[int]chan bool),
		cfg:      cfg,
		httpClient: http.Client{Transport: &http.Transport{
			MaxIdleConns:       cfg.HealthCheck.ConnCount,
			IdleConnTimeout:    time.Duration(cfg.HealthCheck.ConnTimeout) * time.Second,
			DisableCompression: true,
		}},
		dbProvider: provider.NewDBProvider(db),
	}
}

func (hc *HealthCheck) Run() error {
	log.Printf("health-check: initialization services started")

	services, err := hc.dbProvider.Pagination(&models.Service{}, 0, -1)
	if err != nil {
		return err
	}

	for _, iService := range services {
		service, ok := iService.(*models.Service)
		if !ok {
			continue
		}

		if service.IsDeleted() {
			continue
		}

		hc.Lock()
		hc.dones[service.ID] = make(chan bool)

		hc.run(service, hc.dones[service.ID])
		hc.Unlock()

	}

	go hc.watcher()

	log.Printf("health-check: initialization services finished")

	return nil
}

func (hc *HealthCheck) watcher() {
	for {
		select {
		case serviceName := <-hc.syncChan:
			iService, err := hc.dbProvider.GetByName(&models.Service{}, serviceName)
			if err != nil {
				log.Printf("health-check: dbProvider.GetByName() error: %v", err)
				continue
			}

			service, ok := iService.(*models.Service)
			if !ok {
				log.Printf("health-check: %+v is not service", iService)
				continue
			}

			hc.Lock()
			_, ok = hc.dones[service.ID]
			if !ok {
				hc.dones[service.ID] = make(chan bool)
			}

			if service.IsDeleted() {
				hc.dones[service.ID] <- true
				hc.Unlock()
				continue
			}

			hc.run(service, hc.dones[service.ID])
			hc.Unlock()
		default:
			continue
		}
	}
}

func (hc *HealthCheck) run(service *models.Service, done chan bool) {
	ticker := time.NewTicker(time.Duration(service.Frequency) * time.Second)
	go func() {
		for {
			select {
			case <-done:
				log.Printf("health-check: for service %s health check stoped", service.Name)
				return
			case <-ticker.C:
				err := hc.Health(service)
				if err != nil {
					if service.Status == models.Alive {
						service.Status = models.Dead
					}
					//hc.notifier.Send() TODO send message logic
					log.Printf("service with id %d and name %s Health error: %v", service.ID, service.Name, err)
					err := hc.dbProvider.Update(service, service.ID)
					if err != nil {
						//hc.notifier.Send() TODO send message logic
						log.Printf(
							"for service with id %d and name %s status %v not updated",
							service.ID,
							service.Name,
							models.Dead)
					}
					continue
				}

				if service.Status == models.Dead {
					service.Status = models.Alive
					err = hc.dbProvider.Update(service, service.ID)
					if err != nil {
						//hc.notifier.Send() TODO send message logic
						log.Printf(
							"for service with id %d and name %s status %v not updated",
							service.ID,
							service.Name,
							models.Alive,
						)
					}
				}
			}
		}
	}()
}

func (hc *HealthCheck) Health(service *models.Service) error {
	req, err := http.NewRequest(http.MethodGet, service.Link, nil)
	if err != nil {
		return err
	}

	req.Header.Set(consts.AuthorizationHeader, service.Token)

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return err
	}

	return utils.CheckHealthStatusCode(resp.StatusCode, service.Name)
}
