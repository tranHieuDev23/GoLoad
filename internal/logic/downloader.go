package logic

import (
	"context"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

const (
	HTTPResponseHeaderContentType = "Content-Type"
	HTTPMetadataKeyContentType    = "content-type"
)

type Downloader interface {
	Download(ctx context.Context, writer io.Writer) (map[string]any, error)
}

type HTTPDownloader struct {
	url    string
	logger *zap.Logger
}

func NewHTTPDownloader(
	url string,
	logger *zap.Logger,
) Downloader {
	return &HTTPDownloader{
		url:    url,
		logger: logger,
	}
}

func (h HTTPDownloader) Download(ctx context.Context, writer io.Writer) (map[string]any, error) {
	logger := utils.LoggerWithContext(ctx, h.logger)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, http.NoBody)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create http get request")
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to make http get request")
		return nil, err
	}

	defer response.Body.Close()

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to read response and write to writer")
		return nil, err
	}

	metadata := map[string]any{
		HTTPMetadataKeyContentType: response.Header.Get(HTTPResponseHeaderContentType),
	}

	return metadata, nil
}
