package logs

import (
	"bot/internal/database/models"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"go.uber.org/zap"
)

func (l *logs) addMessageToDelete(message *gotgbot.Message) {
	if _, ok := l.data["messages_to_delete"]; !ok {
		l.data["messages_to_delete"] = make([]int64, 0)
	}
	l.data["messages_to_delete"] = append(l.data["messages_to_delete"].([]int64), message.MessageId)
}

func (l *logs) deleteMessages(b *gotgbot.Bot, chatId int64) {
	messagesToDelete, ok := l.data["messages_to_delete"].([]int64)
	if !ok {
		messagesToDelete = make([]int64, 0)
	}
	b.DeleteMessages(chatId, messagesToDelete, nil)
	delete(l.data, "messages_to_delete")
}

func (l *logs) watcher(b *gotgbot.Bot, message *gotgbot.Message, file models.File, ctx context.Context, cancel context.CancelFunc) {
	text := fmt.Sprintf("`[` %s `]`", file.Name)

	go func() {
		l.logger.Info("watcher started", zap.String("path", file.Path))
		defer l.logger.Info("watcher stopped", zap.String("path", file.Path))

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		var offset int64 = 0
		defer func() {
			content, _ := l.readFile(file.Path, offset)
			message.EditText(b, text+fmt.Sprintf("\n\n```console\n%s\n```", content), &gotgbot.EditMessageTextOpts{ParseMode: "Markdown", ReplyMarkup: l.logMarkup(file.ID)})
		}()

		var prevText string = ""
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				content, newOffset := l.readFile(file.Path, offset)
				offset = newOffset

				text := fmt.Sprintf("%s\n\n```console\n%s\n```", text, content)
				if text == prevText {
					continue
				}
				prevText = text
				if _, _, err := message.EditText(b, text, &gotgbot.EditMessageTextOpts{ParseMode: "Markdown", ReplyMarkup: l.logMarkup(file.ID)}); err != nil {
					if strings.Contains(err.(*gotgbot.TelegramError).Description, "message to edit not found") {
						cancel()
					} else {
						l.logger.Error("err", zap.Error(err))
					}
				}
			}
		}
	}()
}

func (l *logs) readFile(path string, offset int64) (string, int64) {
	file, err := os.Open(path)
	if err != nil {
		return err.Error(), 0
	}
	defer file.Close()

	fileSize, _ := file.Seek(0, 2)
	if fileSize < offset {
		offset = 0
	}
	if _, err = file.Seek(offset, 0); err != nil {
		return err.Error(), 0
	}

	var lines []string
	var offsets []int64
	var lastOffset int64 = offset

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if len(line) > 0 {
				lines = append(lines, string(line))
				offsets = append(offsets, lastOffset)
				lastOffset += int64(len(line))
			}
			break
		}

		lines = append(lines, string(line))

		offsets = append(offsets, lastOffset)
		lastOffset += int64(len(line))

		if len(lines) > 10 {
			lines = lines[1:]
			offsets = offsets[1:]
		}
	}

	if len(offsets) > 0 {
		offset = offsets[0]
	}
	return strings.TrimSpace(strings.Join(lines, "")), offset
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
