package provider

import (
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

func (c *MockCfgProvider) GetByName(m models.Model, name string) error {
	c.RLock()
	for _, model := range c.configs {
		if name == model.GetName() {
			m = model
			c.RUnlock()
			return nil
		}
	}
	c.RUnlock()

	return nil
}

func (c *MockCfgProvider) GetByID(m models.Model, id int64) (models.Model, error) {
	c.RLock()
	for _, model := range c.configs {
		if m.GetID() == model.GetID() {
			c.RUnlock()
			return model, nil
		}
	}
	c.RUnlock()

	return m, nil
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

func (c *MockCfgProvider) Update(m models.Model, id int) error {
	c.Lock()
	for i, model := range c.configs {
		if model.GetID() == id {
			c.configs[i] = m
			break
		}
	}
	c.Unlock()

	return nil
}

func (c *MockCfgProvider) Delete(m models.Model, id int) error {
	c.Lock()
	for i, model := range c.configs {
		if model.GetID() == id {
			c.configs[i].SetDeleted()
			break
		}
	}
	c.Unlock()

	return nil
}
