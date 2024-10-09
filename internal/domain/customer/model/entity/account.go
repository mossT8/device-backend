package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID mysqlRecordId

	Email           mysqlText
	PasswordHash    mysqlText
	Salt            mysqlText
	Name            mysqlText
	Verified        mysqlBool
	ReceivesUpdates mysqlBool

	CreatedAt  mysqlDate
	ModifiedAt mysqlDate
}

func NewAccount(email, name string, timestamp time.Time) Account {
	return Account{
		Email:      mysqlText(email),
		Name:       mysqlText(name),
		CreatedAt:  mysqlDate(timestamp),
		ModifiedAt: mysqlDate(timestamp),
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
	a.ID = mysqlRecordId(id)
}

func (a *Account) SetEmail(email string) {
	a.Email = mysqlText(email)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Account) SetPassword(password, salt string) error {
	hashBytes, hErr := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if hErr != nil {
		return hErr
	}
	a.PasswordHash = mysqlText(string(hashBytes))
	a.Salt = mysqlText(salt)
	a.ModifiedAt = mysqlDate(time.Now())
	return nil
}

func (a *Account) SetPasswordHash(passwordHash string) {
	a.PasswordHash = mysqlText(passwordHash)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Account) SetSalt(salt string) {
	a.Salt = mysqlText(salt)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Account) SetName(name string) {
	a.Name = mysqlText(name)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Account) SetVerified(verified bool) {
	a.Verified = mysqlBool(verified)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Account) SetReceivesUpdates(receivesUpdates bool) {
	a.ReceivesUpdates = mysqlBool(receivesUpdates)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Account) SetCreatedAt(createdAt time.Time) {
	a.CreatedAt = mysqlDate(createdAt)
}

func (a *Account) SetModifiedAt(modifiedAt time.Time) {
	a.ModifiedAt = mysqlDate(modifiedAt)
}
