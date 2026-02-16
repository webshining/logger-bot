package logs

import (
	"bot/internal/database/models"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

func (l *logs) watcher(b *gotgbot.Bot, chatId int64, dbFile models.File) {
	text := fmt.Sprintf("`[` %s `]`", dbFile.Name)
	message, _ := b.SendMessage(chatId, text, &gotgbot.SendMessageOpts{ParseMode: "Markdown"})

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var offset int64 = 0
	for {
		<-ticker.C

		content, newOffset := l.readFile(dbFile.Path, offset)
		offset = newOffset

		if _, _, err := b.EditMessageText(text+fmt.Sprintf("\n\n```console\n%s\n```", content), &gotgbot.EditMessageTextOpts{ChatId: chatId, MessageId: message.MessageId, ParseMode: "Markdown"}); err != nil {
			// l.logger.Error("failed to edit message", zap.Error(err))
			// return
		}
	}
}

func (l *logs) readFile(path string, offset int64) (string, int64) {
	file, err := os.Open(path)
	if err != nil {
		l.logger.Error("failed to open file", zap.String("path", path), zap.Error(err))
		return "", 0
	}
	defer file.Close()

	fileSize, _ := file.Seek(0, 2)
	if fileSize < offset {
		offset = 0
	}
	if _, err = file.Seek(offset, 0); err != nil {
		l.logger.Error("failed to seek cursor", zap.String("path", path), zap.Error(err))
		return "", 0
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
	return strings.Join(lines, ""), offset
}
