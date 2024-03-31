package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	go_load "github.com/tranHieuDev23/GoLoad/internal/generated/go_load/v1"
	"github.com/tranHieuDev23/GoLoad/internal/utils"
)

var (
	TabNameDownloadTasks = goqu.T("download_tasks")

	ErrDownloadTaskNotFound = status.Error(codes.NotFound, "download task not found")
)

const (
	ColNameDownloadTaskID             = "id"
	ColNameDownloadTaskOfAccountID    = "of_account_id"
	ColNameDownloadTaskDownloadType   = "download_type"
	ColNameDownloadTaskURL            = "url"
	ColNameDownloadTaskDownloadStatus = "download_status"
	ColNameDownloadTaskMetadata       = "metadata"
)

type DownloadTask struct {
	ID             uint64                 `db:"id" goqu:"skipinsert,skipupdate"`
	OfAccountID    uint64                 `db:"of_account_id" goqu:"skipupdate"`
	DownloadType   go_load.DownloadType   `db:"download_type"`
	URL            string                 `db:"url"`
	DownloadStatus go_load.DownloadStatus `db:"download_status"`
	Metadata       JSON                   `db:"metadata"`
}

type DownloadTaskDataAccessor interface {
	CreateDownloadTask(ctx context.Context, task DownloadTask) (uint64, error)
	GetDownloadTaskListOfAccount(ctx context.Context, accountID, offset, limit uint64) ([]DownloadTask, error)
	GetDownloadTaskCountOfAccount(ctx context.Context, accountID uint64) (uint64, error)
	GetDownloadTask(ctx context.Context, id uint64) (DownloadTask, error)
	GetDownloadTaskWithXLock(ctx context.Context, id uint64) (DownloadTask, error)
	UpdateDownloadTask(ctx context.Context, task DownloadTask) error
	DeleteDownloadTask(ctx context.Context, id uint64) error
	GetPendingDownloadTaskIDList(ctx context.Context) ([]uint64, error)
	UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error
	WithDatabase(database Database) DownloadTaskDataAccessor
}

type downloadTaskDataAccessor struct {
	database Database
	logger   *zap.Logger
}

func NewDownloadTaskDataAccessor(
	database *goqu.Database,
	logger *zap.Logger,
) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   logger,
	}
}

func (d downloadTaskDataAccessor) CreateDownloadTask(ctx context.Context, task DownloadTask) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("task", task))

	result, err := d.database.
		Insert(TabNameDownloadTasks).
		Rows(task).
		Executor().
		ExecContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to create download task")
		return 0, status.Error(codes.Internal, "failed to create download task")
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get last inserted id")
		return 0, status.Error(codes.Internal, "failed to get last inserted id")
	}

	return uint64(lastInsertedID), nil
}

func (d downloadTaskDataAccessor) DeleteDownloadTask(ctx context.Context, id uint64) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	if _, err := d.database.
		Delete(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskID: id}).
		Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to delete download task")
		return status.Error(codes.Internal, "failed to delete download task")
	}

	return nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskCountOfAccount(ctx context.Context, accountID uint64) (uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("account_id", accountID))

	count, err := d.database.
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskOfAccountID: accountID}).
		CountContext(ctx)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to count download task of user")
		return 0, status.Error(codes.Internal, "failed to count download task of user")
	}

	return uint64(count), nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskListOfAccount(
	ctx context.Context,
	accountID uint64,
	offset uint64,
	limit uint64,
) ([]DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).
		With(zap.Uint64("account_id", accountID)).
		With(zap.Uint64("offset", offset)).
		With(zap.Uint64("limit", limit))

	downloadTaskList := make([]DownloadTask, 0)
	if err := d.database.
		Select().
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameAccountPasswordsOfAccountID: accountID}).
		Offset(uint(offset)).
		Limit(uint(limit)).
		Executor().
		ScanStructsContext(ctx, &downloadTaskList); err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task list of account")
		return nil, status.Error(codes.Internal, "failed to get download task list of account")
	}

	return downloadTaskList, nil
}

func (d downloadTaskDataAccessor) GetDownloadTask(ctx context.Context, id uint64) (DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	downloadTask := DownloadTask{}
	found, err := d.database.
		Select().
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskID: id}).
		ScanStructContext(ctx, &downloadTask)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task")
		return DownloadTask{}, status.Error(codes.Internal, "failed to get download task list of account")
	}

	if !found {
		logger.Error("download task not found")
		return DownloadTask{}, ErrDownloadTaskNotFound
	}

	return downloadTask, nil
}

func (d downloadTaskDataAccessor) GetDownloadTaskWithXLock(ctx context.Context, id uint64) (DownloadTask, error) {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Uint64("id", id))

	downloadTask := DownloadTask{}
	found, err := d.database.
		Select().
		From(TabNameDownloadTasks).
		Where(goqu.Ex{ColNameDownloadTaskID: id}).
		ForUpdate(goqu.Wait).
		ScanStructContext(ctx, &downloadTask)
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to get download task")
		return DownloadTask{}, status.Error(codes.Internal, "failed to get download task list of account")
	}

	if !found {
		logger.Error("download task not found")
		return DownloadTask{}, ErrDownloadTaskNotFound
	}

	return downloadTask, nil
}

func (d downloadTaskDataAccessor) UpdateDownloadTask(ctx context.Context, task DownloadTask) error {
	logger := utils.LoggerWithContext(ctx, d.logger).With(zap.Any("task", task))

	if _, err := d.database.
		Update(TabNameDownloadTasks).
		Set(task).
		Where(goqu.Ex{ColNameDownloadTaskID: task.ID}).
		Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to update download task")
		return status.Error(codes.Internal, "failed to update download task")
	}

	return nil
}

func (d downloadTaskDataAccessor) GetPendingDownloadTaskIDList(ctx context.Context) ([]uint64, error) {
	logger := utils.LoggerWithContext(ctx, d.logger)

	downloadTaskIDList := make([]uint64, 0)
	if err := d.database.
		Select(ColNameDownloadTaskID).
		From(TabNameDownloadTasks).
		Where(goqu.Ex{
			ColNameDownloadTaskDownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		}).
		ScanValsContext(ctx, &downloadTaskIDList); err != nil {
		logger.With(zap.Error(err)).Error("failed to get pending download task id list")
		return nil, status.Error(codes.Internal, "failed to get pending download task id list")
	}

	return downloadTaskIDList, nil
}

func (d downloadTaskDataAccessor) UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error {
	logger := utils.LoggerWithContext(ctx, d.logger)

	if _, err := d.database.
		Update(TabNameDownloadTasks).
		Set(goqu.Record{
			ColNameDownloadTaskDownloadStatus: go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING,
		}).
		Where(
			goqu.C(ColNameDownloadTaskDownloadStatus).
				In(go_load.DownloadStatus_DOWNLOAD_STATUS_PENDING, go_load.DownloadStatus_DOWNLOAD_STATUS_FAILED),
		).Executor().
		ExecContext(ctx); err != nil {
		logger.With(zap.Error(err)).Error("failed to update downloading and failed download task status to pending")
		return status.Error(codes.Internal, "failed to update downloading and failed download task status to pending")
	}

	return nil
}

func (d downloadTaskDataAccessor) WithDatabase(database Database) DownloadTaskDataAccessor {
	return &downloadTaskDataAccessor{
		database: database,
		logger:   d.logger,
	}
}
