package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
)

const (
	paidMediaPreviewType = "preview"
	paidMediaPhotoType   = "photo"
	paidMediaVideoType   = "video"
)

// PaidMediaInfo describes paid media attached to a message.
type PaidMediaInfo struct {
	StarCount int         `json:"star_count"`
	PaidMedia []PaidMedia `json:"paid_media"`
}

// PaidMedia marks Telegram paid media objects.
type PaidMedia interface {
	paidMedia()
}

// PaidMediaPreview describes paid media metadata visible before purchase.
type PaidMediaPreview struct {
	Type     string `json:"type"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Duration int    `json:"duration,omitempty"`
}

// PaidMediaPhoto describes a purchased paid photo.
type PaidMediaPhoto struct {
	Type  string      `json:"type"`
	Photo []PhotoSize `json:"photo"`
}

// PaidMediaVideo describes a purchased paid video.
type PaidMediaVideo struct {
	Type  string `json:"type"`
	Video Video  `json:"video"`
}

// PaidMediaPurchased represents an update sent when paid media with payload is purchased.
type PaidMediaPurchased struct {
	From             User   `json:"from"`
	PaidMediaPayload string `json:"paid_media_payload"`
}

func (PaidMediaPreview) paidMedia() {}
func (PaidMediaPhoto) paidMedia()   {}
func (PaidMediaVideo) paidMedia()   {}

// UnmarshalPaidMedia decodes a polymorphic Telegram PaidMedia object.
func UnmarshalPaidMedia(data []byte) (PaidMedia, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case paidMediaPreviewType:
		var media PaidMediaPreview
		if err := json.Unmarshal(data, &media); err != nil {
			return nil, err
		}
		media.Type = paidMediaPreviewType
		return media, nil
	case paidMediaPhotoType:
		var media PaidMediaPhoto
		if err := json.Unmarshal(data, &media); err != nil {
			return nil, err
		}
		media.Type = paidMediaPhotoType
		return media, nil
	case paidMediaVideoType:
		var media PaidMediaVideo
		if err := json.Unmarshal(data, &media); err != nil {
			return nil, err
		}
		media.Type = paidMediaVideoType
		return media, nil
	default:
		return nil, stderrors.New("unsupported paid media type")
	}
}

// UnmarshalJSON decodes paid media info with polymorphic media items.
func (info *PaidMediaInfo) UnmarshalJSON(data []byte) error {
	var payload struct {
		StarCount int               `json:"star_count"`
		PaidMedia []json.RawMessage `json:"paid_media"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	media, err := unmarshalPaidMediaItems(payload.PaidMedia)
	if err != nil {
		return err
	}
	info.StarCount = payload.StarCount
	info.PaidMedia = media
	return nil
}

func unmarshalPaidMediaItems(rawItems []json.RawMessage) ([]PaidMedia, error) {
	if rawItems == nil {
		return nil, nil
	}
	media := make([]PaidMedia, 0, len(rawItems))
	for index, raw := range rawItems {
		if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
			return nil, fmt.Errorf("paid_media[%d] is required", index)
		}
		item, err := UnmarshalPaidMedia(raw)
		if err != nil {
			return nil, fmt.Errorf("paid_media[%d]: %w", index, err)
		}
		media = append(media, item)
	}
	return media, nil
}
