package mq

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/mq/consumer"
	"github.com/tranHieuDev23/GoLoad/internal/dataaccess/mq/producer"
)

var WireSet = wire.NewSet(
	consumer.WireSet,
	producer.WireSet,
)
