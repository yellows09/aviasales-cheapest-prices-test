package main

import (
	"go_requests/internal/bot"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	bot.Init()
}
