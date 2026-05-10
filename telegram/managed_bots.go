package telegram

// ManagedBotCreated contains information about a bot created to be managed by the current bot.
type ManagedBotCreated struct {
	Bot User `json:"bot"`
}

// ManagedBotUpdated contains information about managed bot creation, token replacement, or owner changes.
type ManagedBotUpdated struct {
	User User `json:"user"`
	Bot  User `json:"bot"`
}

// BotAccessSettings describes the access settings of a bot.
type BotAccessSettings struct {
	IsAccessRestricted bool   `json:"is_access_restricted"`
	AddedUsers         []User `json:"added_users,omitempty"`
}

// PreparedKeyboardButton describes a saved keyboard button that can be used by a Mini App user.
type PreparedKeyboardButton struct {
	ID string `json:"id"`
}
