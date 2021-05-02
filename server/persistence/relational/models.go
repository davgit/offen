// Copyright 2020 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

package relational

import (
	"time"

	"github.com/offen/offen/server/persistence"
)

// Event is any analytics event that will be stored in the database. It is
// uniquely tied to an Account and a Secret model.
type Event struct {
	EventID   string `gorm:"primary_key;size:26;unique"`
	Sequence  string `gorm:"size:26"`
	AccountID string `gorm:"size:36"`
	// the secret id is nullable for anonymous events
	SecretID *string `gorm:"size:64"`
	Payload  string  `gorm:"type:text"`
	Secret   Secret  `gorm:"foreignkey:SecretID;association_foreignkey:SecretID"`
}

// A Tombstone replaces an event on its deletion
type Tombstone struct {
	EventID   string  `gorm:"primary_key"`
	AccountID string  `gorm:"size:36"`
	SecretID  *string `gorm:"size:64"`
	Sequence  string  `gorm:"size:26"`
}

// Secret associates a hashed user id - which ties a user and account together
// uniquely - with the encrypted user secret the account owner can use
// to decrypt events stored for that user.
type Secret struct {
	SecretID        string `gorm:"primary_key;size:64;unique"`
	EncryptedSecret string `gorm:"type:text"`
}

// Account stores information about an account.
type Account struct {
	AccountID           string `gorm:"primary_key;size:36;unique"`
	Name                string
	PublicKey           string `gorm:"type:text"`
	EncryptedPrivateKey string `gorm:"type:text"`
	UserSalt            string
	Retired             bool
	AccountStyles       string `gorm:"type:text"`
	Created             time.Time
	Events              []Event `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
}

// AccountUser is a person that can log in and access data related to all
// associated accounts.
type AccountUser struct {
	AccountUserID  string `gorm:"primary_key;size:36;unique"`
	HashedEmail    string
	HashedPassword string
	Salt           string
	AdminLevel     int
	Relationships  []AccountUserRelationship `gorm:"foreignkey:AccountUserID;association_foreignkey:AccountUserID"`
}

// AccountUserRelationship contains the encrypted KeyEncryptionKeys needed for
// an AccountUser to access the data of the account it links to.
type AccountUserRelationship struct {
	RelationshipID                    string `gorm:"primary_key;size:36;unique"`
	AccountUserID                     string `gorm:"size:36"`
	AccountID                         string `gorm:"size:36"`
	PasswordEncryptedKeyEncryptionKey string `gorm:"type:text"`
	EmailEncryptedKeyEncryptionKey    string `gorm:"type:text"`
	OneTimeEncryptedKeyEncryptionKey  string `gorm:"type:text"`
}

func (e *Event) export() persistence.Event {
	return persistence.Event{
		EventID:   e.EventID,
		AccountID: e.AccountID,
		SecretID:  e.SecretID,
		Payload:   e.Payload,
		Secret:    e.Secret.export(),
		Sequence:  e.Sequence,
	}
}

func importEvent(e *persistence.Event) Event {
	return Event{
		EventID:   e.EventID,
		AccountID: e.AccountID,
		SecretID:  e.SecretID,
		Payload:   e.Payload,
		Secret:    importSecret(&e.Secret),
		Sequence:  e.Sequence,
	}
}

func (t *Tombstone) export() persistence.Tombstone {
	return persistence.Tombstone{
		EventID:   t.EventID,
		AccountID: t.AccountID,
		SecretID:  t.SecretID,
		Sequence:  t.Sequence,
	}
}

func importTombstone(t *persistence.Tombstone) *Tombstone {
	return &Tombstone{
		EventID:   t.EventID,
		AccountID: t.AccountID,
		Sequence:  t.Sequence,
		SecretID:  t.SecretID,
	}
}

func (s *Secret) export() persistence.Secret {
	return persistence.Secret{
		SecretID:        s.SecretID,
		EncryptedSecret: s.EncryptedSecret,
	}
}

func importSecret(s *persistence.Secret) Secret {
	return Secret{
		SecretID:        s.SecretID,
		EncryptedSecret: s.EncryptedSecret,
	}
}

func (a *AccountUser) export() persistence.AccountUser {
	var relationships []persistence.AccountUserRelationship
	for _, r := range a.Relationships {
		relationships = append(relationships, r.export())
	}
	return persistence.AccountUser{
		AccountUserID:  a.AccountUserID,
		HashedEmail:    a.HashedEmail,
		HashedPassword: a.HashedPassword,
		Salt:           a.Salt,
		AdminLevel:     persistence.AccountUserAdminLevel(a.AdminLevel),
		Relationships:  relationships,
	}
}

func importAccountUser(a *persistence.AccountUser) AccountUser {
	var relationships []AccountUserRelationship
	for _, r := range a.Relationships {
		relationships = append(relationships, importAccountUserRelationship(&r))
	}
	return AccountUser{
		AccountUserID:  a.AccountUserID,
		HashedEmail:    a.HashedEmail,
		HashedPassword: a.HashedPassword,
		Salt:           a.Salt,
		AdminLevel:     int(a.AdminLevel),
		Relationships:  relationships,
	}
}

func (a *AccountUserRelationship) export() persistence.AccountUserRelationship {
	return persistence.AccountUserRelationship{
		RelationshipID:                    a.RelationshipID,
		AccountUserID:                     a.AccountUserID,
		AccountID:                         a.AccountID,
		PasswordEncryptedKeyEncryptionKey: a.PasswordEncryptedKeyEncryptionKey,
		EmailEncryptedKeyEncryptionKey:    a.EmailEncryptedKeyEncryptionKey,
		OneTimeEncryptedKeyEncryptionKey:  a.OneTimeEncryptedKeyEncryptionKey,
	}
}

func importAccountUserRelationship(a *persistence.AccountUserRelationship) AccountUserRelationship {
	return AccountUserRelationship{
		RelationshipID:                    a.RelationshipID,
		AccountUserID:                     a.AccountUserID,
		AccountID:                         a.AccountID,
		PasswordEncryptedKeyEncryptionKey: a.PasswordEncryptedKeyEncryptionKey,
		EmailEncryptedKeyEncryptionKey:    a.EmailEncryptedKeyEncryptionKey,
		OneTimeEncryptedKeyEncryptionKey:  a.OneTimeEncryptedKeyEncryptionKey,
	}
}

func (a *Account) export() persistence.Account {
	var events []persistence.Event
	for _, e := range a.Events {
		events = append(events, e.export())
	}
	return persistence.Account{
		AccountID:           a.AccountID,
		Name:                a.Name,
		PublicKey:           a.PublicKey,
		EncryptedPrivateKey: a.EncryptedPrivateKey,
		UserSalt:            a.UserSalt,
		Retired:             a.Retired,
		Created:             a.Created,
		Events:              events,
		AccountStyles:       a.AccountStyles,
	}
}

func importAccount(a *persistence.Account) Account {
	events := []Event{}
	for _, e := range a.Events {
		events = append(events, importEvent(&e))
	}
	return Account{
		AccountID:           a.AccountID,
		Name:                a.Name,
		PublicKey:           a.PublicKey,
		EncryptedPrivateKey: a.EncryptedPrivateKey,
		UserSalt:            a.UserSalt,
		Retired:             a.Retired,
		Created:             a.Created,
		Events:              events,
		AccountStyles:       a.AccountStyles,
	}
}
