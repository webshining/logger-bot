package logs

import (
	"bot/internal/database/models"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func (l *logs) logsMarkup() gotgbot.InlineKeyboardMarkup {
	var buttons [][]gotgbot.InlineKeyboardButton
	var row []gotgbot.InlineKeyboardButton

	var files []models.File
	l.database.Find(&files)

	for _, file := range files {
		row = append(row, gotgbot.InlineKeyboardButton{Text: file.Name, CallbackData: fmt.Sprintf("files:%d", file.ID)})
		if len(row) == 2 {
			buttons = append(buttons, row)
			row = []gotgbot.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		buttons = append(buttons, row)
	}
	buttons = append(buttons, []gotgbot.InlineKeyboardButton{gotgbot.InlineKeyboardButton{Text: "Add", CallbackData: "files:add"}})

	return gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}
