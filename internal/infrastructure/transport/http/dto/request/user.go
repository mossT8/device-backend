package request

type User struct {
	AccountId       int64  `json:"accountID"`
	Email           string `json:"email"`
	Cell            string `json:"cell"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Verified        bool   `json:"verified"`
	ReceivesUpdates bool   `json:"receivesUpdates"`
}
