package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/artem-shestakov/telegram-budget/internal/models"
)

var createIncomeBtn = []gotgbot.InlineKeyboardButton{
	{
		Text:         "➕ Добавить новый",
		CallbackData: "_create_income",
	},
}

func (tgb *TgBot) createIncomeInfo(b *gotgbot.Bot, ctx *ext.Context) error {
	b.SendMessage(
		ctx.EffectiveChat.Id,
		`
		Чтобы создать источник дохода отпраьте сообщение
		следующего формата:
		название план
		Примеры:
		💳 Зарплата 10000
		💰 Фриланс

		Для отмены используйте команду /cancel
		`,
		nil,
	)
	return handlers.NextConversationState("income_creating")
}

func (tgb *TgBot) createIncome(b *gotgbot.Bot, ctx *ext.Context) error {
	var income models.Income

	// get world from message
	income_slice := strings.Split(ctx.Message.Text, " ")

	// if worlds more then one
	if len(income_slice) > 1 {
		// Check last word. Is it sute for income's plan or it's part of title
		if plan, err := strconv.ParseFloat(income_slice[len(income_slice)-1], 32); err != nil {
			// last word is not number
			income.Plan = 0
			income.Title = strings.Join(income_slice, " ")
		} else {
			// last word is number
			income.Plan = float32(plan)
			income.Title = strings.Join(income_slice[:len(income_slice)-1], " ")
		}
		// received one word (title)
	} else {
		income.Plan = 0
		income.Title = strings.Join(income_slice, " ")
	}

	income.BudgetId = ctx.EffectiveChat.Id

	incomeId, err := tgb.service.Income.Create(income)
	if err != nil {
		tgb.logger.Errorf("Can't create income for %d. %s", ctx.EffectiveChat.Id, err.Error())
		return handlers.NextConversationState("income_creating")
	}
	tgb.logger.Infof("Income %d created in budget %d", incomeId, ctx.EffectiveChat.Id)
	return handlers.EndConversation()
}

func (tgb *TgBot) getIncomes(b *gotgbot.Bot, ctx *ext.Context) error {
	incomes, err := tgb.service.Income.GetAll(ctx.EffectiveChat.Id)
	if err != nil {
		tgb.logger.Errorln(err.Error())
		return nil
	}
	sendIncomes(incomes, b, ctx, "Ваши источники дохода")

	return nil
}

func (tgb *TgBot) startTopUp(b *gotgbot.Bot, ctx *ext.Context) error {
	incomes, err := tgb.service.Income.GetAll(ctx.EffectiveChat.Id)
	if err != nil {
		tgb.logger.Errorln(err.Error())
		return nil
	}
	sendIncomes(incomes, b, ctx, "Куда отнести доход?")
	return handlers.NextConversationState("topup_select_income")
}

func incomesBtn(incomes []models.Income) [][]gotgbot.InlineKeyboardButton {
	var income_btns [][]gotgbot.InlineKeyboardButton
	for _, income := range incomes {
		income_btn := []gotgbot.InlineKeyboardButton{
			{
				Text:         income.Title,
				CallbackData: fmt.Sprintf("_income_%d", income.ID),
			},
		}
		income_btns = append(income_btns, income_btn)
	}

	return append(income_btns, createIncomeBtn)
}

func sendIncomes(incomes []models.Income, b *gotgbot.Bot, ctx *ext.Context, msg string) error {
	income_btns := incomesBtn(incomes)
	b.SendMessage(ctx.EffectiveChat.Id, msg, &gotgbot.SendMessageOpts{
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: income_btns,
		},
		ParseMode: "html",
	})
	return nil
}
