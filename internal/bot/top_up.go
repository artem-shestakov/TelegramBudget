package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/artem-shestakov/telegram-budget/internal/models"
)

func (tgb *TgBot) startTopUp(b *gotgbot.Bot, ctx *ext.Context) error {
	incomes, err := tgb.service.Income.GetAll(ctx.EffectiveChat.Id)
	if err != nil {
		tgb.logger.Errorln(err.Error())
		return nil
	}
	ctx.Data["msg"] = ctx.EffectiveMessage.Text
	sendIncomes(incomes, b, ctx, "Куда записать поступление?")

	return handlers.NextConversationState("topup_select_income")
}

func (tgb *TgBot) topUp(b *gotgbot.Bot, ctx *ext.Context) error {
	query := ctx.CallbackQuery

	// get data from query
	r, _ := regexp.Compile(`^_income_id(\d+)_name(.*)_msg\+(\d+)\s+(.*)$`)
	queryData := r.FindStringSubmatch(query.Data)

	// Get incomeID and income's amount
	incomeId, _ := strconv.Atoi(queryData[1])
	// incomeTitle, _ := queryData[2]
	topUpAmount, err := strconv.ParseFloat(queryData[3], 64)

	if err != nil {
		tgb.logger.Errorln(err.Error())
		return nil
	}
	// Create top up transaction
	topUpId, err := tgb.service.Income.TopUp(models.TopUp{
		Amount:      topUpAmount,
		Date:        time.Now().Format(time.DateOnly),
		Description: queryData[3],
		IncomeId:    incomeId,
	})
	if err != nil {
		tgb.logger.Errorln(err.Error())
		return handlers.NextConversationState("topup_select_income")
	}

	tgb.logger.Infof("top up %d for income %d created", topUpId, incomeId)
	b.SendMessage(
		ctx.EffectiveChat.Id,
		fmt.Sprintf("<b>%s</b> заработал %.1f (%s)", ctx.CallbackQuery.From.FirstName, topUpAmount, queryData[2]),
		&gotgbot.SendMessageOpts{
			ParseMode: "html",
		})
	ctx.CallbackQuery.Answer(b, nil)
	return handlers.EndConversation()
}
