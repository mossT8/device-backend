package entity

import (
	"time"
)

type User struct {
	ID mysqlRecordId

	AccountId       mysqlRecordId
	Email           mysqlText
	Cell            mysqlText
	FirstName       mysqlText
	LastName        mysqlText
	Verified        mysqlBool
	ReceivesUpdates mysqlBool

	CreatedAt  mysqlDate
	ModifiedAt mysqlDate
}

func NewUser(accountId int64, email string, timestamp time.Time) User {

	return User{
		AccountId:       mysqlRecordId(accountId),
		Email:           mysqlText(email),
		Verified:        mysqlBool(false),
		ReceivesUpdates: mysqlBool(false),
		CreatedAt:       mysqlDate(timestamp),
		ModifiedAt:      mysqlDate(timestamp),
	}
}

// Getters
func (u *User) GetID() int64 {
	return int64(u.ID)
}

func (u *User) GetAccountId() int64 {
	return int64(u.AccountId)
}

func (u *User) GetEmail() string {
	return string(u.Email)
}

func (u *User) GetCell() string {
	return string(u.Cell)
}

func (u *User) GetFirstName() string {
	return string(u.FirstName)
}

func (u *User) GetLastName() string {
	return string(u.LastName)
}

func (u *User) GetVerified() bool {
	return bool(u.Verified)
}

func (u *User) GetReceivesUpdates() bool {
	return bool(u.ReceivesUpdates)
}

func (u *User) GetCreatedAt() time.Time {
	return time.Time(u.CreatedAt)
}

func (u *User) GetModifiedAt() time.Time {
	return time.Time(u.ModifiedAt)
}

// Setters
func (u *User) SetID(id int64) {
	u.ID = mysqlRecordId(id)
}

func (u *User) SetAccountId(accountId int64) {
	u.AccountId = mysqlRecordId(accountId)
	u.ModifiedAt = mysqlDate(time.Now())
}

func (u *User) SetEmail(email string) {
	u.Email = mysqlText(email)
	u.ModifiedAt = mysqlDate(time.Now())
}

func (u *User) SetCell(cell string) {
	u.Cell = mysqlText(cell)
	u.ModifiedAt = mysqlDate(time.Now())
}

func (u *User) SetFirstName(firstName string) {
	u.FirstName = mysqlText(firstName)
	u.ModifiedAt = mysqlDate(time.Now())
}

func (u *User) SetLastName(lastName string) {
	u.LastName = mysqlText(lastName)
	u.ModifiedAt = mysqlDate(time.Now())
}

func (u *User) SetVerified(verified bool) {
	u.Verified = mysqlBool(verified)
	u.ModifiedAt = mysqlDate(time.Now())
}

func (u *User) SetReceivesUpdates(receivesUpdates bool) {
	u.ReceivesUpdates = mysqlBool(receivesUpdates)
	u.ModifiedAt = mysqlDate(time.Now())
}

func (u *User) SetCreatedAt(createdAt time.Time) {
	u.CreatedAt = mysqlDate(createdAt)
}

func (u *User) SetModifiedAt(modifiedAt time.Time) {
	u.ModifiedAt = mysqlDate(modifiedAt)
}
