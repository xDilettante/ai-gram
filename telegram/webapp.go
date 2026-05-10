package telegram

// WebAppData describes data sent from a Web App to the bot.
type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}

// WriteAccessAllowed describes a service message about Web App write access.
type WriteAccessAllowed struct {
	FromRequest        bool   `json:"from_request,omitempty"`
	WebAppName         string `json:"web_app_name,omitempty"`
	FromAttachmentMenu bool   `json:"from_attachment_menu,omitempty"`
}

// SentWebAppMessage describes an inline message sent by a Web App on behalf of a user.
type SentWebAppMessage struct {
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// SentGuestMessage describes an inline message sent by a guest bot.
type SentGuestMessage struct {
	InlineMessageID string `json:"inline_message_id"`
}
