package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func extractIDs(u tgbotapi.Update) (int64, int) {
	if u.CallbackQuery != nil {
		return u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID
	}
	if u.Message != nil {
		return u.Message.Chat.ID, 0
	}
	return 0, 0
}
