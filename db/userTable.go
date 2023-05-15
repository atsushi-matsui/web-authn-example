package db

import (
	"fmt"
	"sync"

	"github.com/atsushi-matsui/web-authn-example/domain"
)

type UserTable struct {
	users map[string]*domain.User
	mu sync.Mutex
}

var uTable *UserTable

func NewUserTable() *UserTable {
	if uTable == nil {
		uTable = &UserTable{
			users: make(map[string]*domain.User),
		}
	}

	return uTable
}

func (table *UserTable) GetUser(name string) (*domain.User, error) {
	table.mu.Lock()
	defer table.mu.Unlock()

	user, ok := table.users[name]
	if !ok {
		return &domain.User{}, fmt.Errorf("error getting user '%s': does not exist", name)
	}

	return user, nil
}

func (table *UserTable) PutUser(user *domain.User) {
	table.mu.Lock()
	defer table.mu.Unlock()

	table.users[user.GetName()] = user
}
