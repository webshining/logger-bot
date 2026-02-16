package logs

import (
	"bot/internal/database/models"
	"fmt"
	"html"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"go.uber.org/zap"
)

func (l *logs) logs(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(b, "Select logs to display:", &gotgbot.SendMessageOpts{ReplyMarkup: l.logsMarkup()})
	return nil
}

func (l *logs) files(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.CallbackQuery.Answer(b, nil)

	data := strings.Split(ctx.CallbackQuery.Data, ":")
	var file models.File
	l.database.First(&file, "id = ?", data[1])

	go l.watcher(b, ctx.EffectiveChat.Id, file)
	return nil
}

func (l *logs) addFile(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.CallbackQuery.Answer(b, nil)

	message, _ := ctx.EffectiveMessage.Reply(b, "`[` Add new file `]`\n\nEnter display name:", &gotgbot.SendMessageOpts{ParseMode: "MarkdownV2"})
	l.addMessageToDelete(message)

	return handlers.NextConversationState("NAME")
}

func (l *logs) addFileName(b *gotgbot.Bot, ctx *ext.Context) error {
	name := ctx.EffectiveMessage.Text
	l.data[ctx.EffectiveChat.Id] = map[string]interface{}{"file_name": name}

	message, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("`[` Add new file %s `]`\n\nEnter file path:", html.EscapeString(name)), &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
	if err != nil {
		l.logger.Error("failed", zap.Error(err))
	}
	l.addMessageToDelete(message)
	l.addMessageToDelete(ctx.EffectiveMessage)

	return handlers.NextConversationState("PATH")
}

func (l *logs) addFilePath(b *gotgbot.Bot, ctx *ext.Context) error {
	name := l.data[ctx.EffectiveChat.Id]["file_name"].(string)
	path := ctx.EffectiveMessage.Text
	l.database.Create(&models.File{Name: name, Path: path})
	l.addMessageToDelete(ctx.EffectiveMessage)

	delete(l.data, ctx.EffectiveChat.Id)
	l.deleteMessages(b, ctx.EffectiveChat.Id)
	return handlers.EndConversation()
}

func (l *logs) addFileCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	l.addMessageToDelete(ctx.EffectiveMessage)

	delete(l.data, ctx.EffectiveChat.Id)
	l.deleteMessages(b, ctx.EffectiveChat.Id)
	return handlers.EndConversation()
}
