package jobs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewExecuteAllPendingDownloadTask,
	NewUpdateDownloadingAndFailedDownloadTaskStatusToPending,
)
