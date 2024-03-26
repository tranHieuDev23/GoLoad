package configs

type DownloadMode string

const (
	DownloadModeLocal DownloadMode = "local"
	DownloadModeS3    DownloadMode = "s3"
)

type Download struct {
	Mode              DownloadMode `yaml:"mode"`
	DownloadDirectory string       `yaml:"download_directory"`
	Bucket            string       `yaml:"bucket"`
	Address           string       `yaml:"address"`
	Username          string       `yaml:"username"`
	Password          string       `yaml:"password"`
}
