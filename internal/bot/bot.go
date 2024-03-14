package bot

import (
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/chatmember"
	"github.com/artem-shestakov/telegram-budget/internal/service"
	"github.com/sirupsen/logrus"
)

type TgBot struct {
	token   string
	updater *ext.Updater
	service *service.Service
	logger  *logrus.Logger
}

func NewTgBot(token string, service *service.Service, logger *logrus.Logger) *TgBot {
	return &TgBot{
		token:   token,
		service: service,
		logger:  logger,
	}
}

func (b *TgBot) Run() {
	bot, err := b.initBot()
	if err != nil {
		panic("Fail to create bot: " + err.Error())
	}
	b.initHandlers()
	b.startPolling(bot)
}

func (b *TgBot) initBot() (*gotgbot.Bot, error) {
	bot, err := gotgbot.NewBot(b.token, nil)
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func (b *TgBot) initHandlers() {

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
	})
	b.updater = ext.NewUpdater(dispatcher, nil)

	// Handlers
	dispatcher.AddHandler(handlers.NewMyChatMember(chatmember.Group, b.CreateBudget))
}

func (b *TgBot) startPolling(bot *gotgbot.Bot) {
	err := b.updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	b.logger.Infof("%s has been started...\n", bot.User.Username)

	b.updater.Idle()
}
