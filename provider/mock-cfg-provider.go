package provider

import (
	"projectionist/consts"
	"projectionist/models"
	"sync"
)

type MockCfgProvider struct {
	*sync.RWMutex
	configs []models.Model
	maxID   int
}

func NewMockCfgProvider(configs []models.Model) *MockCfgProvider {
	return &MockCfgProvider{
		RWMutex: &sync.RWMutex{},
		configs: configs,
	}
}

func (c *MockCfgProvider) Save(m models.Model) error {
	m.SetID(c.maxID + 1)

	c.Lock()
	c.configs = append(c.configs, m)
	c.Unlock()

	return nil
}

func (c *MockCfgProvider) GetByName(m models.Model, name string) (models.Model, error) {
	c.RLock()
	for _, model := range c.configs {
		if name == model.GetName() {
			c.RUnlock()
			return model, nil
		}
	}
	c.RUnlock()

	return nil, consts.ErrNotFound
}

func (c *MockCfgProvider) GetByID(m models.Model, id int64) (models.Model, error) {
	c.RLock()
	for _, model := range c.configs {
		if id == int64(model.GetID()) {
			c.RUnlock()
			return model, nil
		}
	}
	c.RUnlock()

	return m, consts.ErrNotFound
}

func (c *MockCfgProvider) IsExistByName(m models.Model) (error, bool) {
	c.RLock()
	for _, model := range c.configs {
		if m.GetName() == model.GetName() {
			c.RUnlock()
			return nil, true
		}
	}
	c.RUnlock()

	return nil, false
}

func (c *MockCfgProvider) Count(m models.Model) (int, error) {
	c.RLock()
	var count = len(c.configs)
	c.RUnlock()

	return count, nil
}

func (c *MockCfgProvider) Pagination(m models.Model, start, stop int) ([]models.Model, error) {
	var result []models.Model
	c.RLock()
	var lastCfgNumber = len(c.configs)
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

	return result, nil
}

func (c *MockCfgProvider) Update(m models.Model, id int) error {
	c.Lock()
	for i, model := range c.configs {
		if model.GetID() == id {
			c.configs[i] = m
			return nil
		}
	}
	c.Unlock()

	return consts.ErrNotFound
}

func (c *MockCfgProvider) Delete(m models.Model, id int) error {
	c.Lock()
	for i, model := range c.configs {
		if model.GetID() == id {
			c.configs[i].SetDeleted()
			c.configs = append(c.configs[:i], c.configs[i+1:]...)
			return nil
		}
	}
	c.Unlock()

	return consts.ErrNotFound
}
