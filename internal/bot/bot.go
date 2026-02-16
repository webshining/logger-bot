package bot

import (
	"bot/internal/handlers/filters"
	"bot/internal/handlers/logs"
	"bot/internal/settings"
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bot struct {
	bot        *gotgbot.Bot
	dispatcher *ext.Dispatcher
	logger     *zap.Logger
}

func New(settings settings.BotConfig, database *gorm.DB, logger *zap.Logger) *Bot {
	bot, err := gotgbot.NewBot(settings.Token, nil)
	if err != nil {
		logger.Fatal("failed to create new bot:", zap.Error(err))
	}

	dispatcher := ext.NewDispatcher(nil)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.All, filters.Admin(settings.Admin)), -10)
	logs.New(dispatcher, database, logger)

	return &Bot{bot, dispatcher, logger.Named("bot")}
}

func (b *Bot) Run(ctx context.Context) {
	defer b.bot.Close(nil)

	updater := ext.NewUpdater(b.dispatcher, nil)
	if err := updater.StartPolling(b.bot, nil); err != nil {
		b.logger.Fatal("failed to start polling", zap.Error(err))
		return
	}

	b.logger.Info("Bot is now running", zap.String("username", b.bot.Username))
	defer b.logger.Info("Bot is shuted down")

	<-ctx.Done()
}
