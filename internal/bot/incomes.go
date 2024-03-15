package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

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
	sendIncomes(incomes, b, ctx, "Ваши источники дохода", createIncomeBtn)

	return nil
}

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
	r, _ := regexp.Compile(`^_income_id(\d+)_msg\+(\d+)\s+(.*)$`)
	queryData := r.FindStringSubmatch(query.Data)

	incomeId, _ := strconv.Atoi(queryData[1])
	amount, err := strconv.ParseFloat(queryData[2], 64)

	if err != nil {
		tgb.logger.Errorln(err.Error())
		return nil
	}
	topUpId, err := tgb.service.Income.TopUp(models.TopUp{
		Amount:      amount,
		Date:        time.Now().Format(time.DateOnly),
		Description: queryData[3],
		IncomeId:    incomeId,
	})
	if err != nil {
		tgb.logger.Errorln(err.Error())
		return handlers.NextConversationState("topup_select_income")
	}

	tgb.logger.Infof("top up %d for income %d created", topUpId, incomeId)
	b.SendMessage(ctx.EffectiveChat.Id, "Пополнили", nil)
	ctx.CallbackQuery.Answer(b, nil)
	return handlers.EndConversation()
}

func incomesBtn(ctx *ext.Context, incomes []models.Income, btns ...[]gotgbot.InlineKeyboardButton) [][]gotgbot.InlineKeyboardButton {
	var income_btns [][]gotgbot.InlineKeyboardButton
	for _, income := range incomes {
		income_btn := []gotgbot.InlineKeyboardButton{
			{
				Text:         income.Title,
				CallbackData: fmt.Sprintf("_income_id%d_msg%s", income.ID, ctx.Data["msg"]),
			},
		}
		income_btns = append(income_btns, income_btn)
	}

	income_btns = append(income_btns, btns...)
	return income_btns
}

func sendIncomes(incomes []models.Income, b *gotgbot.Bot, ctx *ext.Context, msg string, btns ...[]gotgbot.InlineKeyboardButton) error {
	income_btns := incomesBtn(ctx, incomes, btns...)
	b.SendMessage(ctx.EffectiveChat.Id, msg, &gotgbot.SendMessageOpts{
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: income_btns,
		},
		ParseMode: "html",
	})
	return nil
}
