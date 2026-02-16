package logs

import "github.com/PaulSonOfLars/gotgbot/v2"

func (l *logs) addMessageToDelete(message *gotgbot.Message) {
	l.messagesToDelete[message.Chat.Id] = append(l.messagesToDelete[message.Chat.Id], message.MessageId)
}

func (l *logs) deleteMessages(b *gotgbot.Bot, chatId int64) {
	b.DeleteMessages(chatId, l.messagesToDelete[chatId], nil)
	delete(l.messagesToDelete, chatId)
}
