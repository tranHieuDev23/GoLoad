package jobs

import (
	"context"

	"github.com/tranHieuDev23/GoLoad/internal/logic"
)

type ExecuteAllPendingDownloadTask interface {
	Run(context.Context) error
}

type executeAllPendingDownloadTask struct {
	downloadTaskLogic logic.DownloadTask
}

func NewExecuteAllPendingDownloadTask(
	downloadTaskLogic logic.DownloadTask,
) ExecuteAllPendingDownloadTask {
	return &executeAllPendingDownloadTask{
		downloadTaskLogic: downloadTaskLogic,
	}
}

func (e executeAllPendingDownloadTask) Run(ctx context.Context) error {
	return e.downloadTaskLogic.ExecuteAllPendingDownloadTask(ctx)
}
