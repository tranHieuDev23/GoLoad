package configs

type ExecuteAllPendingDownloadTask struct {
	Schedule         string `yaml:"schedule"`
	ConcurrencyLimit int    `yaml:"concurrency_limit"`
}

type UpdateDownloadingAndFailedDownloadTaskStatusToPending struct {
	Schedule string `yaml:"schedule"`
}

//nolint:lll // Long field names
type Cron struct {
	ExecuteAllPendingDownloadTask                         ExecuteAllPendingDownloadTask                         `yaml:"execute_all_pending_download_task"`
	UpdateDownloadingAndFailedDownloadTaskStatusToPending UpdateDownloadingAndFailedDownloadTaskStatusToPending `yaml:"update_downloading_and_failed_download_task_status_to_pending"`
}
