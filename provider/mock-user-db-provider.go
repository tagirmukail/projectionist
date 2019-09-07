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
	var maxID int

	for id := range m.users {
		if id > maxID {
			maxID = id
		}
	}

	user, ok := model.(*models.User)
	if !ok {
		return fmt.Errorf("model %v not user", model)
	}

	user.ID = maxID + 1

	m.users[user.ID] = user

	return nil
}

func (m *MockUsersDBProvider) GetByID(model models.Model, id int64) (models.Model, error) {
	existModel := m.users[int(id)]

	return existModel, nil
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

func (m *MockUsersDBProvider) IsExistByName(model models.Model) (error, bool) {
	userModel := model.(*models.User)

	for _, user := range m.users {
		existUser := user.(*models.User)
		if userModel.Username == existUser.Username {
			return nil, true
		}
	}

	return nil, false
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

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	for i, user := range users {
		if i >= start && i < stop {
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
	for userID, _ := range m.users {
		if userID == id {
			delete(m.users, id)
			return nil
		}
	}

	return fmt.Errorf("not exist")
}
