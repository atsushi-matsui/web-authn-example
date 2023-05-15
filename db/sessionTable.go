package db

import (
	"fmt"
	"sync"

	"github.com/atsushi-matsui/web-authn-example/domain"
	"github.com/go-webauthn/webauthn/webauthn"
)

type SessionTable struct {
	webAuthnSessions map[uint64]*webauthn.SessionData
	mu sync.Mutex
}

var sTable *SessionTable 

func NewSessionTable() *SessionTable {
	if sTable == nil {
		sTable = &SessionTable{
			webAuthnSessions: make(map[uint64]*webauthn.SessionData),
		}
	}

	return sTable
}

func (table *SessionTable) PullSession(userId uint64) (*webauthn.SessionData, error) {
	table.mu.Lock()
	defer table.mu.Unlock()

	session, ok := table.webAuthnSessions[userId]
	if !ok {
		return &webauthn.SessionData{}, fmt.Errorf("error getting session userId '%d': does not exist", userId)
	}

	// キャッシュとして扱いたいので削除
	table.webAuthnSessions[userId] = nil

	return session, nil
}

func (table *SessionTable) PutSession(user *domain.User, session *webauthn.SessionData) {
	table.mu.Lock()
	defer table.mu.Unlock()

	table.webAuthnSessions[user.GetId()] = session
}
