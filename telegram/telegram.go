package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAdapter struct {
	bot    *tgbotapi.BotAPI
	chatID string
}

func NewTelegramAdapter(botToken string, chatID string) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &TelegramAdapter{bot: bot, chatID: chatID}, nil
}
