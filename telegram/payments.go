package telegram

// LabeledPrice represents a price component in the smallest units of a currency.
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int64  `json:"amount"`
}
