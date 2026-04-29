package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"

	apierrors "ai-gram/errors"
)

func (b *Bot) callMultipart(ctx context.Context, method string, fields map[string]string, files map[string]UploadFile, result any) error {
	if b == nil {
		return stderrors.New("bot is required")
	}
	if ctx == nil {
		return stderrors.New("context is required")
	}
	if strings.TrimSpace(method) == "" {
		return stderrors.New("telegram method is required")
	}

	body, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	go func() {
		err := writeMultipart(multipartWriter, fields, files)
		if closeErr := multipartWriter.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			_ = writer.CloseWithError(err)
			return
		}
		_ = writer.Close()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.endpoint(method), body)
	if err != nil {
		_ = body.Close()
		return stderrors.New("create telegram multipart request")
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	responseBody, err := b.client.Do(ctx, req)
	if err != nil {
		_ = body.Close()
		return err
	}

	var response telegramResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("decode telegram API response: %w", err)
	}
	if !response.OK {
		return &apierrors.APIError{
			Code:        response.ErrorCode,
			Description: b.redactToken(response.Description),
			Parameters:  response.Parameters,
		}
	}
	if result == nil || len(response.Result) == 0 {
		return nil
	}
	if err := json.Unmarshal(response.Result, result); err != nil {
		return fmt.Errorf("decode telegram API result: %w", err)
	}

	return nil
}

func writeMultipart(writer *multipart.Writer, fields map[string]string, files map[string]UploadFile) error {
	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			return fmt.Errorf("write multipart field: %w", err)
		}
	}
	for fieldName, file := range files {
		part, err := createFilePart(writer, fieldName, file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, file.Reader); err != nil {
			return fmt.Errorf("copy multipart file: %w", err)
		}
	}

	return nil
}

func createFilePart(writer *multipart.Writer, fieldName string, file UploadFile) (io.Writer, error) {
	if strings.TrimSpace(file.ContentType) == "" {
		part, err := writer.CreateFormFile(fieldName, file.Name)
		if err != nil {
			return nil, fmt.Errorf("create multipart file part: %w", err)
		}
		return part, nil
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q; filename=%q`, fieldName, file.Name))
	header.Set("Content-Type", file.ContentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, fmt.Errorf("create multipart file part: %w", err)
	}

	return part, nil
}
