package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	go_load "github.com/tranHieuDev23/GoLoad/internal/generated/go_load/v1"
	"github.com/tranHieuDev23/GoLoad/internal/logic"
)

const (
	//nolint:gosec // This is just to specify the metadata name
	AuthTokenMetadataName = "GOLOAD_AUTH"
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

func (a Handler) getAuthTokenMetadata(ctx context.Context) string {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	metadataValues := metadata.Get(AuthTokenMetadataName)
	if len(metadataValues) == 0 {
		return ""
	}

	return metadataValues[0]
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
		Token:        a.getAuthTokenMetadata(ctx),
		DownloadType: request.GetDownloadType(),
		URL:          request.GetUrl(),
	})
	if err != nil {
		return nil, err
	}

	return &go_load.CreateDownloadTaskResponse{
		DownloadTask: output.DownloadTask,
	}, nil
}

func (a Handler) CreateSession(
	ctx context.Context,
	request *go_load.CreateSessionRequest,
) (*go_load.CreateSessionResponse, error) {
	output, err := a.accountLogic.CreateSession(ctx, logic.CreateSessionParams{
		AccountName: request.GetAccountName(),
		Password:    request.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	err = grpc.SetHeader(ctx, metadata.Pairs(AuthTokenMetadataName, output.Token))
	if err != nil {
		return nil, err
	}

	return &go_load.CreateSessionResponse{
		Account: output.Account,
	}, nil
}

func (a Handler) DeleteDownloadTask(
	ctx context.Context,
	request *go_load.DeleteDownloadTaskRequest,
) (*go_load.DeleteDownloadTaskResponse, error) {
	if err := a.downloadTaskLogic.DeleteDownloadTask(ctx, logic.DeleteDownloadTaskParams{
		Token:          a.getAuthTokenMetadata(ctx),
		DownloadTaskID: request.GetDownloadTaskId(),
	}); err != nil {
		return nil, err
	}

	return &go_load.DeleteDownloadTaskResponse{}, nil
}

func (a Handler) GetDownloadTaskFile(
	*go_load.GetDownloadTaskFileRequest,
	go_load.GoLoadService_GetDownloadTaskFileServer,
) error {
	panic("unimplemented")
}

func (a Handler) GetDownloadTaskList(
	ctx context.Context,
	request *go_load.GetDownloadTaskListRequest,
) (*go_load.GetDownloadTaskListResponse, error) {
	output, err := a.downloadTaskLogic.GetDownloadTaskList(ctx, logic.GetDownloadTaskListParams{
		Token:  a.getAuthTokenMetadata(ctx),
		Offset: request.GetOffset(),
		Limit:  request.GetLimit(),
	})
	if err != nil {
		return nil, err
	}

	return &go_load.GetDownloadTaskListResponse{
		DownloadTaskList:       output.DownloadTaskList,
		TotalDownloadTaskCount: output.TotalDownloadTaskCount,
	}, nil
}

func (a Handler) UpdateDownloadTask(
	ctx context.Context,
	request *go_load.UpdateDownloadTaskRequest,
) (*go_load.UpdateDownloadTaskResponse, error) {
	output, err := a.downloadTaskLogic.UpdateDownloadTask(ctx, logic.UpdateDownloadTaskParams{
		Token:          a.getAuthTokenMetadata(ctx),
		DownloadTaskID: request.GetDownloadTaskId(),
		URL:            request.GetUrl(),
	})
	if err != nil {
		return nil, err
	}

	return &go_load.UpdateDownloadTaskResponse{
		DownloadTask: output.DownloadTask,
	}, nil
}
