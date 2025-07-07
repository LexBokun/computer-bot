package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strings"
)

// Входная точка всех команд тг бота
func (b *Bot) route(ctx context.Context, update tgbotapi.Update) {
	if update.Message != nil {
		switch update.Message.Command() {
		case "home":
			err := b.home(ctx, update, SendNewMessage)
			if err != nil {
				slog.Error(err.Error())
			}
		case "monitors":
			err := b.monitors(ctx, update, SendNewMessage)
			if err != nil {
				slog.Error(err.Error())
			}
		default:
			if strings.HasPrefix(update.Message.Text, "/DISPLAY") {
				b.monitorsInfo(ctx, update, SendNewMessage)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⁉️ Что-то пошло не так, походу у разработчика кривые руки.")
				b.bot.Send(msg)
			}
		}
		return
	}

	if update.CallbackQuery != nil {
		// Уведомляем Telegram, что мы обработали callback
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := b.bot.Request(callback); err != nil {
			slog.Error("callback request failed:", err)
		}

		switch update.CallbackQuery.Data {
		case "to_home":
			err := b.home(ctx, update, EditMessage)
			if err != nil {
				slog.Error(err.Error())
			}
		case "to_monitors":
			err := b.monitors(ctx, update, EditMessage)
			if err != nil {
				slog.Error(err.Error())
			}

		default:
			if strings.HasPrefix(update.CallbackQuery.Data, "toggle_DISPLAY") {
				//displayID := strings.TrimPrefix(update.CallbackQuery.Data, "toggle_")
				err := b.toggleMonitor(ctx, update, update.CallbackQuery) // твоя логика переключения
				if err != nil {
					slog.Error("failed to toggle display:", err)
				}
			} else if update.CallbackQuery.Message != nil && strings.HasPrefix(update.CallbackQuery.Message.Text, "/DISPLAY") {
				b.monitorsInfo(ctx, update, EditMessage)
			} else {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				b.bot.Send(msg)
			}
		}
	}
}
