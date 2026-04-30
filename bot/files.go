package bot

import (
	"context"
	stderrors "errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// GetFileParams contains supported parameters for getFile.
type GetFileParams struct {
	FileID string `json:"file_id"`
}

// GetFile fetches Telegram metadata for a file by file_id.
func (b *Bot) GetFile(ctx context.Context, params GetFileParams) (*telegram.File, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var file telegram.File
	if err := b.call(ctx, "getFile", params, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// DownloadFile downloads a Telegram file_path into w.
func (b *Bot) DownloadFile(ctx context.Context, filePath string, w io.Writer) error {
	if b == nil {
		return stderrors.New("bot is required")
	}
	if ctx == nil {
		return stderrors.New("context is required")
	}
	if w == nil {
		return stderrors.New("writer is required")
	}
	if err := validateFilePath(filePath); err != nil {
		return err
	}

	requestURL, err := b.downloadEndpoint(filePath)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return stderrors.New("create telegram file download request")
	}

	return b.client.Copy(ctx, req, w)
}

func (params GetFileParams) validate() error {
	if strings.TrimSpace(params.FileID) == "" {
		return stderrors.New("file_id is required")
	}

	return nil
}

func (b *Bot) downloadEndpoint(filePath string) (string, error) {
	base, err := url.Parse(b.fileBaseURL)
	if err != nil || base.Scheme == "" || base.Host == "" {
		return "", stderrors.New("invalid telegram file base URL")
	}

	segments := strings.Split(filePath, "/")
	escaped := make([]string, 0, len(segments)+2)
	escaped = append(escaped, "bot"+b.token)
	for _, segment := range segments {
		escaped = append(escaped, url.PathEscape(segment))
	}

	base.Path = strings.TrimRight(base.Path, "/") + "/" + strings.Join(escaped, "/")
	base.RawQuery = ""
	base.Fragment = ""

	return base.String(), nil
}

func validateFilePath(filePath string) error {
	if filePath == "" {
		return stderrors.New("file_path is required")
	}
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		return stderrors.New("file_path must be a Telegram file path, not a URL")
	}
	if strings.HasPrefix(filePath, "/") {
		return stderrors.New("file_path must not start with slash")
	}
	if strings.ContainsAny(filePath, "?#") {
		return stderrors.New("file_path must not contain query or fragment")
	}

	for _, segment := range strings.Split(filePath, "/") {
		if segment == "" {
			return stderrors.New("file_path must not contain empty path segments")
		}
		if segment == ".." {
			return stderrors.New("file_path must not contain parent directory segments")
		}
	}

	return nil
}
