package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
)

const (
	richTextBoldType                   = "bold"
	richTextItalicType                 = "italic"
	richTextUnderlineType              = "underline"
	richTextStrikethroughType          = "strikethrough"
	richTextSpoilerType                = "spoiler"
	richTextDateTimeType               = "date_time"
	richTextTextMentionType            = "text_mention"
	richTextSubscriptType              = "subscript"
	richTextSuperscriptType            = "superscript"
	richTextMarkedType                 = "marked"
	richTextCodeType                   = "code"
	richTextCustomEmojiType            = "custom_emoji"
	richTextMathematicalExpressionType = "mathematical_expression"
	richTextURLType                    = "url"
	richTextEmailAddressType           = "email_address"
	richTextPhoneNumberType            = "phone_number"
	richTextBankCardNumberType         = "bank_card_number"
	richTextMentionType                = "mention"
	richTextHashtagType                = "hashtag"
	richTextCashtagType                = "cashtag"
	richTextBotCommandType             = "bot_command"
	richTextAnchorType                 = "anchor"
	richTextAnchorLinkType             = "anchor_link"
	richTextReferenceType              = "reference"
	richTextReferenceLinkType          = "reference_link"

	richBlockParagraphType              = "paragraph"
	richBlockSectionHeadingType         = "heading"
	richBlockPreformattedType           = "pre"
	richBlockFooterType                 = "footer"
	richBlockDividerType                = "divider"
	richBlockMathematicalExpressionType = "mathematical_expression"
	richBlockAnchorType                 = "anchor"
	richBlockListType                   = "list"
	richBlockBlockQuotationType         = "blockquote"
	richBlockPullQuotationType          = "pullquote"
	richBlockCollageType                = "collage"
	richBlockSlideshowType              = "slideshow"
	richBlockTableType                  = "table"
	richBlockDetailsType                = "details"
	richBlockMapType                    = "map"
	richBlockAnimationType              = "animation"
	richBlockAudioType                  = "audio"
	richBlockPhotoType                  = "photo"
	richBlockVideoType                  = "video"
	richBlockVoiceNoteType              = "voice_note"
	richBlockThinkingType               = "thinking"
)

// RichMessage represents a rich formatted message.
type RichMessage struct {
	Blocks []RichBlock `json:"blocks"`
	IsRTL  bool        `json:"is_rtl,omitempty"`
}

// RichText marks Telegram rich text objects.
type RichText interface {
	richText()
}

// RichTextPlain represents a plain string rich text value.
type RichTextPlain string

// RichTextList represents an array of rich text values.
type RichTextList []RichText

// RichTextBold represents bold rich text.
type RichTextBold struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextItalic represents italic rich text.
type RichTextItalic struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextUnderline represents underlined rich text.
type RichTextUnderline struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextStrikethrough represents strikethrough rich text.
type RichTextStrikethrough struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextSpoiler represents spoiler rich text.
type RichTextSpoiler struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextDateTime represents a formatted date and time.
type RichTextDateTime struct {
	Type           string   `json:"type"`
	Text           RichText `json:"text"`
	UnixTime       int64    `json:"unix_time"`
	DateTimeFormat string   `json:"date_time_format"`
}

// RichTextTextMention represents a mention of a Telegram user by identifier.
type RichTextTextMention struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
	User User     `json:"user"`
}

// RichTextSubscript represents subscript rich text.
type RichTextSubscript struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextSuperscript represents superscript rich text.
type RichTextSuperscript struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextMarked represents marked rich text.
type RichTextMarked struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextCode represents monowidth rich text.
type RichTextCode struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichTextCustomEmoji represents a custom emoji.
type RichTextCustomEmoji struct {
	Type            string `json:"type"`
	CustomEmojiID   string `json:"custom_emoji_id"`
	AlternativeText string `json:"alternative_text"`
}

// RichTextMathematicalExpression represents an inline mathematical expression.
type RichTextMathematicalExpression struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

// RichTextURL represents rich text with a link.
type RichTextURL struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
	URL  string   `json:"url"`
}

// RichTextEmailAddress represents rich text with an email address.
type RichTextEmailAddress struct {
	Type         string   `json:"type"`
	Text         RichText `json:"text"`
	EmailAddress string   `json:"email_address"`
}

// RichTextPhoneNumber represents rich text with a phone number.
type RichTextPhoneNumber struct {
	Type        string   `json:"type"`
	Text        RichText `json:"text"`
	PhoneNumber string   `json:"phone_number"`
}

// RichTextBankCardNumber represents rich text with a bank card number.
type RichTextBankCardNumber struct {
	Type           string   `json:"type"`
	Text           RichText `json:"text"`
	BankCardNumber string   `json:"bank_card_number"`
}

// RichTextMention represents a username mention.
type RichTextMention struct {
	Type     string   `json:"type"`
	Text     RichText `json:"text"`
	Username string   `json:"username"`
}

// RichTextHashtag represents a hashtag.
type RichTextHashtag struct {
	Type    string   `json:"type"`
	Text    RichText `json:"text"`
	Hashtag string   `json:"hashtag"`
}

// RichTextCashtag represents a cashtag.
type RichTextCashtag struct {
	Type    string   `json:"type"`
	Text    RichText `json:"text"`
	Cashtag string   `json:"cashtag"`
}

// RichTextBotCommand represents a bot command.
type RichTextBotCommand struct {
	Type       string   `json:"type"`
	Text       RichText `json:"text"`
	BotCommand string   `json:"bot_command"`
}

// RichTextAnchor represents an anchor.
type RichTextAnchor struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// RichTextAnchorLink represents a link to an anchor.
type RichTextAnchorLink struct {
	Type       string   `json:"type"`
	Text       RichText `json:"text"`
	AnchorName string   `json:"anchor_name"`
}

// RichTextReference represents a reference.
type RichTextReference struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
	Name string   `json:"name"`
}

// RichTextReferenceLink represents a link to a reference.
type RichTextReferenceLink struct {
	Type          string   `json:"type"`
	Text          RichText `json:"text"`
	ReferenceName string   `json:"reference_name"`
}

func (RichTextPlain) richText()                  {}
func (RichTextList) richText()                   {}
func (RichTextBold) richText()                   {}
func (RichTextItalic) richText()                 {}
func (RichTextUnderline) richText()              {}
func (RichTextStrikethrough) richText()          {}
func (RichTextSpoiler) richText()                {}
func (RichTextDateTime) richText()               {}
func (RichTextTextMention) richText()            {}
func (RichTextSubscript) richText()              {}
func (RichTextSuperscript) richText()            {}
func (RichTextMarked) richText()                 {}
func (RichTextCode) richText()                   {}
func (RichTextCustomEmoji) richText()            {}
func (RichTextMathematicalExpression) richText() {}
func (RichTextURL) richText()                    {}
func (RichTextEmailAddress) richText()           {}
func (RichTextPhoneNumber) richText()            {}
func (RichTextBankCardNumber) richText()         {}
func (RichTextMention) richText()                {}
func (RichTextHashtag) richText()                {}
func (RichTextCashtag) richText()                {}
func (RichTextBotCommand) richText()             {}
func (RichTextAnchor) richText()                 {}
func (RichTextAnchorLink) richText()             {}
func (RichTextReference) richText()              {}
func (RichTextReferenceLink) richText()          {}

// RichBlockCaption represents a rich block caption and optional credit.
type RichBlockCaption struct {
	Text   RichText `json:"text"`
	Credit RichText `json:"credit,omitempty"`
}

// RichBlockTableCell represents a rich table cell.
type RichBlockTableCell struct {
	Text     RichText `json:"text,omitempty"`
	IsHeader bool     `json:"is_header,omitempty"`
	Colspan  int      `json:"colspan,omitempty"`
	Rowspan  int      `json:"rowspan,omitempty"`
	Align    string   `json:"align,omitempty"`
	Valign   string   `json:"valign,omitempty"`
}

// RichBlockListItem represents an item of a rich list block.
type RichBlockListItem struct {
	Label       string      `json:"label"`
	Blocks      []RichBlock `json:"blocks"`
	HasCheckbox bool        `json:"has_checkbox,omitempty"`
	IsChecked   bool        `json:"is_checked,omitempty"`
	Value       int         `json:"value,omitempty"`
	Type        string      `json:"type,omitempty"`
}

// RichBlock marks Telegram rich message block objects.
type RichBlock interface {
	richBlock()
}

// RichBlockParagraph represents a rich text paragraph.
type RichBlockParagraph struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichBlockSectionHeading represents a rich section heading.
type RichBlockSectionHeading struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
	Size int      `json:"size"`
}

// RichBlockPreformatted represents a rich preformatted block.
type RichBlockPreformatted struct {
	Type     string   `json:"type"`
	Text     RichText `json:"text"`
	Language string   `json:"language,omitempty"`
}

// RichBlockFooter represents a rich footer block.
type RichBlockFooter struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

// RichBlockDivider represents a rich divider block.
type RichBlockDivider struct {
	Type string `json:"type"`
}

// RichBlockMathematicalExpression represents a block mathematical expression.
type RichBlockMathematicalExpression struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

// RichBlockAnchor represents a rich block anchor.
type RichBlockAnchor struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// RichBlockList represents a rich list block.
type RichBlockList struct {
	Type  string              `json:"type"`
	Items []RichBlockListItem `json:"items"`
}

// RichBlockBlockQuotation represents a block quotation.
type RichBlockBlockQuotation struct {
	Type   string      `json:"type"`
	Blocks []RichBlock `json:"blocks"`
	Credit RichText    `json:"credit,omitempty"`
}

// RichBlockPullQuotation represents a pull quotation.
type RichBlockPullQuotation struct {
	Type   string   `json:"type"`
	Text   RichText `json:"text"`
	Credit RichText `json:"credit,omitempty"`
}

// RichBlockCollage represents a rich collage block.
type RichBlockCollage struct {
	Type    string            `json:"type"`
	Blocks  []RichBlock       `json:"blocks"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockSlideshow represents a rich slideshow block.
type RichBlockSlideshow struct {
	Type    string            `json:"type"`
	Blocks  []RichBlock       `json:"blocks"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockTable represents a rich table block.
type RichBlockTable struct {
	Type       string                 `json:"type"`
	Cells      [][]RichBlockTableCell `json:"cells"`
	IsBordered bool                   `json:"is_bordered,omitempty"`
	IsStriped  bool                   `json:"is_striped,omitempty"`
	Caption    RichText               `json:"caption,omitempty"`
}

// RichBlockDetails represents an expandable details block.
type RichBlockDetails struct {
	Type    string      `json:"type"`
	Summary RichText    `json:"summary"`
	Blocks  []RichBlock `json:"blocks"`
	IsOpen  bool        `json:"is_open,omitempty"`
}

// RichBlockMap represents a rich map block.
type RichBlockMap struct {
	Type     string            `json:"type"`
	Location Location          `json:"location"`
	Zoom     int               `json:"zoom"`
	Width    int               `json:"width"`
	Height   int               `json:"height"`
	Caption  *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockAnimation represents a rich animation block.
type RichBlockAnimation struct {
	Type       string            `json:"type"`
	Animation  Animation         `json:"animation"`
	HasSpoiler bool              `json:"has_spoiler,omitempty"`
	Caption    *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockAudio represents a rich audio block.
type RichBlockAudio struct {
	Type    string            `json:"type"`
	Audio   Audio             `json:"audio"`
	Caption *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockPhoto represents a rich photo block.
type RichBlockPhoto struct {
	Type       string            `json:"type"`
	Photo      []PhotoSize       `json:"photo"`
	HasSpoiler bool              `json:"has_spoiler,omitempty"`
	Caption    *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockVideo represents a rich video block.
type RichBlockVideo struct {
	Type       string            `json:"type"`
	Video      Video             `json:"video"`
	HasSpoiler bool              `json:"has_spoiler,omitempty"`
	Caption    *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockVoiceNote represents a rich voice note block.
type RichBlockVoiceNote struct {
	Type      string            `json:"type"`
	VoiceNote Voice             `json:"voice_note"`
	Caption   *RichBlockCaption `json:"caption,omitempty"`
}

// RichBlockThinking represents a rich draft-only thinking block.
type RichBlockThinking struct {
	Type string   `json:"type"`
	Text RichText `json:"text"`
}

func (RichBlockParagraph) richBlock()              {}
func (RichBlockSectionHeading) richBlock()         {}
func (RichBlockPreformatted) richBlock()           {}
func (RichBlockFooter) richBlock()                 {}
func (RichBlockDivider) richBlock()                {}
func (RichBlockMathematicalExpression) richBlock() {}
func (RichBlockAnchor) richBlock()                 {}
func (RichBlockList) richBlock()                   {}
func (RichBlockBlockQuotation) richBlock()         {}
func (RichBlockPullQuotation) richBlock()          {}
func (RichBlockCollage) richBlock()                {}
func (RichBlockSlideshow) richBlock()              {}
func (RichBlockTable) richBlock()                  {}
func (RichBlockDetails) richBlock()                {}
func (RichBlockMap) richBlock()                    {}
func (RichBlockAnimation) richBlock()              {}
func (RichBlockAudio) richBlock()                  {}
func (RichBlockPhoto) richBlock()                  {}
func (RichBlockVideo) richBlock()                  {}
func (RichBlockVoiceNote) richBlock()              {}
func (RichBlockThinking) richBlock()               {}

// UnmarshalJSON decodes RichMessage blocks as polymorphic rich block objects.
func (m *RichMessage) UnmarshalJSON(data []byte) error {
	var payload struct {
		Blocks []json.RawMessage `json:"blocks"`
		IsRTL  bool              `json:"is_rtl"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	blocks, err := decodeRichBlockSlice(payload.Blocks)
	if err != nil {
		return err
	}
	m.Blocks = blocks
	m.IsRTL = payload.IsRTL
	return nil
}

// UnmarshalRichText decodes a polymorphic Telegram RichText value.
func UnmarshalRichText(data []byte) (RichText, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	switch trimmed[0] {
	case '"':
		var text string
		if err := json.Unmarshal(trimmed, &text); err != nil {
			return nil, err
		}
		return RichTextPlain(text), nil
	case '[':
		var rawItems []json.RawMessage
		if err := json.Unmarshal(trimmed, &rawItems); err != nil {
			return nil, err
		}
		items := make(RichTextList, 0, len(rawItems))
		for _, raw := range rawItems {
			item, err := UnmarshalRichText(raw)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		return items, nil
	case '{':
		return unmarshalRichTextObject(trimmed)
	default:
		return nil, stderrors.New("rich text must be a string, array, object, or null")
	}
}

// UnmarshalRichBlock decodes a polymorphic Telegram RichBlock object.
func UnmarshalRichBlock(data []byte) (RichBlock, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}
	if trimmed[0] != '{' {
		return nil, stderrors.New("rich block must be an object")
	}
	return unmarshalRichBlockObject(trimmed)
}

func unmarshalRichTextObject(data []byte) (RichText, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case richTextBoldType:
		text, err := decodeRichTextField(data)
		return RichTextBold{Type: richTextBoldType, Text: text}, err
	case richTextItalicType:
		text, err := decodeRichTextField(data)
		return RichTextItalic{Type: richTextItalicType, Text: text}, err
	case richTextUnderlineType:
		text, err := decodeRichTextField(data)
		return RichTextUnderline{Type: richTextUnderlineType, Text: text}, err
	case richTextStrikethroughType:
		text, err := decodeRichTextField(data)
		return RichTextStrikethrough{Type: richTextStrikethroughType, Text: text}, err
	case richTextSpoilerType:
		text, err := decodeRichTextField(data)
		return RichTextSpoiler{Type: richTextSpoilerType, Text: text}, err
	case richTextDateTimeType:
		var payload struct {
			Text           json.RawMessage `json:"text"`
			UnixTime       int64           `json:"unix_time"`
			DateTimeFormat string          `json:"date_time_format"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextDateTime{Type: richTextDateTimeType, Text: text, UnixTime: payload.UnixTime, DateTimeFormat: payload.DateTimeFormat}, err
	case richTextTextMentionType:
		var payload struct {
			Text json.RawMessage `json:"text"`
			User User            `json:"user"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextTextMention{Type: richTextTextMentionType, Text: text, User: payload.User}, err
	case richTextSubscriptType:
		text, err := decodeRichTextField(data)
		return RichTextSubscript{Type: richTextSubscriptType, Text: text}, err
	case richTextSuperscriptType:
		text, err := decodeRichTextField(data)
		return RichTextSuperscript{Type: richTextSuperscriptType, Text: text}, err
	case richTextMarkedType:
		text, err := decodeRichTextField(data)
		return RichTextMarked{Type: richTextMarkedType, Text: text}, err
	case richTextCodeType:
		text, err := decodeRichTextField(data)
		return RichTextCode{Type: richTextCodeType, Text: text}, err
	case richTextCustomEmojiType:
		var payload RichTextCustomEmoji
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		payload.Type = richTextCustomEmojiType
		return payload, nil
	case richTextMathematicalExpressionType:
		var payload RichTextMathematicalExpression
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		payload.Type = richTextMathematicalExpressionType
		return payload, nil
	case richTextURLType:
		var payload struct {
			Text json.RawMessage `json:"text"`
			URL  string          `json:"url"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextURL{Type: richTextURLType, Text: text, URL: payload.URL}, err
	case richTextEmailAddressType:
		var payload struct {
			Text         json.RawMessage `json:"text"`
			EmailAddress string          `json:"email_address"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextEmailAddress{Type: richTextEmailAddressType, Text: text, EmailAddress: payload.EmailAddress}, err
	case richTextPhoneNumberType:
		var payload struct {
			Text        json.RawMessage `json:"text"`
			PhoneNumber string          `json:"phone_number"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextPhoneNumber{Type: richTextPhoneNumberType, Text: text, PhoneNumber: payload.PhoneNumber}, err
	case richTextBankCardNumberType:
		var payload struct {
			Text           json.RawMessage `json:"text"`
			BankCardNumber string          `json:"bank_card_number"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextBankCardNumber{Type: richTextBankCardNumberType, Text: text, BankCardNumber: payload.BankCardNumber}, err
	case richTextMentionType:
		var payload struct {
			Text     json.RawMessage `json:"text"`
			Username string          `json:"username"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextMention{Type: richTextMentionType, Text: text, Username: payload.Username}, err
	case richTextHashtagType:
		var payload struct {
			Text    json.RawMessage `json:"text"`
			Hashtag string          `json:"hashtag"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextHashtag{Type: richTextHashtagType, Text: text, Hashtag: payload.Hashtag}, err
	case richTextCashtagType:
		var payload struct {
			Text    json.RawMessage `json:"text"`
			Cashtag string          `json:"cashtag"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextCashtag{Type: richTextCashtagType, Text: text, Cashtag: payload.Cashtag}, err
	case richTextBotCommandType:
		var payload struct {
			Text       json.RawMessage `json:"text"`
			BotCommand string          `json:"bot_command"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextBotCommand{Type: richTextBotCommandType, Text: text, BotCommand: payload.BotCommand}, err
	case richTextAnchorType:
		var payload RichTextAnchor
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		payload.Type = richTextAnchorType
		return payload, nil
	case richTextAnchorLinkType:
		var payload struct {
			Text       json.RawMessage `json:"text"`
			AnchorName string          `json:"anchor_name"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextAnchorLink{Type: richTextAnchorLinkType, Text: text, AnchorName: payload.AnchorName}, err
	case richTextReferenceType:
		var payload struct {
			Text json.RawMessage `json:"text"`
			Name string          `json:"name"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextReference{Type: richTextReferenceType, Text: text, Name: payload.Name}, err
	case richTextReferenceLinkType:
		var payload struct {
			Text          json.RawMessage `json:"text"`
			ReferenceName string          `json:"reference_name"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichTextReferenceLink{Type: richTextReferenceLinkType, Text: text, ReferenceName: payload.ReferenceName}, err
	default:
		return nil, stderrors.New("unsupported rich text type")
	}
}

func unmarshalRichBlockObject(data []byte) (RichBlock, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case richBlockParagraphType:
		text, err := decodeRichTextField(data)
		return RichBlockParagraph{Type: richBlockParagraphType, Text: text}, err
	case richBlockSectionHeadingType:
		var payload struct {
			Text json.RawMessage `json:"text"`
			Size int             `json:"size"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichBlockSectionHeading{Type: richBlockSectionHeadingType, Text: text, Size: payload.Size}, err
	case richBlockPreformattedType:
		var payload struct {
			Text     json.RawMessage `json:"text"`
			Language string          `json:"language"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		return RichBlockPreformatted{Type: richBlockPreformattedType, Text: text, Language: payload.Language}, err
	case richBlockFooterType:
		text, err := decodeRichTextField(data)
		return RichBlockFooter{Type: richBlockFooterType, Text: text}, err
	case richBlockDividerType:
		return RichBlockDivider{Type: richBlockDividerType}, nil
	case richBlockMathematicalExpressionType:
		var payload RichBlockMathematicalExpression
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		payload.Type = richBlockMathematicalExpressionType
		return payload, nil
	case richBlockAnchorType:
		var payload RichBlockAnchor
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		payload.Type = richBlockAnchorType
		return payload, nil
	case richBlockListType:
		var payload struct {
			Items []RichBlockListItem `json:"items"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		return RichBlockList{Type: richBlockListType, Items: payload.Items}, nil
	case richBlockBlockQuotationType:
		var payload struct {
			Blocks []json.RawMessage `json:"blocks"`
			Credit json.RawMessage   `json:"credit"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		blocks, err := decodeRichBlockSlice(payload.Blocks)
		if err != nil {
			return nil, err
		}
		credit, err := decodeOptionalRichText(payload.Credit)
		return RichBlockBlockQuotation{Type: richBlockBlockQuotationType, Blocks: blocks, Credit: credit}, err
	case richBlockPullQuotationType:
		var payload struct {
			Text   json.RawMessage `json:"text"`
			Credit json.RawMessage `json:"credit"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		text, err := UnmarshalRichText(payload.Text)
		if err != nil {
			return nil, err
		}
		credit, err := decodeOptionalRichText(payload.Credit)
		return RichBlockPullQuotation{Type: richBlockPullQuotationType, Text: text, Credit: credit}, err
	case richBlockCollageType:
		blocks, caption, err := decodeRichBlockContainer(data)
		return RichBlockCollage{Type: richBlockCollageType, Blocks: blocks, Caption: caption}, err
	case richBlockSlideshowType:
		blocks, caption, err := decodeRichBlockContainer(data)
		return RichBlockSlideshow{Type: richBlockSlideshowType, Blocks: blocks, Caption: caption}, err
	case richBlockTableType:
		var payload struct {
			Cells      [][]RichBlockTableCell `json:"cells"`
			IsBordered bool                   `json:"is_bordered"`
			IsStriped  bool                   `json:"is_striped"`
			Caption    json.RawMessage        `json:"caption"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		caption, err := decodeOptionalRichText(payload.Caption)
		return RichBlockTable{Type: richBlockTableType, Cells: payload.Cells, IsBordered: payload.IsBordered, IsStriped: payload.IsStriped, Caption: caption}, err
	case richBlockDetailsType:
		var payload struct {
			Summary json.RawMessage   `json:"summary"`
			Blocks  []json.RawMessage `json:"blocks"`
			IsOpen  bool              `json:"is_open"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		summary, err := UnmarshalRichText(payload.Summary)
		if err != nil {
			return nil, err
		}
		blocks, err := decodeRichBlockSlice(payload.Blocks)
		return RichBlockDetails{Type: richBlockDetailsType, Summary: summary, Blocks: blocks, IsOpen: payload.IsOpen}, err
	case richBlockMapType:
		var payload struct {
			Location Location        `json:"location"`
			Zoom     int             `json:"zoom"`
			Width    int             `json:"width"`
			Height   int             `json:"height"`
			Caption  json.RawMessage `json:"caption"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		caption, err := decodeOptionalRichBlockCaption(payload.Caption)
		return RichBlockMap{Type: richBlockMapType, Location: payload.Location, Zoom: payload.Zoom, Width: payload.Width, Height: payload.Height, Caption: caption}, err
	case richBlockAnimationType:
		var payload struct {
			Animation  Animation       `json:"animation"`
			HasSpoiler bool            `json:"has_spoiler"`
			Caption    json.RawMessage `json:"caption"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		caption, err := decodeOptionalRichBlockCaption(payload.Caption)
		return RichBlockAnimation{Type: richBlockAnimationType, Animation: payload.Animation, HasSpoiler: payload.HasSpoiler, Caption: caption}, err
	case richBlockAudioType:
		var payload struct {
			Audio   Audio           `json:"audio"`
			Caption json.RawMessage `json:"caption"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		caption, err := decodeOptionalRichBlockCaption(payload.Caption)
		return RichBlockAudio{Type: richBlockAudioType, Audio: payload.Audio, Caption: caption}, err
	case richBlockPhotoType:
		var payload struct {
			Photo      []PhotoSize     `json:"photo"`
			HasSpoiler bool            `json:"has_spoiler"`
			Caption    json.RawMessage `json:"caption"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		caption, err := decodeOptionalRichBlockCaption(payload.Caption)
		return RichBlockPhoto{Type: richBlockPhotoType, Photo: payload.Photo, HasSpoiler: payload.HasSpoiler, Caption: caption}, err
	case richBlockVideoType:
		var payload struct {
			Video      Video           `json:"video"`
			HasSpoiler bool            `json:"has_spoiler"`
			Caption    json.RawMessage `json:"caption"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		caption, err := decodeOptionalRichBlockCaption(payload.Caption)
		return RichBlockVideo{Type: richBlockVideoType, Video: payload.Video, HasSpoiler: payload.HasSpoiler, Caption: caption}, err
	case richBlockVoiceNoteType:
		var payload struct {
			VoiceNote Voice           `json:"voice_note"`
			Caption   json.RawMessage `json:"caption"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, err
		}
		caption, err := decodeOptionalRichBlockCaption(payload.Caption)
		return RichBlockVoiceNote{Type: richBlockVoiceNoteType, VoiceNote: payload.VoiceNote, Caption: caption}, err
	case richBlockThinkingType:
		text, err := decodeRichTextField(data)
		return RichBlockThinking{Type: richBlockThinkingType, Text: text}, err
	default:
		return nil, stderrors.New("unsupported rich block type")
	}
}

func decodeRichTextField(data []byte) (RichText, error) {
	var payload struct {
		Text json.RawMessage `json:"text"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return UnmarshalRichText(payload.Text)
}

func decodeOptionalRichText(raw json.RawMessage) (RichText, error) {
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, nil
	}
	return UnmarshalRichText(raw)
}

func decodeRichBlockSlice(rawItems []json.RawMessage) ([]RichBlock, error) {
	blocks := make([]RichBlock, 0, len(rawItems))
	for _, raw := range rawItems {
		block, err := UnmarshalRichBlock(raw)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func decodeRichBlockContainer(data []byte) ([]RichBlock, *RichBlockCaption, error) {
	var payload struct {
		Blocks  []json.RawMessage `json:"blocks"`
		Caption json.RawMessage   `json:"caption"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, nil, err
	}
	blocks, err := decodeRichBlockSlice(payload.Blocks)
	if err != nil {
		return nil, nil, err
	}
	caption, err := decodeOptionalRichBlockCaption(payload.Caption)
	if err != nil {
		return nil, nil, err
	}
	return blocks, caption, nil
}

func decodeOptionalRichBlockCaption(raw json.RawMessage) (*RichBlockCaption, error) {
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, nil
	}
	var caption RichBlockCaption
	if err := json.Unmarshal(raw, &caption); err != nil {
		return nil, err
	}
	return &caption, nil
}

// UnmarshalJSON decodes RichBlockCaption text fields as polymorphic rich text.
func (caption *RichBlockCaption) UnmarshalJSON(data []byte) error {
	var payload struct {
		Text   json.RawMessage `json:"text"`
		Credit json.RawMessage `json:"credit"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	text, err := UnmarshalRichText(payload.Text)
	if err != nil {
		return err
	}
	credit, err := decodeOptionalRichText(payload.Credit)
	if err != nil {
		return err
	}
	caption.Text = text
	caption.Credit = credit
	return nil
}

// UnmarshalJSON decodes RichBlockTableCell text as polymorphic rich text.
func (cell *RichBlockTableCell) UnmarshalJSON(data []byte) error {
	var payload struct {
		Text     json.RawMessage `json:"text"`
		IsHeader bool            `json:"is_header"`
		Colspan  int             `json:"colspan"`
		Rowspan  int             `json:"rowspan"`
		Align    string          `json:"align"`
		Valign   string          `json:"valign"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	text, err := decodeOptionalRichText(payload.Text)
	if err != nil {
		return err
	}
	cell.Text = text
	cell.IsHeader = payload.IsHeader
	cell.Colspan = payload.Colspan
	cell.Rowspan = payload.Rowspan
	cell.Align = payload.Align
	cell.Valign = payload.Valign
	return nil
}

// UnmarshalJSON decodes RichBlockListItem blocks as polymorphic rich blocks.
func (item *RichBlockListItem) UnmarshalJSON(data []byte) error {
	var payload struct {
		Label       string            `json:"label"`
		Blocks      []json.RawMessage `json:"blocks"`
		HasCheckbox bool              `json:"has_checkbox"`
		IsChecked   bool              `json:"is_checked"`
		Value       int               `json:"value"`
		Type        string            `json:"type"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	blocks, err := decodeRichBlockSlice(payload.Blocks)
	if err != nil {
		return err
	}
	item.Label = payload.Label
	item.Blocks = blocks
	item.HasCheckbox = payload.HasCheckbox
	item.IsChecked = payload.IsChecked
	item.Value = payload.Value
	item.Type = payload.Type
	return nil
}
