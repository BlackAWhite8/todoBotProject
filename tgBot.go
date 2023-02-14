package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"tgBotProject/database/storage"
	"tgBotProject/errors/e"
	"tgBotProject/processor/event"
)

const (
	tokenErr   = "bot token not found or wrong, please check token and try again"
	messageErr = "message sending failed"
)

func main() {
	var (
		uTable   storage.UserData
		tTable   storage.Tasks
		listener event.BotListener
	)

	if err := godotenv.Load(); err != nil {
		e.WrapErr(".env file not found", err)
		log.Fatal()
	}

	token, _ := os.LookupEnv("TOKEN")
	dbPath, _ := os.LookupEnv("DATABASE")
	db, err := sql.Open("mysql", dbPath)
	defer db.Close()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		e.WrapErr(tokenErr, err)
		log.Fatal()
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 10
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			uTable.UserID = update.Message.Chat.ID
			uTable.UserName = update.Message.Chat.UserName
			tTable.TaskID = update.Message.Chat.ID
			listener.UserDatabase = uTable
			listener.TaskDatabase = tTable
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg = listener.Reply(update.Message, msg, db)
			if _, err := bot.Send(msg); err != nil {
				e.WrapErr(messageErr, err)
			}
		}
	}
}
