package bot

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() {
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		Handle(bot, update)
	}
}
