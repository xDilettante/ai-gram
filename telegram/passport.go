package telegram

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strings"
)

const (
	passportElementErrorSourceDataField        = "data"
	passportElementErrorSourceFrontSide        = "front_side"
	passportElementErrorSourceReverseSide      = "reverse_side"
	passportElementErrorSourceSelfie           = "selfie"
	passportElementErrorSourceFile             = "file"
	passportElementErrorSourceFiles            = "files"
	passportElementErrorSourceTranslationFile  = "translation_file"
	passportElementErrorSourceTranslationFiles = "translation_files"
	passportElementErrorSourceUnspecified      = "unspecified"
)

// PassportData describes Telegram Passport data shared with the bot by the user.
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

// PassportFile represents a file uploaded to Telegram Passport.
type PassportFile struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	FileDate     int64  `json:"file_date"`
}

// EncryptedPassportElement describes a Telegram Passport element shared with the bot.
type EncryptedPassportElement struct {
	Type        string         `json:"type"`
	Data        string         `json:"data,omitempty"`
	PhoneNumber string         `json:"phone_number,omitempty"`
	Email       string         `json:"email,omitempty"`
	Files       []PassportFile `json:"files,omitempty"`
	FrontSide   *PassportFile  `json:"front_side,omitempty"`
	ReverseSide *PassportFile  `json:"reverse_side,omitempty"`
	Selfie      *PassportFile  `json:"selfie,omitempty"`
	Translation []PassportFile `json:"translation,omitempty"`
	Hash        string         `json:"hash"`
}

// EncryptedCredentials describes data required for decrypting and authenticating Telegram Passport elements.
type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

// PassportElementError marks Telegram Passport error objects sent to setPassportDataErrors.
type PassportElementError interface {
	passportElementError()
}

// PassportElementErrorDataField represents an issue in one of the submitted data fields.
type PassportElementErrorDataField struct {
	Source    string `json:"source"`
	Type      string `json:"type"`
	FieldName string `json:"field_name"`
	DataHash  string `json:"data_hash"`
	Message   string `json:"message"`
}

// PassportElementErrorFrontSide represents an issue with the front side of a document.
type PassportElementErrorFrontSide struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorReverseSide represents an issue with the reverse side of a document.
type PassportElementErrorReverseSide struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorSelfie represents an issue with the selfie with a document.
type PassportElementErrorSelfie struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorFile represents an issue with a document scan.
type PassportElementErrorFile struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorFiles represents an issue with a list of document scans.
type PassportElementErrorFiles struct {
	Source     string   `json:"source"`
	Type       string   `json:"type"`
	FileHashes []string `json:"file_hashes"`
	Message    string   `json:"message"`
}

// PassportElementErrorTranslationFile represents an issue with one translated document file.
type PassportElementErrorTranslationFile struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorTranslationFiles represents an issue with translated document files.
type PassportElementErrorTranslationFiles struct {
	Source     string   `json:"source"`
	Type       string   `json:"type"`
	FileHashes []string `json:"file_hashes"`
	Message    string   `json:"message"`
}

// PassportElementErrorUnspecified represents an issue in an unspecified place.
type PassportElementErrorUnspecified struct {
	Source      string `json:"source"`
	Type        string `json:"type"`
	ElementHash string `json:"element_hash"`
	Message     string `json:"message"`
}

func (PassportElementErrorDataField) passportElementError()        {}
func (PassportElementErrorFrontSide) passportElementError()        {}
func (PassportElementErrorReverseSide) passportElementError()      {}
func (PassportElementErrorSelfie) passportElementError()           {}
func (PassportElementErrorFile) passportElementError()             {}
func (PassportElementErrorFiles) passportElementError()            {}
func (PassportElementErrorTranslationFile) passportElementError()  {}
func (PassportElementErrorTranslationFiles) passportElementError() {}
func (PassportElementErrorUnspecified) passportElementError()      {}

// MarshalJSON encodes PassportElementErrorDataField with the required source field.
func (e PassportElementErrorDataField) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceDataField
	type alias PassportElementErrorDataField
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorFrontSide with the required source field.
func (e PassportElementErrorFrontSide) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceFrontSide
	type alias PassportElementErrorFrontSide
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorReverseSide with the required source field.
func (e PassportElementErrorReverseSide) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceReverseSide
	type alias PassportElementErrorReverseSide
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorSelfie with the required source field.
func (e PassportElementErrorSelfie) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceSelfie
	type alias PassportElementErrorSelfie
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorFile with the required source field.
func (e PassportElementErrorFile) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceFile
	type alias PassportElementErrorFile
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorFiles with the required source field.
func (e PassportElementErrorFiles) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceFiles
	type alias PassportElementErrorFiles
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorTranslationFile with the required source field.
func (e PassportElementErrorTranslationFile) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceTranslationFile
	type alias PassportElementErrorTranslationFile
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorTranslationFiles with the required source field.
func (e PassportElementErrorTranslationFiles) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceTranslationFiles
	type alias PassportElementErrorTranslationFiles
	return json.Marshal(alias(e))
}

// MarshalJSON encodes PassportElementErrorUnspecified with the required source field.
func (e PassportElementErrorUnspecified) MarshalJSON() ([]byte, error) {
	e.Source = passportElementErrorSourceUnspecified
	type alias PassportElementErrorUnspecified
	return json.Marshal(alias(e))
}

// ValidatePassportElementErrors checks whether passport errors can be sent to Telegram.
func ValidatePassportElementErrors(errors []PassportElementError) error {
	for index, elementError := range errors {
		if err := ValidatePassportElementError(elementError); err != nil {
			return fmt.Errorf("errors[%d] is invalid: %w", index, err)
		}
	}
	return nil
}

// ValidatePassportElementError checks whether a passport error can be sent to Telegram.
func ValidatePassportElementError(elementError PassportElementError) error {
	if elementError == nil || isNilInterfaceValue(elementError) {
		return stderrors.New("passport element error must not be nil")
	}

	switch value := elementError.(type) {
	case PassportElementErrorDataField:
		return validatePassportElementErrorDataField(value)
	case *PassportElementErrorDataField:
		return validatePassportElementErrorDataField(*value)
	case PassportElementErrorFrontSide:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case *PassportElementErrorFrontSide:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case PassportElementErrorReverseSide:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case *PassportElementErrorReverseSide:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case PassportElementErrorSelfie:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case *PassportElementErrorSelfie:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case PassportElementErrorFile:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case *PassportElementErrorFile:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case PassportElementErrorFiles:
		return validatePassportElementErrorFileHashes(value.Type, value.FileHashes, value.Message)
	case *PassportElementErrorFiles:
		return validatePassportElementErrorFileHashes(value.Type, value.FileHashes, value.Message)
	case PassportElementErrorTranslationFile:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case *PassportElementErrorTranslationFile:
		return validatePassportElementErrorSingleFile(value.Type, value.FileHash, value.Message)
	case PassportElementErrorTranslationFiles:
		return validatePassportElementErrorFileHashes(value.Type, value.FileHashes, value.Message)
	case *PassportElementErrorTranslationFiles:
		return validatePassportElementErrorFileHashes(value.Type, value.FileHashes, value.Message)
	case PassportElementErrorUnspecified:
		return validatePassportElementErrorUnspecified(value)
	case *PassportElementErrorUnspecified:
		return validatePassportElementErrorUnspecified(*value)
	default:
		return stderrors.New("unsupported passport element error")
	}
}

func validatePassportElementErrorDataField(elementError PassportElementErrorDataField) error {
	if err := validatePassportErrorTypeAndMessage(elementError.Type, elementError.Message); err != nil {
		return err
	}
	if strings.TrimSpace(elementError.FieldName) == "" {
		return stderrors.New("field_name is required")
	}
	if strings.TrimSpace(elementError.DataHash) == "" {
		return stderrors.New("data_hash is required")
	}
	return nil
}

func validatePassportElementErrorSingleFile(elementType string, fileHash string, message string) error {
	if err := validatePassportErrorTypeAndMessage(elementType, message); err != nil {
		return err
	}
	if strings.TrimSpace(fileHash) == "" {
		return stderrors.New("file_hash is required")
	}
	return nil
}

func validatePassportElementErrorFileHashes(elementType string, fileHashes []string, message string) error {
	if err := validatePassportErrorTypeAndMessage(elementType, message); err != nil {
		return err
	}
	if len(fileHashes) == 0 {
		return stderrors.New("file_hashes must not be empty")
	}
	for index, hash := range fileHashes {
		if strings.TrimSpace(hash) == "" {
			return fmt.Errorf("file_hashes[%d] is required", index)
		}
	}
	return nil
}

func validatePassportElementErrorUnspecified(elementError PassportElementErrorUnspecified) error {
	if err := validatePassportErrorTypeAndMessage(elementError.Type, elementError.Message); err != nil {
		return err
	}
	if strings.TrimSpace(elementError.ElementHash) == "" {
		return stderrors.New("element_hash is required")
	}
	return nil
}

func validatePassportErrorTypeAndMessage(elementType string, message string) error {
	if strings.TrimSpace(elementType) == "" {
		return stderrors.New("type is required")
	}
	if strings.TrimSpace(message) == "" {
		return stderrors.New("message is required")
	}
	return nil
}
