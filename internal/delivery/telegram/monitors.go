package telegram

import (
	"context"
	"errors"
	"fmt"
	setenable "github.com/LexBokun/ControlAgent/internal/application/service/command/set-enable"
	listDisplays "github.com/LexBokun/ControlAgent/internal/application/service/query/list-displays"
	"github.com/LexBokun/ControlAgent/internal/domain/display"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strings"
)

func (b *Bot) monitors(ctx context.Context, u tgbotapi.Update, mode SendMode) error {
	chatID, messageID := extractIDs(u)
	var text strings.Builder

	text.WriteString("🖥 Ваши мониторы:")

	result, err := b.list.Handle(ctx, listDisplays.Query{})
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	for _, v := range result.Displays {
		enable := "Включен"
		if !v.Enable {
			enable = "Выключен"
		}
		id := strings.TrimPrefix(v.ID, `\\.\`)
		text.WriteString(fmt.Sprintf("\n\nИмя: %v.\nID: /%v.\nСтатус: %v.", v.Name, id, enable))
	}
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "to_home"),
		),
	)

	switch mode {
	case SendNewMessage:
		msg := tgbotapi.NewMessage(chatID, text.String())
		msg.ReplyMarkup = markup
		if _, err := b.bot.Send(msg); err != nil {
			slog.Error(err.Error())
			return err
		}
	case EditMessage:
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text.String())
		edit.ReplyMarkup = &markup
		if _, err := b.bot.Send(edit); err != nil {
			slog.Error(err.Error())
			return err
		}
	}
	return nil
}

func (b *Bot) monitorsInfo(ctx context.Context, u tgbotapi.Update, mode SendMode) error {
	chatID, messageID := extractIDs(u)

	// Получаем название дисплея из команды
	cmdParts := strings.Split(u.Message.Text, "/")
	if len(cmdParts) < 2 {
		return errors.New("команда без имени дисплея")
	}
	displayName := cmdParts[1]

	text := "🖥 Информация о " + displayName + "\n\nСостояние: Включен." // Заглушка

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Переключить режим", "toggle_"+displayName),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "to_home"),
		),
	)

	switch mode {
	case SendNewMessage:
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = markup
		b.bot.Send(msg)
	case EditMessage:
		edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
		edit.ReplyMarkup = &markup
		b.bot.Send(edit)
	}
	return nil
}

func (b *Bot) toggleMonitor(ctx context.Context, u tgbotapi.Update, callback *tgbotapi.CallbackQuery) error {
	result, err := b.list.Handle(ctx, listDisplays.Query{})
	if err != nil {
		slog.Error("list.Handle failed", slog.String("err", err.Error()))
		return err
	}

	if callback.Message == nil {
		slog.Error("callback message is nil")
		return errors.New("callback message is nil")
	}

	callbackParts := strings.Split(callback.Data, "_")
	if len(callbackParts) != 2 {
		slog.Error("неверный callback_data", slog.String("data", callback.Data))
		return errors.New("invalid callback_data")
	}
	displayID := `\\.\` + callbackParts[1]

	// Ищем нужный дисплей
	var display display.Display
	found := false
	for _, d := range result.Displays {
		if d.ID == displayID {
			display = d
			found = true
			break
		}
	}

	if !found {
		slog.Error("дисплей не найден", slog.String("id", displayID))
		b.bot.Request(tgbotapi.NewCallback(callback.ID, "Дисплей не найден"))
		return fmt.Errorf("дисплей %s не найден", displayID)
	}

	newState := !display.Enable

	err = b.setEnable.Handle(ctx, setenable.Command{
		ID:     display.ID,
		Enable: newState,
	})
	if err != nil {
		slog.Error("failed to toggle display", slog.String("err", err.Error()))
		b.bot.Request(tgbotapi.NewCallback(callback.ID, "Ошибка переключения"))
		return err
	}

	// Формируем новое сообщение и клавиатуру
	stateText := "Выкл"
	if newState {
		stateText = "Вкл"
	}

	displayName := strings.TrimPrefix(display.ID, `\\.\`)

	newText := fmt.Sprintf("🖥 Информация о %s\n\nСостояние: %s.", displayName, stateText)

	newKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Переключить Режим", "toggle_"+displayName),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", "to_home"),
		),
	)

	editMsg := tgbotapi.NewEditMessageTextAndMarkup(
		callback.Message.Chat.ID,
		callback.Message.MessageID,
		newText,
		newKeyboard,
	)

	_, err = b.bot.Send(editMsg)
	if err != nil {
		slog.Error("failed to edit message", slog.String("err", err.Error()))
		return err
	}

	return nil
}
