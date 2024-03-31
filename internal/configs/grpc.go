package configs

import (
	"github.com/dustin/go-humanize"
)

type GetDownloadTaskFile struct {
	ResponseBufferSize string `yaml:"response_buffer_size"`
}

func (g GetDownloadTaskFile) GetResponseBufferSizeInBytes() (uint64, error) {
	return humanize.ParseBytes(g.ResponseBufferSize)
}

type GRPC struct {
	Address             string              `yaml:"address"`
	GetDownloadTaskFile GetDownloadTaskFile `yaml:"get_download_task_file"`
}
