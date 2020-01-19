package healtchecker

import (
	"database/sql"
	"fmt"
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc/grpclog"
	"net/http"
	apps "projectionist/apps/notifier"
	"projectionist/utils"
	"sync"
	"time"

	"projectionist/config"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
)

const EveryDurationPtrn = "@every %s"

type HealthCheck struct {
	notifier *apps.Notifier
	syncChan chan string // chan with service name
	errChan  chan error
	sync.Mutex
	cronEntries map[int]cron.EntryID // map[service id]cron entry id
	httpClient  http.Client
	cfg         *config.Config
	dbProvider  provider.IDBProvider
	crontab     *cron.Cron
}

func NewHealthCkeck(cfg *config.Config, db *sql.DB, syncChan chan string) *HealthCheck {
	return &HealthCheck{
		notifier:    apps.NewNotifier(cfg),
		syncChan:    syncChan,
		errChan:     make(chan error),
		Mutex:       sync.Mutex{},
		cfg:         cfg,
		crontab:     cron.New(),
		cronEntries: make(map[int]cron.EntryID),
		httpClient: http.Client{Transport: &http.Transport{
			MaxIdleConns:       cfg.HealthCheck.ConnCount,
			IdleConnTimeout:    time.Duration(cfg.HealthCheck.ConnTimeout) * time.Second,
			DisableCompression: true,
		}},
		dbProvider: provider.NewDBProvider(db),
	}
}

// Run - run health check on services
func (hc *HealthCheck) Run() error {
	grpclog.Info("health-check: initialization services started")

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

		err = hc.plan(service)
		if err != nil {
			return err
		}

	}

	hc.crontab.Start()

	hc.watcher()

	err = hc.watchErr()
	if err != nil {
		grpclog.Errorf("healtchecker.watchErr error: %v", err)
		return err
	}

	return nil
}

// Stop - stop healt checker cron tasks
func (hc *HealthCheck) Stop() {
	hc.crontab.Stop()
}

// plan - add cron entry by service and change service health status
func (hc *HealthCheck) plan(service *models.Service) error {
	entryID, err := hc.crontab.AddFunc(fmt.Sprintf(EveryDurationPtrn, time.Duration(service.Frequency)*time.Second), func() {
		err := hc.Health(service)
		if err != nil {
			if service.Status == models.Alive {
				service.Status = models.Dead
			}
			//hc.notifier.Send() TODO send message logic
			grpclog.Errorf("service with id %d and name %s Health error: %v", service.ID, service.Name, err)
			err := hc.dbProvider.Update(service, service.ID)
			if err != nil {
				//hc.notifier.Send() TODO send message logic
				grpclog.Errorf(
					"for service with id %d and name %s status %v not updated",
					service.ID,
					service.Name,
					models.Dead)
			}
			return
		}

		if service.Status == models.Dead {
			service.Status = models.Alive
			err = hc.dbProvider.Update(service, service.ID)
			if err != nil {
				//hc.notifier.Send() TODO send message logic
				grpclog.Errorf(
					"for service with id %d and name %s status %v not updated",
					service.ID,
					service.Name,
					models.Alive,
				)
			}
		}
	})
	if err != nil {
		return err
	}

	hc.Lock()
	hc.cronEntries[service.ID] = entryID
	hc.Unlock()

	return nil
}

// watcher - watch if serivce added, updated, deleted run this logic
func (hc *HealthCheck) watcher() {
	for {
		select {
		case serviceName := <-hc.syncChan:
			iService, err := hc.dbProvider.GetByName(&models.Service{}, serviceName)
			if err != nil {
				grpclog.Errorf("health-check: dbProvider.GetByName() error: %v", err)
				hc.errChan <- err
				break
			}

			service, ok := iService.(*models.Service)
			if !ok {
				grpclog.Errorf("health-check: %+v is not service", iService)
				break
			}

			hc.Lock()
			entryID, ok := hc.cronEntries[service.ID]
			hc.Unlock()
			if !ok {
				err := hc.plan(service)
				hc.errChan <- err
				break
			}

			if service.IsDeleted() {
				hc.crontab.Remove(entryID)
				break
			}
		default:
			break
		}
	}
}

// watchErr - watcher on errors
func (hc *HealthCheck) watchErr() error {
	for {
		select {
		case err := <-hc.errChan:
			if err != nil {
				return err
			}
		default:
		}
	}
}

// Health - send request by service healt link and check service status
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
