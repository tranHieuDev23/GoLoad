package grpc

import (
	"context"

	"github.com/tranHieuDev23/GoLoad/internal/generated/grpc/go_load"
	"github.com/tranHieuDev23/GoLoad/internal/logic"
)

type Handler struct {
	go_load.UnimplementedGoLoadServiceServer
	accountLogic      logic.Account
	downloadTaskLogic logic.DownloadTask
}

func NewHandler(
	accountLogic logic.Account,
	downloadTaskLogic logic.DownloadTask,
) go_load.GoLoadServiceServer {
	return &Handler{
		accountLogic:      accountLogic,
		downloadTaskLogic: downloadTaskLogic,
	}
}

func (a Handler) CreateAccount(
	ctx context.Context,
	request *go_load.CreateAccountRequest,
) (*go_load.CreateAccountResponse, error) {
	output, err := a.accountLogic.CreateAccount(ctx, logic.CreateAccountParams{
		AccountName: request.GetAccountName(),
		Password:    request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &go_load.CreateAccountResponse{
		AccountId: output.ID,
	}, nil
}

func (a Handler) CreateDownloadTask(
	ctx context.Context,
	request *go_load.CreateDownloadTaskRequest,
) (*go_load.CreateDownloadTaskResponse, error) {
	output, err := a.downloadTaskLogic.CreateDownloadTask(ctx, logic.CreateDownloadTaskParams{
		Token:        request.GetToken(),
		DownloadType: request.GetDownloadType(),
		URL:          request.GetUrl(),
	})
	if err != nil {
		return nil, err
	}

	return &go_load.CreateDownloadTaskResponse{
		DownloadTask: &output.DownloadTask,
	}, nil
}

func (a Handler) CreateSession(
	ctx context.Context,
	request *go_load.CreateSessionRequest,
) (*go_load.CreateSessionResponse, error) {
	token, err := a.accountLogic.CreateSession(ctx, logic.CreateSessionParams{
		AccountName: request.GetAccountName(),
		Password:    request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &go_load.CreateSessionResponse{
		Token: token,
	}, nil
}

func (a Handler) DeleteDownloadTask(
	context.Context,
	*go_load.DeleteDownloadTaskRequest,
) (*go_load.DeleteDownloadTaskResponse, error) {
	panic("unimplemented")
}

func (a Handler) GetDownloadTaskFile(
	*go_load.GetDownloadTaskFileRequest,
	go_load.GoLoadService_GetDownloadTaskFileServer,
) error {
	panic("unimplemented")
}

func (a Handler) GetDownloadTaskList(
	context.Context,
	*go_load.GetDownloadTaskListRequest,
) (*go_load.GetDownloadTaskListResponse, error) {
	panic("unimplemented")
}

func (a Handler) UpdateDownloadTask(
	context.Context,
	*go_load.UpdateDownloadTaskRequest,
) (*go_load.UpdateDownloadTaskResponse, error) {
	panic("unimplemented")
}
