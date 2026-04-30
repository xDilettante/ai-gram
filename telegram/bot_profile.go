package telegram

// BotName describes a localized bot name.
type BotName struct {
	Name string `json:"name"`
}

// BotDescription describes a localized bot description.
type BotDescription struct {
	Description string `json:"description"`
}

// BotShortDescription describes a localized bot short description.
type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}
