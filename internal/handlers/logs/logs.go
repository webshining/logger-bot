package logs

import (
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
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

func (l *logs) setupHandlers() {
	l.dispatcher.AddHandlerToGroup(handlers.NewCommand("logs", l.files), 10)
	l.dispatcher.AddHandlerToGroup(handlers.NewConversation(
		[]ext.Handler{handlers.NewCallback(callbackquery.Equal("files:add"), l.addFile)},
		map[string][]ext.Handler{
			"NAME": {handlers.NewMessage(noCommands, l.addFileName)},
			"PATH": {handlers.NewMessage(noCommands, l.addFilePath)},
		},
		&handlers.ConversationOpts{
			Exits:        []ext.Handler{handlers.NewCommand("cancel", l.addFileCancel)},
			StateStorage: conversation.NewInMemoryStorage(conversation.KeyStrategySenderAndChat),
			AllowReEntry: true,
		},
	), 10)
	l.dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("files:"), l.filesCallback), 10)
	l.dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("file:"), l.fileCallback), 10)
}
