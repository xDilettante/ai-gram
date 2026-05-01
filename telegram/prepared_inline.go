package telegram

// PreparedInlineMessage describes an inline message saved for a Mini App user.
type PreparedInlineMessage struct {
	ID             string `json:"id"`
	ExpirationDate int64  `json:"expiration_date"`
}
