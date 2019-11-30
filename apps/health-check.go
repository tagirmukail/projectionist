package apps

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"projectionist/config"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"projectionist/utils"
)

type HealthCheck struct {
	httpClient http.Client
	cfg        *config.Config
	dbProvider provider.IDBProvider
}

func NewHealthCkeck(cfg *config.Config, db *sql.DB) *HealthCheck {
	return &HealthCheck{
		cfg: cfg,
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

		hc.run(service)

	}
	log.Printf("health-check: initialization services finished")

	return nil
}

func (hc *HealthCheck) run(service *models.Service) {
	ticker := time.NewTicker(time.Duration(service.Frequency) * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				// todo implement graceful stop
				return
			case <-ticker.C:
				err := hc.Health(service)
				if err != nil {
					if service.Status == models.Alive {
						service.Status = models.Dead
					}
					// todo send message
					log.Printf("service with id %d and name %s Health error: %v", service.ID, service.Name, err)
					err := hc.dbProvider.Update(service, service.ID)
					if err != nil {
						// todo send message
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
						// todo send message
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
