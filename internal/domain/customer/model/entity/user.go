package entity

import (
	"time"
)

type userRecordId int64
type userText string
type userBool bool
type userDate time.Time

type User struct {
	ID userRecordId

	AccountId       accountRecordId
	Email           userText
	Cell            userText
	FirstName       *userText
	LastName        *userText
	Verified        userBool
	ReceivesUpdates userBool

	CreatedAt  userDate
	ModifiedAt userDate
}

func NewUser(accountId int64, email string, timestamp time.Time) User {

	return User{
		AccountId:       accountRecordId(accountId),
		Email:           userText(email),
		Verified:        userBool(false),
		ReceivesUpdates: userBool(false),
		CreatedAt:       userDate(timestamp),
		ModifiedAt:      userDate(timestamp),
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
	if u.FirstName != nil {
		return string(*u.FirstName)
	}
	return ""
}

func (u *User) GetLastName() string {
	if u.LastName != nil {
		return string(*u.LastName)
	}
	return ""
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
	u.ID = userRecordId(id)
}

func (u *User) SetAccountId(accountId int64) {
	u.AccountId = accountRecordId(accountId)
	u.ModifiedAt = userDate(time.Now())
}

func (u *User) SetEmail(email string) {
	u.Email = userText(email)
	u.ModifiedAt = userDate(time.Now())
}

func (u *User) SetCell(cell string) {
	u.Cell = userText(cell)
	u.ModifiedAt = userDate(time.Now())
}

func (u *User) SetFirstName(firstName string) {
	ft := userText(firstName)
	u.FirstName = &ft
	u.ModifiedAt = userDate(time.Now())
}

func (u *User) SetLastName(lastName string) {
	lt := userText(lastName)
	u.LastName = &lt
	u.ModifiedAt = userDate(time.Now())
}

func (u *User) SetVerified(verified bool) {
	u.Verified = userBool(verified)
	u.ModifiedAt = userDate(time.Now())
}

func (u *User) SetReceivesUpdates(receivesUpdates bool) {
	u.ReceivesUpdates = userBool(receivesUpdates)
	u.ModifiedAt = userDate(time.Now())
}

func (u *User) SetCreatedAt(createdAt time.Time) {
	u.CreatedAt = userDate(createdAt)
}

func (u *User) SetModifiedAt(modifiedAt time.Time) {
	u.ModifiedAt = userDate(modifiedAt)
}
