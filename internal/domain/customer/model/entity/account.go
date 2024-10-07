package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type accountRecordId int64
type accountText string
type accountBool bool
type accountDate time.Time

type Account struct {
	ID accountRecordId

	Email           accountText
	PasswordHash    accountText
	Salt            accountText
	Name            accountText
	Verified        accountBool
	ReceivesUpdates accountBool

	CreatedAt  accountDate
	ModifiedAt accountDate
}

func NewAccount(email, name string, timestamp time.Time) Account {
	return Account{
		Email:      accountText(email),
		Name:       accountText(name),
		CreatedAt:  accountDate(timestamp),
		ModifiedAt: accountDate(timestamp),
	}
}

func (a *Account) GetID() int64 {
	return int64(a.ID)
}

func (a *Account) GetEmail() string {
	return string(a.Email)
}

func (a *Account) GetPasswordHash() string {
	return string(a.PasswordHash)
}

func (a *Account) GetSalt() string {
	return string(a.Salt)
}

func (a *Account) GetName() string {
	return string(a.Name)
}

func (a *Account) GetVerified() bool {
	return bool(a.Verified)
}

func (a *Account) GetReceivesUpdates() bool {
	return bool(a.ReceivesUpdates)
}

func (a *Account) GetCreatedAt() time.Time {
	return time.Time(a.CreatedAt)
}

func (a *Account) GetModifiedAt() time.Time {
	return time.Time(a.ModifiedAt)
}

func (a *Account) SetID(id int64) {
	a.ID = accountRecordId(id)
}

func (a *Account) SetEmail(email string) {
	a.Email = accountText(email)
	a.ModifiedAt = accountDate(time.Now())
}

func (a *Account) SetPassword(password, salt string) error {
	hashBytes, hErr := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if hErr != nil {
		return hErr
	}
	a.PasswordHash = accountText(string(hashBytes))
	a.Salt = accountText(salt)
	a.ModifiedAt = accountDate(time.Now())
	return nil
}

func (a *Account) SetPasswordHash(passwordHash string) {
	a.PasswordHash = accountText(passwordHash)
	a.ModifiedAt = accountDate(time.Now())
}

func (a *Account) SetSalt(salt string) {
	a.Salt = accountText(salt)
	a.ModifiedAt = accountDate(time.Now())
}

func (a *Account) SetName(name string) {
	a.Name = accountText(name)
	a.ModifiedAt = accountDate(time.Now())
}

func (a *Account) SetVerified(verified bool) {
	a.Verified = accountBool(verified)
	a.ModifiedAt = accountDate(time.Now())
}

func (a *Account) SetReceivesUpdates(receivesUpdates bool) {
	a.ReceivesUpdates = accountBool(receivesUpdates)
	a.ModifiedAt = accountDate(time.Now())
}

func (a *Account) SetCreatedAt(createdAt time.Time) {
	a.CreatedAt = accountDate(createdAt)
}

func (a *Account) SetModifiedAt(modifiedAt time.Time) {
	a.ModifiedAt = accountDate(modifiedAt)
}
