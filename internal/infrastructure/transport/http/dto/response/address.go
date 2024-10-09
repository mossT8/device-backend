package response

import "time"

type Address struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	AddressLine1 string    `json:"addressLine1"`
	AddressLine2 string    `json:"addressLine2"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	PostalCode   string    `json:"postalCode"`
	Country      string    `json:"country"`
	CreatedAt    time.Time `json:"createdAt"`
	ModifiedAt   time.Time `json:"modifiedAt"`
}
