package event

import (
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"tgBotProject/database/storage"
)

const (
	todoMessage1   = "please write a task you want to do"
	todoMessage2   = "your task added to my database"
	deleteMessage1 = "please choose a note from the list,input a number of chosen one"
	changeRequest  = "incorrect answer,please choose a number"
)

type BotListener struct {
	lastBotMsg   string
	UserDatabase storage.UserData
	TaskDatabase storage.Tasks
}

func (l *BotListener) Reply(userMsg *tgbotapi.Message, newBotMsg tgbotapi.MessageConfig, db *sql.DB) tgbotapi.MessageConfig {
	noteList := l.showList(db)
	switch userMsg.Text {
	case "/todo":
		newBotMsg.Text = todoMessage1
		l.lastBotMsg = newBotMsg.Text
		err := storage.Save(&l.UserDatabase, context.Background(), db)
		if err != nil {
			panic(err)
		}
	case "/list":
		if len(noteList) == 0 {
			newBotMsg.Text = "you have no any saved notes already"
		} else {
			newBotMsg.Text = noteList
		}
	case "/delete":
		if len(noteList) == 0 {
			newBotMsg.Text = "you have no any saved notes already"
		} else {
			newBotMsg.Text = deleteMessage1 + "\n" + noteList
			l.lastBotMsg = deleteMessage1
		}
	default:
		if l.lastBotMsg == todoMessage1 {
			newBotMsg.Text = todoMessage2
			l.lastBotMsg = newBotMsg.Text
			l.TaskDatabase.TaskText = userMsg.Text
			err := storage.Save(&l.TaskDatabase, context.Background(), db)
			if err != nil {
				panic(err)
			}
		} else if l.lastBotMsg == deleteMessage1 {
			chosenNumber, err := strconv.Atoi(userMsg.Text)
			if err != nil {
				newBotMsg.Text = changeRequest
				break
			}
			ok, err := l.TaskDatabase.Delete(context.Background(), db, chosenNumber)
			if err != nil {
				panic(err)
			} else if !ok {
				newBotMsg.Text = "No such number in the list"
			} else {
				newBotMsg.Text = "Note successfully deleted"
				l.lastBotMsg = newBotMsg.Text
			}
		} else {
			newBotMsg.Text = userMsg.Text
			l.lastBotMsg = newBotMsg.Text
		}
	}
	return newBotMsg
}

func (l *BotListener) showList(db *sql.DB) string {
	_, taskList, err := storage.Get(&l.TaskDatabase, context.Background(), db)
	answer := ""
	if err != nil {
		panic(err)
	}
	for ind, task := range taskList {
		num := strconv.Itoa(ind + 1)
		answer += num + ". " + task + "\n"
	}
	return answer
}
