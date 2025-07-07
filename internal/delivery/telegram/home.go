package telegram

import (
	"context"
	"fmt"
	listDisplays "github.com/LexBokun/ControlAgent/internal/application/service/query/list-displays"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

func (b *Bot) home(ctx context.Context, u tgbotapi.Update, mode SendMode) error {
	chatID, messageID := extractIDs(u)

	var username string
	if u.Message != nil && u.Message.From != nil {
		username = u.Message.From.UserName
		if username == "" {
			username = fmt.Sprintf("%s %s", u.Message.From.FirstName, u.Message.From.LastName)
		}
	} else if u.CallbackQuery != nil && u.CallbackQuery.From != nil {
		username = u.CallbackQuery.From.UserName
		if username == "" {
			username = fmt.Sprintf("%s %s", u.CallbackQuery.From.FirstName, u.CallbackQuery.From.LastName)
		}
	} else {
		slog.Warn("Update doesn't contain a valid message or callback sender")
		return nil
	}

	result, err := b.list.Handle(ctx, listDisplays.Query{})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	text := fmt.Sprintf("🏠 Главное меню:\n\n👤 @%v\n\nКоличество мониторов: %v.", username, len(result.Displays))
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🖥 Мои мониторы", "to_monitors"),
		),
	)

	switch mode {
	case SendNewMessage:
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = markup
		if _, err := b.bot.Send(msg); err != nil {
			slog.Error(err.Error())
			return err
		}
	case EditMessage:
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ReplyMarkup = &markup
		if _, err := b.bot.Send(edit); err != nil {
			slog.Error(err.Error())
			return err
		}
	}
	return nil
}
