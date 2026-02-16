package filters

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func Admin(id int64) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		if ctx.EffectiveChat.Id == id {
			return nil
		}

		return ext.EndGroups
	}
}
