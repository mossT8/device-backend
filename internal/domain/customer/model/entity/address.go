package entity

import "time"

type addressRecordId int64
type addressText string
type addressBool bool
type addressDate time.Time

type Address struct {
	ID addressRecordId

	AccountId    accountRecordId
	Name         addressText
	AddressLine1 addressText
	AddressLine2 *addressText
	City         addressText
	State        *addressText
	PostalCode   addressText
	Country      addressText
	Verified     addressBool

	CreatedAt  addressDate
	ModifiedAt addressDate
}

func NewAddress(accountId int64, timestamp time.Time) Address {

	return Address{
		AccountId:  accountRecordId(accountId),
		CreatedAt:  addressDate(timestamp),
		ModifiedAt: addressDate(timestamp),
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
	a.ID = addressRecordId(id)
}

func (a *Address) SetAccountId(accountId int64) {
	a.AccountId = accountRecordId(accountId)
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetName(name string) {
	a.Name = addressText(name)
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetAddressLine1(addressLine1 string) {
	a.AddressLine1 = addressText(addressLine1)
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetAddressLine2(addressLine2 string) {
	at := addressText(addressLine2)
	a.AddressLine2 = &at
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetCity(city string) {
	a.City = addressText(city)
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetState(state string) {
	st := addressText(state)
	a.State = &st
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetPostalCode(postalCode string) {
	a.PostalCode = addressText(postalCode)
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetCountry(country string) {
	a.Country = addressText(country)
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetVerified(verified bool) {
	a.Verified = addressBool(verified)
	a.ModifiedAt = addressDate(time.Now())
}

func (a *Address) SetCreatedAt(createdAt time.Time) {
	a.CreatedAt = addressDate(createdAt)
}

func (a *Address) SetModifiedAt(modifiedAt time.Time) {
	a.ModifiedAt = addressDate(modifiedAt)
}
