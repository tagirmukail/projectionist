package provider

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"projectionist/models"
	"sort"
)

type MockUsersDBProvider struct {
	users map[int]models.Model
}

func NewMockUsersDBProvider(users map[int]models.Model) *MockUsersDBProvider {
	var resultUsers = make(map[int]models.Model)
	for id, iUser := range users {
		user, ok := iUser.(*models.User)
		if !ok {
			continue
		}
		passwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			log.Println(err)
			return nil
		}
		user.Password = string(passwd)
		resultUsers[id] = user
	}
	return &MockUsersDBProvider{
		users: resultUsers,
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

	passwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		return err
	}
	user.Password = string(passwd)

	m.users[user.ID] = user

	return nil
}

func (m *MockUsersDBProvider) GetByID(model models.Model, id int64) (models.Model, error) {
	existModel := m.users[int(id)]

	return existModel, nil
}

func (m *MockUsersDBProvider) GetByName(model models.Model, name string) (models.Model, error) {
	for _, iUser := range m.users {
		if iUser.GetName() == name {
			model = iUser
			return model, nil
		}
	}

	return nil, sql.ErrNoRows
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
