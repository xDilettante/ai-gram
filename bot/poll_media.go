package bot

import (
	stderrors "errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// InputPollMedia describes media accepted in sendPoll media fields.
type InputPollMedia interface {
	inputPollMedia()
}

// InputPollOptionMedia describes media accepted in InputPollOption.media.
type InputPollOptionMedia interface {
	inputPollOptionMedia()
}

// InputMediaLocation describes a location input media item.
type InputMediaLocation struct {
	Type               string  `json:"type"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	HorizontalAccuracy float64 `json:"horizontal_accuracy,omitempty"`
}

// InputMediaSticker describes a sticker input media item.
type InputMediaSticker struct {
	Type  string  `json:"type"`
	Media FileRef `json:"media"`
	Emoji string  `json:"emoji,omitempty"`
}

// InputMediaVenue describes a venue input media item.
type InputMediaVenue struct {
	Type            string  `json:"type"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FoursquareID    string  `json:"foursquare_id,omitempty"`
	FoursquareType  string  `json:"foursquare_type,omitempty"`
	GooglePlaceID   string  `json:"google_place_id,omitempty"`
	GooglePlaceType string  `json:"google_place_type,omitempty"`
}

func (InputMediaPhoto) inputPollMedia()     {}
func (InputMediaVideo) inputPollMedia()     {}
func (InputMediaAnimation) inputPollMedia() {}
func (InputMediaAudio) inputPollMedia()     {}
func (InputMediaDocument) inputPollMedia()  {}
func (InputMediaLivePhoto) inputPollMedia() {}
func (InputMediaLocation) inputPollMedia()  {}
func (InputMediaVenue) inputPollMedia()     {}

func (InputMediaPhoto) inputPollOptionMedia()     {}
func (InputMediaVideo) inputPollOptionMedia()     {}
func (InputMediaAnimation) inputPollOptionMedia() {}
func (InputMediaLivePhoto) inputPollOptionMedia() {}
func (InputMediaLocation) inputPollOptionMedia()  {}
func (InputMediaSticker) inputPollOptionMedia()   {}
func (InputMediaVenue) inputPollOptionMedia()     {}

// MediaLocation creates a location input media item.
func MediaLocation(latitude float64, longitude float64) InputMediaLocation {
	return InputMediaLocation{Type: "location", Latitude: latitude, Longitude: longitude}
}

// MediaSticker creates a sticker input media item.
func MediaSticker(media FileRef) InputMediaSticker {
	return InputMediaSticker{Type: "sticker", Media: media}
}

// MediaVenue creates a venue input media item.
func MediaVenue(latitude float64, longitude float64, title string, address string) InputMediaVenue {
	return InputMediaVenue{Type: "venue", Latitude: latitude, Longitude: longitude, Title: title, Address: address}
}

type inputPollMediaPayload struct {
	Type                        string                   `json:"type"`
	Media                       string                   `json:"media,omitempty"`
	Photo                       string                   `json:"photo,omitempty"`
	Thumbnail                   string                   `json:"thumbnail,omitempty"`
	Cover                       string                   `json:"cover,omitempty"`
	StartTimestamp              int                      `json:"start_timestamp,omitempty"`
	Caption                     string                   `json:"caption,omitempty"`
	ParseMode                   string                   `json:"parse_mode,omitempty"`
	CaptionEntities             []telegram.MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia       bool                     `json:"show_caption_above_media,omitempty"`
	HasSpoiler                  bool                     `json:"has_spoiler,omitempty"`
	Width                       int                      `json:"width,omitempty"`
	Height                      int                      `json:"height,omitempty"`
	Duration                    int                      `json:"duration,omitempty"`
	SupportsStreaming           bool                     `json:"supports_streaming,omitempty"`
	Performer                   string                   `json:"performer,omitempty"`
	Title                       string                   `json:"title,omitempty"`
	DisableContentTypeDetection bool                     `json:"disable_content_type_detection,omitempty"`
	Latitude                    float64                  `json:"latitude,omitempty"`
	Longitude                   float64                  `json:"longitude,omitempty"`
	HorizontalAccuracy          float64                  `json:"horizontal_accuracy,omitempty"`
	Address                     string                   `json:"address,omitempty"`
	FoursquareID                string                   `json:"foursquare_id,omitempty"`
	FoursquareType              string                   `json:"foursquare_type,omitempty"`
	GooglePlaceID               string                   `json:"google_place_id,omitempty"`
	GooglePlaceType             string                   `json:"google_place_type,omitempty"`
	Emoji                       string                   `json:"emoji,omitempty"`
}

type inputPollOptionPayload struct {
	Text          string                   `json:"text"`
	TextParseMode string                   `json:"text_parse_mode,omitempty"`
	TextEntities  []telegram.MessageEntity `json:"text_entities,omitempty"`
	Media         any                      `json:"media,omitempty"`
}

func buildInputPollOptionPayload(option telegram.InputPollOption) (inputPollOptionPayload, error) {
	return buildInputPollOptionPayloadForSend(option, "options", nil)
}

func buildInputPollOptionPayloadForSend(option telegram.InputPollOption, field string, files map[string]UploadFile) (inputPollOptionPayload, error) {
	payload := inputPollOptionPayload{
		Text:          option.Text,
		TextParseMode: option.TextParseMode,
		TextEntities:  option.TextEntities,
	}
	if option.Media != nil && !isNilPollMediaInterface(option.Media) {
		media, err := buildInputPollOptionMediaPayload(option.Media, field+".media", files)
		if err != nil {
			return inputPollOptionPayload{}, err
		}
		payload.Media = media
	}
	return payload, nil
}

func buildInputPollMediaPayload(media InputPollMedia, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	switch item := media.(type) {
	case nil:
		return inputPollMediaPayload{}, nil
	case InputMediaPhoto:
		return buildPollMediaPhotoPayload(item, field, files)
	case *InputMediaPhoto:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaPhotoPayload(*item, field, files)
	case InputMediaVideo:
		return buildPollMediaVideoPayload(item, field, files)
	case *InputMediaVideo:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaVideoPayload(*item, field, files)
	case InputMediaAnimation:
		return buildPollMediaAnimationPayload(item, field, files)
	case *InputMediaAnimation:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaAnimationPayload(*item, field, files)
	case InputMediaAudio:
		return buildPollMediaAudioPayload(item, field, files)
	case *InputMediaAudio:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaAudioPayload(*item, field, files)
	case InputMediaDocument:
		return buildPollMediaDocumentPayload(item, field, files)
	case *InputMediaDocument:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaDocumentPayload(*item, field, files)
	case InputMediaLivePhoto:
		return buildPollMediaLivePhotoPayload(item, field, files)
	case *InputMediaLivePhoto:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaLivePhotoPayload(*item, field, files)
	case InputMediaLocation:
		return buildPollMediaLocationPayload(item, field)
	case *InputMediaLocation:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaLocationPayload(*item, field)
	case InputMediaVenue:
		return buildPollMediaVenuePayload(item, field)
	case *InputMediaVenue:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaVenuePayload(*item, field)
	default:
		return inputPollMediaPayload{}, stderrors.New("unsupported input poll media type")
	}
}

func buildInputPollOptionMediaPayload(media any, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	switch item := media.(type) {
	case nil:
		return inputPollMediaPayload{}, nil
	case InputMediaPhoto:
		return buildPollMediaPhotoPayload(item, field, files)
	case *InputMediaPhoto:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaPhotoPayload(*item, field, files)
	case InputMediaVideo:
		return buildPollMediaVideoPayload(item, field, files)
	case *InputMediaVideo:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaVideoPayload(*item, field, files)
	case InputMediaAnimation:
		return buildPollMediaAnimationPayload(item, field, files)
	case *InputMediaAnimation:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaAnimationPayload(*item, field, files)
	case InputMediaLivePhoto:
		return buildPollMediaLivePhotoPayload(item, field, files)
	case *InputMediaLivePhoto:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaLivePhotoPayload(*item, field, files)
	case InputMediaLocation:
		return buildPollMediaLocationPayload(item, field)
	case *InputMediaLocation:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaLocationPayload(*item, field)
	case InputMediaSticker:
		return buildPollMediaStickerPayload(item, field, files)
	case *InputMediaSticker:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaStickerPayload(*item, field, files)
	case InputMediaVenue:
		return buildPollMediaVenuePayload(item, field)
	case *InputMediaVenue:
		if item == nil {
			return inputPollMediaPayload{}, nil
		}
		return buildPollMediaVenuePayload(*item, field)
	case InputMediaAudio, *InputMediaAudio, InputMediaDocument, *InputMediaDocument:
		return inputPollMediaPayload{}, stderrors.New("audio and document media are not supported in poll options")
	default:
		return inputPollMediaPayload{}, stderrors.New("unsupported input poll option media type")
	}
}

func validateInputPollMedia(media InputPollMedia, field string) error {
	if media == nil || isNilPollMediaInterface(media) {
		return nil
	}
	_, err := buildInputPollMediaPayload(media, field, make(map[string]UploadFile))
	return err
}

func validateInputPollOptionMedia(media any, field string) error {
	if media == nil || isNilPollMediaInterface(media) {
		return nil
	}
	_, err := buildInputPollOptionMediaPayload(media, field, make(map[string]UploadFile))
	return err
}

func buildPollMediaPhotoPayload(media InputMediaPhoto, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	if err := validateInputMediaType(media.Type, "photo"); err != nil {
		return inputPollMediaPayload{}, err
	}
	mediaValue, err := pollMediaFileValue(media.Media, field+".media", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	return inputPollMediaPayload{
		Type:                  mediaType(media.Type, "photo"),
		Media:                 mediaValue,
		Caption:               media.Caption,
		ParseMode:             media.ParseMode,
		CaptionEntities:       media.CaptionEntities,
		ShowCaptionAboveMedia: media.ShowCaptionAboveMedia,
		HasSpoiler:            media.HasSpoiler,
	}, validateCaptionFormatting(media.ParseMode, media.CaptionEntities)
}

func buildPollMediaVideoPayload(media InputMediaVideo, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	if err := validateInputMediaVideo(media); err != nil {
		return inputPollMediaPayload{}, err
	}
	mediaValue, err := pollMediaFileValue(media.Media, field+".media", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	cover, err := optionalPollMediaFileValue(media.Cover, field+".cover", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	thumbnail, err := optionalPollMediaFileValue(media.Thumbnail, field+".thumbnail", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	return inputPollMediaPayload{
		Type:                  mediaType(media.Type, "video"),
		Media:                 mediaValue,
		Thumbnail:             thumbnail,
		Cover:                 cover,
		StartTimestamp:        media.StartTimestamp,
		Caption:               media.Caption,
		ParseMode:             media.ParseMode,
		CaptionEntities:       media.CaptionEntities,
		ShowCaptionAboveMedia: media.ShowCaptionAboveMedia,
		Width:                 media.Width,
		Height:                media.Height,
		Duration:              media.Duration,
		SupportsStreaming:     media.SupportsStreaming,
		HasSpoiler:            media.HasSpoiler,
	}, nil
}

func buildPollMediaAnimationPayload(media InputMediaAnimation, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	if err := validateInputMediaAnimation(media); err != nil {
		return inputPollMediaPayload{}, err
	}
	mediaValue, err := pollMediaFileValue(media.Media, field+".media", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	thumbnail, err := optionalPollMediaFileValue(media.Thumbnail, field+".thumbnail", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	return inputPollMediaPayload{
		Type:                  mediaType(media.Type, "animation"),
		Media:                 mediaValue,
		Thumbnail:             thumbnail,
		Caption:               media.Caption,
		ParseMode:             media.ParseMode,
		CaptionEntities:       media.CaptionEntities,
		ShowCaptionAboveMedia: media.ShowCaptionAboveMedia,
		Width:                 media.Width,
		Height:                media.Height,
		Duration:              media.Duration,
		HasSpoiler:            media.HasSpoiler,
	}, nil
}

func buildPollMediaAudioPayload(media InputMediaAudio, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	if err := validateInputMediaAudio(media); err != nil {
		return inputPollMediaPayload{}, err
	}
	mediaValue, err := pollMediaFileValue(media.Media, field+".media", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	thumbnail, err := optionalPollMediaFileValue(media.Thumbnail, field+".thumbnail", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	return inputPollMediaPayload{
		Type:            mediaType(media.Type, "audio"),
		Media:           mediaValue,
		Thumbnail:       thumbnail,
		Caption:         media.Caption,
		ParseMode:       media.ParseMode,
		CaptionEntities: media.CaptionEntities,
		Duration:        media.Duration,
		Performer:       media.Performer,
		Title:           media.Title,
	}, nil
}

func buildPollMediaDocumentPayload(media InputMediaDocument, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	if err := validateInputMediaDocument(media); err != nil {
		return inputPollMediaPayload{}, err
	}
	mediaValue, err := pollMediaFileValue(media.Media, field+".media", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	thumbnail, err := optionalPollMediaFileValue(media.Thumbnail, field+".thumbnail", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	return inputPollMediaPayload{
		Type:                        mediaType(media.Type, "document"),
		Media:                       mediaValue,
		Thumbnail:                   thumbnail,
		Caption:                     media.Caption,
		ParseMode:                   media.ParseMode,
		CaptionEntities:             media.CaptionEntities,
		DisableContentTypeDetection: media.DisableContentTypeDetection,
	}, nil
}

func buildPollMediaLivePhotoPayload(media InputMediaLivePhoto, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	if err := validateInputMediaType(media.Type, "live_photo"); err != nil {
		return inputPollMediaPayload{}, err
	}
	mediaValue, err := pollMediaFileValue(media.Media, field+".media", false, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	photoValue, err := pollMediaFileValue(media.Photo, field+".photo", false, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	if err := validateCaptionFormatting(media.ParseMode, media.CaptionEntities); err != nil {
		return inputPollMediaPayload{}, err
	}
	return inputPollMediaPayload{
		Type:                  mediaType(media.Type, "live_photo"),
		Media:                 mediaValue,
		Photo:                 photoValue,
		Caption:               media.Caption,
		ParseMode:             media.ParseMode,
		CaptionEntities:       media.CaptionEntities,
		ShowCaptionAboveMedia: media.ShowCaptionAboveMedia,
		HasSpoiler:            media.HasSpoiler,
	}, nil
}

func buildPollMediaLocationPayload(media InputMediaLocation, field string) (inputPollMediaPayload, error) {
	if err := validateInputMediaType(media.Type, "location"); err != nil {
		return inputPollMediaPayload{}, err
	}
	if err := validateLatitude(media.Latitude); err != nil {
		return inputPollMediaPayload{}, err
	}
	if err := validateLongitude(media.Longitude); err != nil {
		return inputPollMediaPayload{}, err
	}
	if media.HorizontalAccuracy < 0 {
		return inputPollMediaPayload{}, stderrors.New("horizontal_accuracy must not be negative")
	}
	return inputPollMediaPayload{
		Type:               mediaType(media.Type, "location"),
		Latitude:           media.Latitude,
		Longitude:          media.Longitude,
		HorizontalAccuracy: media.HorizontalAccuracy,
	}, nil
}

func buildPollMediaStickerPayload(media InputMediaSticker, field string, files map[string]UploadFile) (inputPollMediaPayload, error) {
	if err := validateInputMediaType(media.Type, "sticker"); err != nil {
		return inputPollMediaPayload{}, err
	}
	mediaValue, err := pollMediaFileValue(media.Media, field+".media", true, files)
	if err != nil {
		return inputPollMediaPayload{}, err
	}
	return inputPollMediaPayload{
		Type:  mediaType(media.Type, "sticker"),
		Media: mediaValue,
		Emoji: media.Emoji,
	}, nil
}

func buildPollMediaVenuePayload(media InputMediaVenue, field string) (inputPollMediaPayload, error) {
	if err := validateInputMediaType(media.Type, "venue"); err != nil {
		return inputPollMediaPayload{}, err
	}
	if err := validateLatitude(media.Latitude); err != nil {
		return inputPollMediaPayload{}, err
	}
	if err := validateLongitude(media.Longitude); err != nil {
		return inputPollMediaPayload{}, err
	}
	if strings.TrimSpace(media.Title) == "" {
		return inputPollMediaPayload{}, stderrors.New("title is required")
	}
	if strings.TrimSpace(media.Address) == "" {
		return inputPollMediaPayload{}, stderrors.New("address is required")
	}
	return inputPollMediaPayload{
		Type:            mediaType(media.Type, "venue"),
		Latitude:        media.Latitude,
		Longitude:       media.Longitude,
		Title:           media.Title,
		Address:         media.Address,
		FoursquareID:    media.FoursquareID,
		FoursquareType:  media.FoursquareType,
		GooglePlaceID:   media.GooglePlaceID,
		GooglePlaceType: media.GooglePlaceType,
	}, nil
}

func pollMediaFileValue(ref FileRef, field string, allowURL bool, files map[string]UploadFile) (string, error) {
	if err := ref.validate(field); err != nil {
		return "", err
	}
	if !allowURL && (ref.kind == fileRefURL || strings.Contains(ref.value, "://")) {
		return "", fmt.Errorf("%s must be a file_id because URLs are not supported for this media type", field)
	}
	if ref.isUpload() {
		if files == nil {
			return "", fmt.Errorf("%s FileUpload requires multipart sendPoll", field)
		}
		name := pollMediaAttachName(field)
		files[name] = ref.upload
		return "attach://" + name, nil
	}
	return ref.value, nil
}

func optionalPollMediaFileValue(ref FileRef, field string, allowURL bool, files map[string]UploadFile) (string, error) {
	if !ref.isSet() {
		return "", nil
	}
	return pollMediaFileValue(ref, field, allowURL, files)
}

func pollMediaAttachName(field string) string {
	var builder strings.Builder
	for _, r := range field {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			builder.WriteRune(r)
			continue
		}
		if builder.Len() > 0 && !strings.HasSuffix(builder.String(), "_") {
			builder.WriteByte('_')
		}
	}
	name := strings.Trim(builder.String(), "_")
	if name == "" {
		return "poll_media"
	}
	return name
}

func isNilPollMediaInterface(value any) bool {
	if value == nil {
		return true
	}
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflectValue.IsNil()
	default:
		return false
	}
}
