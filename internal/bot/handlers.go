package bot

import (
	"fmt"
	"go_requests/internal/aviasales"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		handleCallback(bot, update.CallbackQuery)
		return
	}
	if update.Message == nil {
		return
	}
	chatId := update.Message.Chat.ID
	text := update.Message.Text

	if text == "/start" {
		states[chatId] = &UserState{Step: StepFrom}
		msg := tgbotapi.NewMessage(chatId, "Откуда летим? Код города, например MOW")
		msg.ReplyMarkup = replyKeyboard()
		bot.Send(msg)
		return
	}

	state, ok := states[chatId]
	if !ok {
		bot.Send(tgbotapi.NewMessage(chatId, "Напиши /start"))
		return
	}

	switch state.Step {
	case StepFrom:
		state.From = text
		state.Step = StepTo
		msg := tgbotapi.NewMessage(
			chatId, fmt.Sprintf("Город вылета: %s. Теперь введи город прибытия, например HKT", text))
		msg.ReplyMarkup = replyKeyboard()
		bot.Send(msg)
	case StepTo:
		state.To = text
		state.Step = StepDeparture
		msg := tgbotapi.NewMessage(
			chatId, fmt.Sprintf("Выбран город прилета: %s. Теперь введи месяц отбытия, например 2026-05", text))
		msg.ReplyMarkup = replyKeyboard()
		bot.Send(msg)
	case StepDeparture:
		state.Departure = text
		state.Step = StepReturn
		msg := tgbotapi.NewMessage(
			chatId, fmt.Sprintf("Выбрана дата вылета: %s. Выбери дату прилета, например 2026-06", text))
		msg.ReplyMarkup = replyKeyboard()
		bot.Send(msg)
	case StepReturn:
		state.Return = text
		state.Step = StepTransfers

		// сначала восстанавливаем ReplyKeyboard
		confirmMsg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Дата возрата: %s", text))
		confirmMsg.ReplyMarkup = replyKeyboard()
		bot.Send(confirmMsg)

		// потом отдельное сообщение с inline-кнопками
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Только прямые", "direct"),
				tgbotapi.NewInlineKeyboardButtonData("Любые", "any"),
			),
		)
		msg := tgbotapi.NewMessage(chatId, "Пересадки?")
		msg.ReplyMarkup = keyboard
		bot.Send(msg)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery) {
	chatId := cb.Message.Chat.ID
	state := states[chatId]

	direct := cb.Data == "direct"
	result, err := aviasales.DoSearch(state.From, state.To, state.Departure, state.Return, direct)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatId, "Ошибка во время запроса"))
		return
	}
	msg := tgbotapi.NewMessage(chatId, result)
	msg.ParseMode = "HTML"
	bot.Send(msg)
	delete(states, chatId)
}

func replyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заново"),
		),
	)
}
