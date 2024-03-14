package bot

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (tgb *TgBot) CreateBudget(b *gotgbot.Bot, ctx *ext.Context) error {
	id, err := tgb.service.CreateBudget(ctx.EffectiveChat.Id, ctx.EffectiveChat.Title)
	if err != nil {
		tgb.logger.Errorf(err.Error())
		return err
	}
	tgb.logger.Infof("Budget %d was created", id)
	b.SendMessage(ctx.EffectiveChat.Id, "TEST", nil)
	return nil
}
