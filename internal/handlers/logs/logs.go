package logs

import (
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type logs struct {
	data         map[string]interface{}
	realtimeLogs map[uint]context.CancelFunc

	dispatcher *ext.Dispatcher
	database   *gorm.DB
	logger     *zap.Logger
}

func New(dispatcher *ext.Dispatcher, database *gorm.DB, logger *zap.Logger) {
	l := &logs{
		data:         make(map[string]interface{}),
		realtimeLogs: make(map[uint]context.CancelFunc),

		dispatcher: dispatcher,
		database:   database,
		logger:     logger.Named("logs"),
	}

	l.setupHandlers()
}
