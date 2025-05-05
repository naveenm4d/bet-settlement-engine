package entities

type Account struct {
	UserID  string `json:"user_id"`
	Balance int64  `json:"balance"`
}
