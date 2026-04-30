package telegram

// Game represents a Telegram game created through BotFather.
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	Animation    *Animation      `json:"animation,omitempty"`
}

// CallbackGame is a placeholder object used by game-launching inline keyboard buttons.
type CallbackGame struct{}

// GameHighScore represents one row of a game's high score table.
type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}
