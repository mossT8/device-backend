package request

type Account struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	ReceivesUpdates bool   `json:"receivesUpdates"`
}
