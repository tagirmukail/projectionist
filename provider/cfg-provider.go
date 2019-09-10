package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"projectionist/consts"
	"projectionist/models"
	"sync"
)

type CfgProvider struct {
	*sync.RWMutex
	configs []models.Model
	maxID   int
}

func NewCfgProvider() (*CfgProvider, error) {
	var cfgProvider = &CfgProvider{
		RWMutex: &sync.RWMutex{},
		configs: []models.Model{},
	}

	var err = filepath.Walk(consts.PathSaveCfgs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != "json" {
			log.Printf("NewCfgProvider() error: file %v extension not json", path)
			return nil
		}

		bfile, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		var cfg = models.Configuration{}

		err = json.Unmarshal(bfile, &cfg)
		if err != nil {
			return err
		}

		if cfg.ID > cfgProvider.maxID {
			cfgProvider.maxID = cfg.ID
		}

		cfgProvider.configs = append(cfgProvider.configs, &cfg)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return cfgProvider, nil
}

func (c *CfgProvider) Save(m models.Model) error {
	m.SetID(c.maxID + 1)

	c.Lock()
	c.configs = append(c.configs, m)
	c.Unlock()

	return m.Save()
}

func (c *CfgProvider) GetByName(m models.Model, name string) error {
	c.RLock()
	defer c.RUnlock()
	for _, model := range c.configs {
		if name == model.GetName() {
			m = model
			return nil
		}
	}

	return fmt.Errorf("configuration with name %v not exist", name)
}

func (c *CfgProvider) GetByID(m models.Model, id int64) (models.Model, error) {
	c.RLock()
	defer c.RUnlock()
	for _, model := range c.configs {
		if m.GetID() == model.GetID() {
			return model, nil
		}
	}

	return nil, fmt.Errorf("configuration with id %v not exist", id)
}

func (c *CfgProvider) IsExistByName(m models.Model) (error, bool) {
	c.RLock()
	defer c.RUnlock()
	for _, model := range c.configs {
		if m.GetName() == model.GetName() {
			return nil, true
		}
	}

	return nil, false
}

func (c *CfgProvider) Count(m models.Model) (int, error) {
	c.RLock()
	var count = len(c.configs)
	c.RUnlock()

	return count, nil
}

func (c *CfgProvider) Pagination(m models.Model, start, stop int) ([]models.Model, error) {
	var result []models.Model
	c.RLock()
	var lastCfgNumber = len(c.configs) - 1
	c.RUnlock()
	if stop > lastCfgNumber {
		stop = lastCfgNumber
	}

	if start < 0 {
		start = 0
	}

	c.RLock()
	for i := start; i <= stop; i++ {
		result = append(result, c.configs[i])
	}
	c.RUnlock()

	return result, nil
}

func (c *CfgProvider) Update(m models.Model, id int) error {
	c.Lock()
	for i, model := range c.configs {
		if model.GetID() == id {
			c.configs[i] = m
			break
		}
	}
	c.Unlock()

	return m.Update(id)
}

func (c *CfgProvider) Delete(m models.Model, id int) error {
	c.Lock()
	for i, model := range c.configs {
		if model.GetID() == id {
			c.configs[i].SetDeleted()
			break
		}
	}
	c.Unlock()

	return m.Delete(id)
}