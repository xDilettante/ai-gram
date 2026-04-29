package bot

import (
	"encoding/json"
	stderrors "errors"
	"strings"
)

type chatIDKind uint8

const (
	chatIDUnset chatIDKind = iota
	chatIDInt
	chatIDString
)

// ChatID identifies a Telegram chat by numeric ID or username string.
type ChatID struct {
	kind     chatIDKind
	intID    int64
	stringID string
}

// ChatIDInt creates a numeric chat ID.
func ChatIDInt(id int64) ChatID {
	return ChatID{kind: chatIDInt, intID: id}
}

// ChatIDString creates a string chat ID, such as a channel username.
func ChatIDString(id string) ChatID {
	return ChatID{kind: chatIDString, stringID: id}
}

// MarshalJSON encodes ChatID as either a JSON number or a JSON string.
func (id ChatID) MarshalJSON() ([]byte, error) {
	if !id.valid() {
		return nil, stderrors.New("chat_id is required")
	}

	switch id.kind {
	case chatIDInt:
		return json.Marshal(id.intID)
	case chatIDString:
		return json.Marshal(id.stringID)
	default:
		return nil, stderrors.New("chat_id is required")
	}
}

func (id ChatID) valid() bool {
	switch id.kind {
	case chatIDInt:
		return id.intID != 0
	case chatIDString:
		return strings.TrimSpace(id.stringID) != ""
	default:
		return false
	}
}
