package telegram

import (
	"context"
	setenable "github.com/LexBokun/ControlAgent/internal/application/service/command/set-enable"
	listdisplays "github.com/LexBokun/ControlAgent/internal/application/service/query/list-displays"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SendMode int

const (
	SendNewMessage SendMode = iota
	EditMessage
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	list      *listdisplays.QueryHandler
	setEnable *setenable.CommandHandler
}

func New(
	c Config,
	list *listdisplays.QueryHandler,
	setEnable *setenable.CommandHandler,
) *Bot {
	bot, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	bot.Debug = c.Debug

	return &Bot{
		bot:       bot,
		list:      list,
		setEnable: setEnable,
	}
}

func (b *Bot) Start(ctx context.Context) error {
	slog.Info("Телеграм бот авторизован", "username", b.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.bot.GetUpdatesChan(u)
	go func(u tgbotapi.UpdatesChannel) {
		for {
			select {
			case update := <-u:
				b.route(ctx, update)
			case <-ctx.Done():
				slog.Debug("Телеграм бот завершил работу")
				return
			}
		}
	}(updates)

	return nil
}
