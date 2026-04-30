package telegram

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPassportDataMessageDecoding(t *testing.T) {
	payload := []byte(`{
		"message_id":7,
		"chat":{"id":123,"type":"private"},
		"date":100,
		"passport_data":{
			"data":[{
				"type":"passport",
				"data":"enc-data",
				"files":[{"file_id":"file-1","file_unique_id":"unique-1","file_size":10,"file_date":101}],
				"front_side":{"file_id":"front","file_unique_id":"front-u","file_size":11,"file_date":102},
				"reverse_side":{"file_id":"reverse","file_unique_id":"reverse-u","file_size":12,"file_date":103},
				"selfie":{"file_id":"selfie","file_unique_id":"selfie-u","file_size":13,"file_date":104},
				"translation":[{"file_id":"translation","file_unique_id":"translation-u","file_size":14,"file_date":105}],
				"hash":"hash-1"
			},{
				"type":"phone_number",
				"phone_number":"+10000000000",
				"hash":"hash-2"
			},{
				"type":"email",
				"email":"user@example.test",
				"hash":"hash-3"
			}],
			"credentials":{"data":"cred-data","hash":"cred-hash","secret":"cred-secret"}
		}
	}`)

	var message Message
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if message.PassportData == nil {
		t.Fatal("expected passport data")
	}
	if len(message.PassportData.Data) != 3 {
		t.Fatalf("unexpected element count: %d", len(message.PassportData.Data))
	}
	first := message.PassportData.Data[0]
	if first.Type != "passport" || first.Data != "enc-data" || first.FrontSide == nil || first.ReverseSide == nil || first.Selfie == nil || len(first.Files) != 1 || len(first.Translation) != 1 || first.Hash != "hash-1" {
		t.Fatalf("unexpected passport element: %+v", first)
	}
	if message.PassportData.Data[1].PhoneNumber != "+10000000000" || message.PassportData.Data[2].Email != "user@example.test" {
		t.Fatalf("unexpected simple elements: %+v", message.PassportData.Data)
	}
	if message.PassportData.Credentials.Data != "cred-data" || message.PassportData.Credentials.Hash != "cred-hash" || message.PassportData.Credentials.Secret != "cred-secret" {
		t.Fatalf("unexpected credentials: %+v", message.PassportData.Credentials)
	}
}

func TestPassportElementErrorMarshalSources(t *testing.T) {
	tests := []struct {
		name     string
		error    PassportElementError
		source   string
		required map[string]any
	}{
		{name: "data field", error: PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", DataHash: "data-hash", Message: "invalid"}, source: "data", required: map[string]any{"field_name": "first_name", "data_hash": "data-hash"}},
		{name: "front side", error: PassportElementErrorFrontSide{Type: "passport", FileHash: "file-hash", Message: "invalid"}, source: "front_side", required: map[string]any{"file_hash": "file-hash"}},
		{name: "reverse side", error: PassportElementErrorReverseSide{Type: "driver_license", FileHash: "file-hash", Message: "invalid"}, source: "reverse_side", required: map[string]any{"file_hash": "file-hash"}},
		{name: "selfie", error: PassportElementErrorSelfie{Type: "passport", FileHash: "file-hash", Message: "invalid"}, source: "selfie", required: map[string]any{"file_hash": "file-hash"}},
		{name: "file", error: PassportElementErrorFile{Type: "utility_bill", FileHash: "file-hash", Message: "invalid"}, source: "file", required: map[string]any{"file_hash": "file-hash"}},
		{name: "files", error: PassportElementErrorFiles{Type: "utility_bill", FileHashes: []string{"hash-1", "hash-2"}, Message: "invalid"}, source: "files", required: map[string]any{"file_hashes": []any{"hash-1", "hash-2"}}},
		{name: "translation file", error: PassportElementErrorTranslationFile{Type: "passport", FileHash: "file-hash", Message: "invalid"}, source: "translation_file", required: map[string]any{"file_hash": "file-hash"}},
		{name: "translation files", error: PassportElementErrorTranslationFiles{Type: "passport", FileHashes: []string{"hash-1"}, Message: "invalid"}, source: "translation_files", required: map[string]any{"file_hashes": []any{"hash-1"}}},
		{name: "unspecified", error: PassportElementErrorUnspecified{Type: "passport", ElementHash: "element-hash", Message: "invalid"}, source: "unspecified", required: map[string]any{"element_hash": "element-hash"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePassportElementError(tt.error); err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
			body, err := json.Marshal(tt.error)
			if err != nil {
				t.Fatalf("marshal error: %v", err)
			}
			var payload map[string]any
			if err := json.Unmarshal(body, &payload); err != nil {
				t.Fatalf("decode marshaled error: %v", err)
			}
			if payload["source"] != tt.source || payload["type"] == "" || payload["message"] == "" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
			for name, want := range tt.required {
				got := payload[name]
				switch wantValue := want.(type) {
				case []any:
					gotValues, ok := got.([]any)
					if !ok || len(gotValues) != len(wantValue) {
						t.Fatalf("unexpected %s: %#v", name, got)
					}
				default:
					if got != wantValue {
						t.Fatalf("unexpected %s: %#v", name, got)
					}
				}
			}
		})
	}
}

func TestPassportElementErrorValidation(t *testing.T) {
	tests := []struct {
		name  string
		error PassportElementError
	}{
		{name: "nil", error: nil},
		{name: "missing type", error: PassportElementErrorDataField{FieldName: "first_name", DataHash: "hash", Message: "invalid"}},
		{name: "missing message", error: PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", DataHash: "hash"}},
		{name: "data field missing field", error: PassportElementErrorDataField{Type: "personal_details", DataHash: "hash", Message: "invalid"}},
		{name: "data field missing hash", error: PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", Message: "invalid"}},
		{name: "single file missing hash", error: PassportElementErrorFile{Type: "utility_bill", Message: "invalid"}},
		{name: "files empty", error: PassportElementErrorFiles{Type: "utility_bill", Message: "invalid"}},
		{name: "files blank hash", error: PassportElementErrorFiles{Type: "utility_bill", FileHashes: []string{"hash", " "}, Message: "invalid"}},
		{name: "unspecified missing element hash", error: PassportElementErrorUnspecified{Type: "passport", Message: "invalid"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassportElementError(tt.error)
			if err == nil {
				t.Fatal("expected error")
			}
			if strings.Contains(err.Error(), "enc-data") || strings.Contains(err.Error(), "cred-secret") {
				t.Fatalf("passport payload leaked in error: %v", err)
			}
		})
	}
}
