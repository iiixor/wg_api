package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func initBot() {
	// Канал для сообщений
	messageChannel = make(chan Message)

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		lg.Fatalf("Ошибка при создании бота: %v", err)
	}
	lg.Printf("Авторизован под именем: %s", bot.Self.UserName)

	// Запуск бота в отдельной горутине
	go runBot(bot)
}

func runBot(bot *tgbotapi.BotAPI) {
	for {
		select {
		case msg := <-messageChannel:
			// Создаем сообщение для отправки
			tgMessage := tgbotapi.NewMessage(msg.ChatID, msg.Content)

			tgMessage.ParseMode = "MarkdownV2"
			// Отправляем сообщение в Telegram
			_, err := bot.Send(tgMessage)
			if err != nil {
				lg.Printf("Ошибка при отправке сообщения: %v", err)
			} else {
				lg.Printf("Сообщение отправлено в чат %d", msg.ChatID)
			}
		}
	}
}

func sendMessage(chatID int64, content string) {
	messageChannel <- Message{
		ChatID:  chatID,
		Content: content,
	}
}
