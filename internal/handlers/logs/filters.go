package logs

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func noCommands(msg *gotgbot.Message) bool {
	return message.Text(msg) && !message.Command(msg)
}
