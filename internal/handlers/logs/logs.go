package logs

import (
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type logs struct {
	data             map[int64]map[string]interface{}
	realtimeLogs     map[int64]context.CancelFunc
	messagesToDelete map[int64][]int64

	database *gorm.DB
	logger   *zap.Logger
}

func New(dispatcher *ext.Dispatcher, database *gorm.DB, logger *zap.Logger) {
	l := &logs{
		data:             make(map[int64]map[string]interface{}),
		realtimeLogs:     make(map[int64]context.CancelFunc),
		messagesToDelete: make(map[int64][]int64),

		database: database,
		logger:   logger.Named("logs"),
	}
	dispatcher.AddHandlerToGroup(handlers.NewCommand("logs", l.logs), 10)
	dispatcher.AddHandlerToGroup(handlers.NewConversation(
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
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("files:"), l.files), 10)
}

func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}
