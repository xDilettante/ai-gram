package telegram

import (
	"encoding/json"
	"testing"
)

func TestRichMessageDecodesTextAndBlockVariants(t *testing.T) {
	var message RichMessage
	if err := json.Unmarshal([]byte(`{
		"blocks": [
			{
				"type": "heading",
				"text": ["Hello ", {"type": "bold", "text": "world"}],
				"size": 2
			},
			{
				"type": "table",
				"cells": [[
					{"text": "Metric", "is_header": true, "align": "center", "valign": "top"},
					{"text": {"type": "code", "text": "42"}}
				]],
				"is_bordered": true,
				"caption": "Results"
			}
		],
		"is_rtl": true
	}`), &message); err != nil {
		t.Fatalf("decode rich message: %v", err)
	}

	if !message.IsRTL || len(message.Blocks) != 2 {
		t.Fatalf("unexpected rich message: %+v", message)
	}

	heading, ok := message.Blocks[0].(RichBlockSectionHeading)
	if !ok || heading.Type != "heading" || heading.Size != 2 {
		t.Fatalf("unexpected heading block: %#v", message.Blocks[0])
	}
	parts, ok := heading.Text.(RichTextList)
	if !ok || len(parts) != 2 {
		t.Fatalf("unexpected heading text: %#v", heading.Text)
	}
	if plain, ok := parts[0].(RichTextPlain); !ok || string(plain) != "Hello " {
		t.Fatalf("unexpected plain heading part: %#v", parts[0])
	}
	if bold, ok := parts[1].(RichTextBold); !ok {
		t.Fatalf("unexpected bold heading part: %#v", parts[1])
	} else if text, ok := bold.Text.(RichTextPlain); !ok || string(text) != "world" {
		t.Fatalf("unexpected bold text: %#v", bold.Text)
	}

	table, ok := message.Blocks[1].(RichBlockTable)
	if !ok || table.Type != "table" || !table.IsBordered || len(table.Cells) != 1 || len(table.Cells[0]) != 2 {
		t.Fatalf("unexpected table block: %#v", message.Blocks[1])
	}
	if firstCellText, ok := table.Cells[0][0].Text.(RichTextPlain); !ok || string(firstCellText) != "Metric" || !table.Cells[0][0].IsHeader || table.Cells[0][0].Align != "center" || table.Cells[0][0].Valign != "top" {
		t.Fatalf("unexpected first table cell: %#v", table.Cells[0][0])
	}
	if caption, ok := table.Caption.(RichTextPlain); !ok || string(caption) != "Results" {
		t.Fatalf("unexpected table caption: %#v", table.Caption)
	}
}

func TestRichMessageDecodesAllBotAPI101Variants(t *testing.T) {
	textPayloads := map[string]string{
		"italic":                  `{"type":"italic","text":"x"}`,
		"underline":               `{"type":"underline","text":"x"}`,
		"strikethrough":           `{"type":"strikethrough","text":"x"}`,
		"spoiler":                 `{"type":"spoiler","text":"x"}`,
		"date_time":               `{"type":"date_time","text":"x","unix_time":123,"date_time_format":"wDT"}`,
		"text_mention":            `{"type":"text_mention","text":"x","user":{"id":1,"is_bot":false,"first_name":"Alice"}}`,
		"subscript":               `{"type":"subscript","text":"x"}`,
		"superscript":             `{"type":"superscript","text":"x"}`,
		"marked":                  `{"type":"marked","text":"x"}`,
		"custom_emoji":            `{"type":"custom_emoji","custom_emoji_id":"emoji-id","alternative_text":"👍"}`,
		"mathematical_expression": `{"type":"mathematical_expression","expression":"x^2"}`,
		"url":                     `{"type":"url","text":"x","url":"https://example.com"}`,
		"email_address":           `{"type":"email_address","text":"x","email_address":"user@example.com"}`,
		"phone_number":            `{"type":"phone_number","text":"x","phone_number":"+10000000000"}`,
		"bank_card_number":        `{"type":"bank_card_number","text":"x","bank_card_number":"4242424242424242"}`,
		"mention":                 `{"type":"mention","text":"x","username":"alice"}`,
		"hashtag":                 `{"type":"hashtag","text":"x","hashtag":"tag"}`,
		"cashtag":                 `{"type":"cashtag","text":"x","cashtag":"USD"}`,
		"bot_command":             `{"type":"bot_command","text":"x","bot_command":"/start"}`,
		"anchor":                  `{"type":"anchor","name":"chapter-1"}`,
		"anchor_link":             `{"type":"anchor_link","text":"x","anchor_name":"chapter-1"}`,
		"reference":               `{"type":"reference","text":"x","name":"note-1"}`,
		"reference_link":          `{"type":"reference_link","text":"x","reference_name":"note-1"}`,
	}
	for name, payload := range textPayloads {
		t.Run("text "+name, func(t *testing.T) {
			text, err := UnmarshalRichText([]byte(payload))
			if err != nil {
				t.Fatalf("decode rich text: %v", err)
			}
			if text == nil {
				t.Fatal("expected rich text")
			}
		})
	}

	blockPayloads := map[string]string{
		"paragraph":               `{"type":"paragraph","text":"x"}`,
		"pre":                     `{"type":"pre","text":"x","language":"go"}`,
		"footer":                  `{"type":"footer","text":"x"}`,
		"divider":                 `{"type":"divider"}`,
		"mathematical_expression": `{"type":"mathematical_expression","expression":"x^2"}`,
		"anchor":                  `{"type":"anchor","name":"chapter-1"}`,
		"list":                    `{"type":"list","items":[{"label":"-","blocks":[{"type":"paragraph","text":"x"}],"has_checkbox":true,"is_checked":true,"value":1,"type":"1"}]}`,
		"blockquote":              `{"type":"blockquote","blocks":[{"type":"paragraph","text":"x"}],"credit":"Alice"}`,
		"pullquote":               `{"type":"pullquote","text":"x","credit":"Alice"}`,
		"collage":                 `{"type":"collage","blocks":[{"type":"photo","photo":[{"file_id":"p","file_unique_id":"u","width":1,"height":1}]}],"caption":{"text":"x"}}`,
		"slideshow":               `{"type":"slideshow","blocks":[{"type":"video","video":{"file_id":"v","file_unique_id":"u","width":1,"height":1,"duration":1}}],"caption":{"text":"x","credit":"Alice"}}`,
		"details":                 `{"type":"details","summary":"More","blocks":[{"type":"paragraph","text":"x"}],"is_open":true}`,
		"map":                     `{"type":"map","location":{"latitude":41.9,"longitude":12.5},"zoom":14,"width":320,"height":200,"caption":{"text":"Map"}}`,
		"animation":               `{"type":"animation","animation":{"file_id":"a","file_unique_id":"u","width":1,"height":1,"duration":1},"has_spoiler":true,"caption":{"text":"Animation"}}`,
		"audio":                   `{"type":"audio","audio":{"file_id":"a","file_unique_id":"u","duration":1},"caption":{"text":"Audio"}}`,
		"photo":                   `{"type":"photo","photo":[{"file_id":"p","file_unique_id":"u","width":1,"height":1}],"has_spoiler":true,"caption":{"text":"Photo"}}`,
		"video":                   `{"type":"video","video":{"file_id":"v","file_unique_id":"u","width":1,"height":1,"duration":1},"has_spoiler":true,"caption":{"text":"Video"}}`,
		"voice_note":              `{"type":"voice_note","voice_note":{"file_id":"v","file_unique_id":"u","duration":1},"caption":{"text":"Voice"}}`,
		"thinking":                `{"type":"thinking","text":"Thinking"}`,
	}
	for name, payload := range blockPayloads {
		t.Run("block "+name, func(t *testing.T) {
			block, err := UnmarshalRichBlock([]byte(payload))
			if err != nil {
				t.Fatalf("decode rich block: %v", err)
			}
			if block == nil {
				t.Fatal("expected rich block")
			}
		})
	}
}

func TestMessageDecodesRichMessage(t *testing.T) {
	var message Message
	if err := json.Unmarshal([]byte(`{
		"message_id": 1,
		"date": 123,
		"chat": {"id": 10, "type": "private"},
		"rich_message": {
			"blocks": [{"type": "paragraph", "text": "Hello"}]
		}
	}`), &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if message.RichMessage == nil || len(message.RichMessage.Blocks) != 1 {
		t.Fatalf("unexpected rich message: %+v", message.RichMessage)
	}
	if paragraph, ok := message.RichMessage.Blocks[0].(RichBlockParagraph); !ok {
		t.Fatalf("unexpected block: %#v", message.RichMessage.Blocks[0])
	} else if text, ok := paragraph.Text.(RichTextPlain); !ok || string(text) != "Hello" {
		t.Fatalf("unexpected paragraph text: %#v", paragraph.Text)
	}
}
