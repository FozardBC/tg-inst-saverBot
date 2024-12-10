package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"gliza/inst"
)

func main() {
	// Укажите токен вашего Telegram-бота

	botToken := "7234404767:AAElVFDeAFH1SBZGnNqttwCCWKzCoDUQFOI"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации бота: %v", err)
	}

	bot.Debug = true
	log.Printf("Бот %s успешно запущен", bot.Self.UserName)

	// Настраиваем канал обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			link := update.Message.Text

			// Проверяем, является ли сообщение ссылкой на Instagram
			if strings.Contains(link, "instagram.com") {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Обрабатываю вашу ссылку, подождите немного..."))

				// Загружаем контент
				filePath, err := inst.File(link)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при скачивании контента: %v", err)))
					continue
				}

				//Отправляем файл пользователю
				//img := tgbotapi.NewPhotoToChannel("-2483058849", tgbotapi.FilePath(filePath))
				photo := tgbotapi.NewPhoto(-1002483058849, tgbotapi.FilePath(filePath))
				//file := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
				//file.ChatID = update.Message.Chat.ID
				_, sendErr := bot.Send(photo)
				if sendErr != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при отправке файла: %v", sendErr)))
				}
				//if c == 1 {

				// 	//Отправляем файл пользователю
				// 	//img := tgbotapi.NewPhotoToChannel("-2483058849", tgbotapi.FilePath(filePath))
				// 	photo := tgbotapi.NewPhoto(-1002483058849, tgbotapi.FilePath(filePath))
				// 	//file := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FilePath(filePath))
				// 	//file.ChatID = update.Message.Chat.ID
				// 	_, sendErr := bot.Send(photo)
				// 	if sendErr != nil {
				// 		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при отправке файла: %v", sendErr)))
				// 	}

				// 	//Удаляем временный файл
				// 	os.Remove(filePath)
				// }
				// if c == 2 {
				// 	video := tgbotapi.NewVideo(-1002483058849, tgbotapi.FilePath(filePath))
				// 	_, sendErr := bot.Send(video)
				// 	if sendErr != nil {
				// 		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка при отправке файла: %v", sendErr)))
				// 	}
				// }
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, отправьте ссылку на Instagram."))
			}
		}
	}
}
