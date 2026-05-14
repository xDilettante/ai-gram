package callback

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xDilettante/ai-gram/telegram"
)

const (
	// MaxDataBytes is Telegram's callback_data size limit.
	MaxDataBytes = 64

	// ActionConfirm is the conventional action name for confirmation callbacks.
	ActionConfirm = "confirm"
	// ActionCancel is the conventional action name for cancellation callbacks.
	ActionCancel = "cancel"
)

var (
	// ErrEmptyData reports an empty callback payload.
	ErrEmptyData = errors.New("callback data is empty")
	// ErrTooLong reports a callback payload that exceeds Telegram's 64-byte limit.
	ErrTooLong = errors.New("callback data exceeds 64 bytes")
	// ErrInvalidFormat reports malformed callback data.
	ErrInvalidFormat = errors.New("callback data format is invalid")
	// ErrInvalidSegment reports an unsupported character in a callback segment.
	ErrInvalidSegment = errors.New("callback data segment is invalid")
	// ErrInvalidPage reports an invalid pagination value.
	ErrInvalidPage = errors.New("callback page is invalid")
	// ErrInvalidExpiry reports an invalid expiry value.
	ErrInvalidExpiry = errors.New("callback expiry is invalid")
)

// Data is a typed callback payload.
//
// Encoded values use a compact ASCII format:
//
//	namespace:action[:i=id][:p=page][:e=unix_seconds]
//
// Namespace and action are required. ID, page, and expiry are optional.
type Data struct {
	Namespace string
	Action    string
	ID        string
	Page      int
	HasPage   bool
	ExpiresAt time.Time
}

// New creates callback data with the required namespace and action fields.
func New(namespace string, action string) Data {
	return Data{Namespace: namespace, Action: action}
}

// Confirm creates callback data for a conventional confirmation action.
func Confirm(namespace string, id string) Data {
	return New(namespace, ActionConfirm).WithID(id)
}

// Cancel creates callback data for a conventional cancellation action.
func Cancel(namespace string, id string) Data {
	return New(namespace, ActionCancel).WithID(id)
}

// ForPage creates callback data for a paginated action.
func ForPage(namespace string, action string, page int) Data {
	return New(namespace, action).WithPage(page)
}

// WithID returns a copy with an ID field.
func (d Data) WithID(id string) Data {
	d.ID = id
	return d
}

// WithPage returns a copy with a page field.
func (d Data) WithPage(page int) Data {
	d.Page = page
	d.HasPage = true
	return d
}

// WithExpiry returns a copy with an absolute expiry time.
func (d Data) WithExpiry(expiresAt time.Time) Data {
	d.ExpiresAt = expiresAt
	return d
}

// WithTTL returns a copy that expires after ttl relative to now.
func (d Data) WithTTL(now time.Time, ttl time.Duration) Data {
	if ttl <= 0 {
		d.ExpiresAt = now
		return d
	}
	d.ExpiresAt = now.Add(ttl)
	return d
}

// NextPage returns a copy with the page incremented by one.
func (d Data) NextPage() Data {
	if !d.HasPage {
		return d.WithPage(1)
	}
	return d.WithPage(d.Page + 1)
}

// PreviousPage returns a copy with the page decremented by one, clamped to zero.
func (d Data) PreviousPage() Data {
	if !d.HasPage || d.Page <= 0 {
		return d.WithPage(0)
	}
	return d.WithPage(d.Page - 1)
}

// Match reports whether namespace and action match. Empty arguments act as wildcards.
func (d Data) Match(namespace string, action string) bool {
	if namespace != "" && d.Namespace != namespace {
		return false
	}
	if action != "" && d.Action != action {
		return false
	}
	return true
}

// Expired reports whether the callback data is expired at now.
func (d Data) Expired(now time.Time) bool {
	return !d.ExpiresAt.IsZero() && !now.Before(d.ExpiresAt)
}

// Encode returns a Telegram callback_data string.
func (d Data) Encode() (string, error) {
	if err := validateRequired("namespace", d.Namespace); err != nil {
		return "", err
	}
	if err := validateRequired("action", d.Action); err != nil {
		return "", err
	}

	parts := []string{d.Namespace, d.Action}
	if d.ID != "" {
		if err := validateSegment("id", d.ID); err != nil {
			return "", err
		}
		parts = append(parts, "i="+d.ID)
	}
	if d.HasPage {
		if d.Page < 0 {
			return "", fmt.Errorf("%w: %d", ErrInvalidPage, d.Page)
		}
		parts = append(parts, "p="+strconv.Itoa(d.Page))
	}
	if !d.ExpiresAt.IsZero() {
		unix := d.ExpiresAt.Unix()
		if unix <= 0 {
			return "", fmt.Errorf("%w: %d", ErrInvalidExpiry, unix)
		}
		parts = append(parts, "e="+strconv.FormatInt(unix, 10))
	}

	encoded := strings.Join(parts, ":")
	if len(encoded) > MaxDataBytes {
		return "", fmt.Errorf("%w: got %d bytes", ErrTooLong, len(encoded))
	}
	return encoded, nil
}

// Must returns encoded callback data or panics.
func Must(data Data) string {
	encoded, err := data.Encode()
	if err != nil {
		panic(err)
	}
	return encoded
}

// Button creates an inline keyboard callback button from typed data.
func Button(text string, data Data) (telegram.InlineKeyboardButton, error) {
	encoded, err := data.Encode()
	if err != nil {
		return telegram.InlineKeyboardButton{}, err
	}
	return telegram.InlineButtonCallback(text, encoded), nil
}

// MustButton creates an inline keyboard callback button or panics.
func MustButton(text string, data Data) telegram.InlineKeyboardButton {
	return telegram.InlineButtonCallback(text, Must(data))
}

// Parse parses a Telegram callback_data string.
func Parse(raw string) (Data, error) {
	if raw == "" {
		return Data{}, ErrEmptyData
	}
	if len(raw) > MaxDataBytes {
		return Data{}, fmt.Errorf("%w: got %d bytes", ErrTooLong, len(raw))
	}

	parts := strings.Split(raw, ":")
	if len(parts) < 2 {
		return Data{}, ErrInvalidFormat
	}

	data := Data{Namespace: parts[0], Action: parts[1]}
	if err := validateRequired("namespace", data.Namespace); err != nil {
		return Data{}, err
	}
	if err := validateRequired("action", data.Action); err != nil {
		return Data{}, err
	}

	var seenID, seenPage, seenExpiry bool
	for _, part := range parts[2:] {
		key, value, ok := strings.Cut(part, "=")
		if !ok || key == "" || value == "" {
			return Data{}, ErrInvalidFormat
		}

		switch key {
		case "i":
			if seenID {
				return Data{}, ErrInvalidFormat
			}
			if err := validateSegment("id", value); err != nil {
				return Data{}, err
			}
			data.ID = value
			seenID = true
		case "p":
			if seenPage {
				return Data{}, ErrInvalidFormat
			}
			page, err := strconv.Atoi(value)
			if err != nil || page < 0 {
				return Data{}, fmt.Errorf("%w: %q", ErrInvalidPage, value)
			}
			data.Page = page
			data.HasPage = true
			seenPage = true
		case "e":
			if seenExpiry {
				return Data{}, ErrInvalidFormat
			}
			unix, err := strconv.ParseInt(value, 10, 64)
			if err != nil || unix <= 0 {
				return Data{}, fmt.Errorf("%w: %q", ErrInvalidExpiry, value)
			}
			data.ExpiresAt = time.Unix(unix, 0)
			seenExpiry = true
		default:
			return Data{}, ErrInvalidFormat
		}
	}

	return data, nil
}

func validateRequired(name string, value string) error {
	if value == "" {
		return fmt.Errorf("%w: %s is required", ErrInvalidFormat, name)
	}
	return validateSegment(name, value)
}

func validateSegment(name string, value string) error {
	for _, r := range value {
		if r > 127 || !validSegmentByte(byte(r)) {
			return fmt.Errorf("%w: %s contains %q", ErrInvalidSegment, name, r)
		}
	}
	return nil
}

func validSegmentByte(ch byte) bool {
	return ch >= 'a' && ch <= 'z' ||
		ch >= 'A' && ch <= 'Z' ||
		ch >= '0' && ch <= '9' ||
		ch == '_' || ch == '-' || ch == '.' || ch == '~'
}
