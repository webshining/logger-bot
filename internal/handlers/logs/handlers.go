package logs

import (
	"bot/internal/database/models"
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"go.uber.org/zap"
)

func (l *logs) files(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.EffectiveMessage.Reply(b, "Select logs to display:", &gotgbot.SendMessageOpts{ReplyMarkup: l.logsMarkup()})
	return nil
}

func (l *logs) filesCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.CallbackQuery.Answer(b, nil)

	data := strings.Split(ctx.CallbackQuery.Data, ":")
	var file models.File
	l.database.First(&file, "id = ?", data[1])

	fileContent, _ := l.readFile(file.Path, 0)
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("`[` %s `]`\n\n```console\n%s\n```", file.Name, fileContent), &gotgbot.SendMessageOpts{ParseMode: "Markdown", ReplyMarkup: l.logMarkup(file.ID)})
	if err != nil {
		l.logger.Error("filed to send message", zap.Error(err))
	}
	return nil
}

func (l *logs) fileCallback(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.CallbackQuery.Answer(b, nil)

	data := strings.Split(ctx.CallbackQuery.Data, ":")
	var file models.File
	if result := l.database.Find(&file, "id = ?", data[2]); result.Error != nil {
		ctx.EffectiveMessage.Delete(b, nil)
		return nil
	}

	switch data[1] {
	case "delete":
		if cancel, ok := l.realtimeLogs[file.ID]; ok {
			cancel()
			delete(l.realtimeLogs, file.ID)
		}
		ctx.EffectiveMessage.Delete(b, nil)
		l.database.Delete(&models.File{}, "id = ?", file.ID)
	case "realtime":
		if cancel, ok := l.realtimeLogs[file.ID]; ok {
			cancel()
			delete(l.realtimeLogs, file.ID)
		} else {
			realtimeContext, cancel := context.WithCancel(context.Background())
			l.realtimeLogs[file.ID] = cancel
			l.watcher(b, ctx.EffectiveMessage, file, realtimeContext, cancel)
		}
	}

	return nil
}

func (l *logs) addFile(b *gotgbot.Bot, ctx *ext.Context) error {
	ctx.CallbackQuery.Answer(b, nil)

	message, _ := ctx.EffectiveMessage.Reply(b, "`[` Add new file `]`\n\nEnter display name:", &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
	l.addMessageToDelete(message)

	return handlers.NextConversationState("NAME")
}

func (l *logs) addFileName(b *gotgbot.Bot, ctx *ext.Context) error {
	name := ctx.EffectiveMessage.Text
	l.data["file_name"] = name

	message, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("`[` Add new file %s `]`\n\nEnter file path:", html.EscapeString(name)), &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
	if err != nil {
		l.logger.Error("failed", zap.Error(err))
	}
	l.addMessageToDelete(message)
	l.addMessageToDelete(ctx.EffectiveMessage)

	return handlers.NextConversationState("PATH")
}

func (l *logs) addFilePath(b *gotgbot.Bot, ctx *ext.Context) error {
	name := l.data["file_name"].(string)
	path := ctx.EffectiveMessage.Text
	l.database.Create(&models.File{Name: name, Path: path})
	l.addMessageToDelete(ctx.EffectiveMessage)

	delete(l.data, "file_name")
	l.deleteMessages(b, ctx.EffectiveChat.Id)
	return handlers.EndConversation()
}

func (l *logs) addFileCancel(b *gotgbot.Bot, ctx *ext.Context) error {
	l.addMessageToDelete(ctx.EffectiveMessage)

	delete(l.data, "file_name")
	l.deleteMessages(b, ctx.EffectiveChat.Id)
	return handlers.EndConversation()
}
