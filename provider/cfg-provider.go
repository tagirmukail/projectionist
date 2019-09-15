package provider

import (
	"encoding/json"
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

		cfgProvider.Lock()
		cfgProvider.configs = append(cfgProvider.configs, &cfg)
		cfgProvider.Unlock()

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
	for _, model := range c.configs {
		if name == model.GetName() {
			m = model
			c.RUnlock()
			return nil
		}
	}
	c.RUnlock()

	return m.GetByName(name)
}

func (c *CfgProvider) GetByID(m models.Model, id int64) (models.Model, error) {
	c.RLock()
	for _, model := range c.configs {
		if id == int64(model.GetID()) {
			c.RUnlock()
			return model, nil
		}
	}
	c.RUnlock()

	return m, m.GetByID(id)
}

func (c *CfgProvider) IsExistByName(m models.Model) (error, bool) {
	c.RLock()
	for _, model := range c.configs {
		if m.GetName() == model.GetName() {
			c.RUnlock()
			return nil, true
		}
	}
	c.RUnlock()

	return m.IsExistByName()
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

	var maxResultCount = stop - start
	var resultCount int

	c.RLock()
	for i := 0; i < len(c.configs); i++ {
		if resultCount >= maxResultCount {
			break
		}

		if c.configs[i].IsDeleted() {
			continue
		}

		if i < start {
			continue
		}

		result = append(result, c.configs[i])
		resultCount++
	}
	c.RUnlock()

	if len(result) == 0 {
		return m.Pagination(start, stop)
	}

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
			c.configs = append(c.configs[:i], c.configs[i+1:]...)
			break
		}
	}
	c.Unlock()

	return m.Delete(id)
}
