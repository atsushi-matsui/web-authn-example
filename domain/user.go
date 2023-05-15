package domain

import (
	"encoding/binary"
	"math/rand"

	"github.com/go-webauthn/webauthn/webauthn"
)


type User struct {
	id uint64
	name string
	displayName string
	credentials []webauthn.Credential
}

func NewUser(name string, displayName string) *User {
	user := &User{}
	user.id = rand.Uint64()
	user.name = name
	user.displayName = displayName

	return user
}

func (user User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(user.id))
	return buf
}

func (user User) WebAuthnName() string {
	return user.name
}

func (user User) WebAuthnDisplayName() string {
	return user.displayName
}

func (user User) WebAuthnIcon() string {
	return ""
}

func (user User) WebAuthnCredentials() []webauthn.Credential {
	return user.credentials
}

func (user *User) AddCredential(cred webauthn.Credential) {
	user.credentials = append(user.credentials, cred)
}

func (user User) GetId() uint64 {
	return user.id
} 

func (user User) GetName() string {
	return user.name
} 
