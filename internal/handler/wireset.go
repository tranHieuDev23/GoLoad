package handler

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/GoLoad/internal/handler/consumers"
	"github.com/tranHieuDev23/GoLoad/internal/handler/grpc"
	"github.com/tranHieuDev23/GoLoad/internal/handler/http"
)

var WireSet = wire.NewSet(
	grpc.WireSet,
	http.WireSet,
	consumers.WireSet,
)
