package bot

import (
	"bot/internal/handlers/filters"
	"bot/internal/handlers/logs"
	"bot/internal/settings"
	"context"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bot struct {
	bot        *gotgbot.Bot
	dispatcher *ext.Dispatcher
	updater    *ext.Updater
	webhook    bool
	logger     *zap.Logger
}

func New(settings settings.BotConfig, database *gorm.DB, logger *zap.Logger) *Bot {
	bot, err := gotgbot.NewBot(settings.Token, nil)
	if err != nil {
		logger.Fatal("failed to create new bot:", zap.Error(err))
	}

	dispatcher := ext.NewDispatcher(nil)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.All, filters.Admin(settings.Admin)), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.All, filters.Admin(settings.Admin)), -10)
	logs.New(dispatcher, database, logger)

	updater := ext.NewUpdater(dispatcher, nil)
	webhook := false
	if settings.Webhook.Domain != nil && settings.Webhook.Secret != nil && settings.Webhook.Path != nil && settings.Webhook.Host != nil && settings.Webhook.Port != nil {
		if err := updater.StartWebhook(bot, *settings.Webhook.Path, ext.WebhookOpts{ListenAddr: fmt.Sprintf("%s:%d", *settings.Webhook.Host, *settings.Webhook.Port), SecretToken: *settings.Webhook.Secret}); err != nil {
			logger.Fatal("failed to start webhook", zap.Error(err))
		}
		if err := updater.SetAllBotWebhooks(*settings.Webhook.Domain, &gotgbot.SetWebhookOpts{DropPendingUpdates: true, SecretToken: *settings.Webhook.Secret}); err != nil {
			logger.Fatal("failed to set webhook", zap.Error(err))
		}
		webhook = true
	}

	return &Bot{bot, dispatcher, updater, webhook, logger.Named("bot")}
}

func (b *Bot) Run(ctx context.Context) {
	defer b.bot.Close(nil)

	if !b.webhook {
		if err := b.updater.StartPolling(b.bot, &ext.PollingOpts{EnableWebhookDeletion: true}); err != nil {
			b.logger.Fatal("failed to start polling", zap.Error(err))
			return
		}
	}

	b.logger.Info("Bot is now running", zap.String("username", b.bot.Username), zap.Bool("webhook", b.webhook))
	defer b.logger.Info("Bot is shuted down")

	<-ctx.Done()
}
