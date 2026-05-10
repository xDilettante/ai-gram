package telegram

import (
	"encoding/json"
	stderrors "errors"
)

// ChatMemberResult decodes an official ChatMember JSON union.
type ChatMemberResult struct {
	ChatMember
}

// UnmarshalJSON decodes an official ChatMember variant by status.
func (r *ChatMemberResult) UnmarshalJSON(data []byte) error {
	member, err := UnmarshalChatMember(data)
	if err != nil {
		return err
	}
	r.ChatMember = member
	return nil
}

// UnmarshalChatMember decodes an official ChatMember variant by status.
func UnmarshalChatMember(data []byte) (ChatMember, error) {
	var meta struct {
		Status ChatMemberStatus `json:"status"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Status {
	case ChatMemberStatusCreator:
		var member ChatMemberOwner
		if err := json.Unmarshal(data, &member); err != nil {
			return nil, err
		}
		return member, nil
	case ChatMemberStatusAdministrator:
		var member ChatMemberAdministrator
		if err := json.Unmarshal(data, &member); err != nil {
			return nil, err
		}
		return member, nil
	case ChatMemberStatusMember:
		var member ChatMemberMember
		if err := json.Unmarshal(data, &member); err != nil {
			return nil, err
		}
		return member, nil
	case ChatMemberStatusRestricted:
		var member ChatMemberRestricted
		if err := json.Unmarshal(data, &member); err != nil {
			return nil, err
		}
		return member, nil
	case ChatMemberStatusLeft:
		var member ChatMemberLeft
		if err := json.Unmarshal(data, &member); err != nil {
			return nil, err
		}
		return member, nil
	case ChatMemberStatusKicked:
		var member ChatMemberBanned
		if err := json.Unmarshal(data, &member); err != nil {
			return nil, err
		}
		return member, nil
	default:
		return nil, stderrors.New("unsupported chat member status")
	}
}

// UnmarshalJSON decodes polymorphic old/new chat member fields.
func (u *ChatMemberUpdated) UnmarshalJSON(data []byte) error {
	type chatMemberUpdatedAlias ChatMemberUpdated
	payload := struct {
		OldChatMember json.RawMessage `json:"old_chat_member"`
		NewChatMember json.RawMessage `json:"new_chat_member"`
		*chatMemberUpdatedAlias
	}{chatMemberUpdatedAlias: (*chatMemberUpdatedAlias)(u)}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.OldChatMember) > 0 {
		oldMember, err := UnmarshalChatMember(payload.OldChatMember)
		if err != nil {
			return err
		}
		u.OldChatMember = oldMember
	}
	if len(payload.NewChatMember) > 0 {
		newMember, err := UnmarshalChatMember(payload.NewChatMember)
		if err != nil {
			return err
		}
		u.NewChatMember = newMember
	}
	return nil
}

func (m ChatMemberOwner) ChatMemberStatus() ChatMemberStatus { return m.Status }
func (m ChatMemberOwner) ChatMemberUser() User               { return m.User }
func (m ChatMemberOwner) isChatMember()                      {}

func (m ChatMemberAdministrator) ChatMemberStatus() ChatMemberStatus { return m.Status }
func (m ChatMemberAdministrator) ChatMemberUser() User               { return m.User }
func (m ChatMemberAdministrator) isChatMember()                      {}

func (m ChatMemberMember) ChatMemberStatus() ChatMemberStatus { return m.Status }
func (m ChatMemberMember) ChatMemberUser() User               { return m.User }
func (m ChatMemberMember) isChatMember()                      {}

func (m ChatMemberRestricted) ChatMemberStatus() ChatMemberStatus { return m.Status }
func (m ChatMemberRestricted) ChatMemberUser() User               { return m.User }
func (m ChatMemberRestricted) isChatMember()                      {}

func (m ChatMemberLeft) ChatMemberStatus() ChatMemberStatus { return m.Status }
func (m ChatMemberLeft) ChatMemberUser() User               { return m.User }
func (m ChatMemberLeft) isChatMember()                      {}

func (m ChatMemberBanned) ChatMemberStatus() ChatMemberStatus { return m.Status }
func (m ChatMemberBanned) ChatMemberUser() User               { return m.User }
func (m ChatMemberBanned) isChatMember()                      {}
