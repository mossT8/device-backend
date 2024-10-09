package entity

import (
	"time"
)

type Address struct {
	ID mysqlRecordId

	AccountId    mysqlRecordId
	Name         mysqlText
	AddressLine1 mysqlText
	AddressLine2 *mysqlText
	City         mysqlText
	State        *mysqlText
	PostalCode   mysqlText
	Country      mysqlText
	Verified     mysqlBool

	CreatedAt  mysqlDate
	ModifiedAt mysqlDate
}

func NewAddress(accountId int64, timestamp time.Time) Address {

	return Address{
		AccountId:  mysqlRecordId(accountId),
		CreatedAt:  mysqlDate(timestamp),
		ModifiedAt: mysqlDate(timestamp),
	}
}

// Getters
func (a *Address) GetID() int64 {
	return int64(a.ID)
}

func (a *Address) GetAccountId() int64 {
	return int64(a.AccountId)
}

func (a *Address) GetName() string {
	return string(a.Name)
}

func (a *Address) GetAddressLine1() string {
	return string(a.AddressLine1)
}

func (a *Address) GetAddressLine2() string {
	if a.AddressLine2 != nil {
		return string(*a.AddressLine2)
	}
	return ""
}

func (a *Address) GetCity() string {
	return string(a.City)
}

func (a *Address) GetState() string {
	if a.State != nil {
		return string(*a.State)
	}
	return ""
}

func (a *Address) GetPostalCode() string {
	return string(a.PostalCode)
}

func (a *Address) GetCountry() string {
	return string(a.Country)
}

func (a *Address) GetVerified() bool {
	return bool(a.Verified)
}

func (a *Address) GetCreatedAt() time.Time {
	return time.Time(a.CreatedAt)
}

func (a *Address) GetModifiedAt() time.Time {
	return time.Time(a.ModifiedAt)
}

// Setters
func (a *Address) SetID(id int64) {
	a.ID = mysqlRecordId(id)
}

func (a *Address) SetAccountId(accountId int64) {
	a.AccountId = mysqlRecordId(accountId)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetName(name string) {
	a.Name = mysqlText(name)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetAddressLine1(addressLine1 string) {
	a.AddressLine1 = mysqlText(addressLine1)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetAddressLine2(addressLine2 string) {
	at := mysqlText(addressLine2)
	a.AddressLine2 = &at
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetCity(city string) {
	a.City = mysqlText(city)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetState(state string) {
	st := mysqlText(state)
	a.State = &st
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetPostalCode(postalCode string) {
	a.PostalCode = mysqlText(postalCode)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetCountry(country string) {
	a.Country = mysqlText(country)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetVerified(verified bool) {
	a.Verified = mysqlBool(verified)
	a.ModifiedAt = mysqlDate(time.Now())
}

func (a *Address) SetCreatedAt(createdAt time.Time) {
	a.CreatedAt = mysqlDate(createdAt)
}

func (a *Address) SetModifiedAt(modifiedAt time.Time) {
	a.ModifiedAt = mysqlDate(modifiedAt)
}
