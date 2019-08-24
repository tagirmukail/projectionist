package provider

import (
	"fmt"
	"projectionist/models"
	"sort"
)

type MockUsersDBProvider struct {
	users map[int]models.Model
}

func NewMockUsersDBProvider(users map[int]models.Model) *MockUsersDBProvider {
	return &MockUsersDBProvider{
		users: users,
	}
}

func (m *MockUsersDBProvider) Save(model models.Model) error {
	user, ok := model.(*models.User)
	if !ok {
		return fmt.Errorf("model %v not user", model)
	}

	m.users[user.ID] = model

	return nil
}

func (m *MockUsersDBProvider) GetByID(model models.Model, id int64) error {
	model = m.users[int(id)]
	return nil
}

func (m *MockUsersDBProvider) GetByName(model models.Model, name string) error {
	for _, iUser := range m.users {
		var user, ok = iUser.(*models.User)
		if !ok {
			return fmt.Errorf("model %v not user", iUser)
		}

		if user.Username == name {
			model = user
			return nil
		}
	}

	return fmt.Errorf("model with name %s not exist", name)
}

func (m *MockUsersDBProvider) IsExist(model models.Model) (error, bool) {
	user := model.(*models.User)
	_, ok := m.users[user.ID]
	if !ok {
		return fmt.Errorf("model id: %d not exist", user.ID), false
	}
	return nil, true
}

func (m *MockUsersDBProvider) Count(model models.Model) (int, error) {
	return len(m.users), nil
}

func (m *MockUsersDBProvider) Pagination(model models.Model, start, stop int) ([]models.Model, error) {
	var (
		users  []*models.User
		result []models.Model
	)

	for _, model := range m.users {
		user, ok := model.(*models.User)
		if !ok {
			return nil, fmt.Errorf("in users exist not user model")
		}

		users = append(users, user)
	}

	sort.Slice(&users, func(i, j int) bool {
		return users[i].ID > users[j].ID
	})

	for i, user := range users {
		if i+1 >= start && i+1 < stop {
			result = append(result, user)
		}
	}

	return result, nil
}

func (m *MockUsersDBProvider) Update(model models.Model, id int) error {
	m.users[id] = model

	return nil
}

func (m *MockUsersDBProvider) Delete(model models.Model, id int) error {
	delete(m.users, id)
	return nil
}
