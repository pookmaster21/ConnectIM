package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/pookmaster21/ConnectIM/db"
	. "github.com/pookmaster21/ConnectIM/types"
)

var commands = [3]string{"/start", "/help", "<username>:<message>"}

var (
	bot    *tgbotapi.BotAPI
	logger *Logger
)

func Run(
	ctx context.Context,
	wg *sync.WaitGroup,
	recieveMsgChan *chan *Message,
	sendingMsgChan *chan *Message,
	TOKEN string,
) {
	logger = NewLogger()

	defer func() {
		logger.Info("Closed the Telegram bot")
		wg.Done()
	}()

	err := initBot(TOKEN)
	if err != nil {
		logger.Error("Error creating bot")
		return
	}

	logger.Info("Telegram bot connected: %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	wg.Add(1)
	go SendMessages(ctx, wg, recieveMsgChan)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			msg := handleUpdates(ctx, &update, sendingMsgChan)
			if msg.Text == "" {
				continue
			}

			_, err := bot.Send(msg)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func initBot(TOKEN string) error {
	var err error

	bot, err = tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		return err
	}

	bot.Debug = false

	return nil
}

func handleUpdates(ctx context.Context,
	update *tgbotapi.Update,
	msgChan *chan *Message,
) (msg tgbotapi.MessageConfig) {
	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if update.Message == nil { // ignore any non-Message updates
		return
	}

	if update.Message.IsCommand() {
		// Extract the command from the Message.
		switch update.Message.Command() {
		case "start":
			fallthrough
		case "help":
			msg.Text = "I understand:\n"
			for _, command := range commands {
				msg.Text += "\t" + command + "\n"
			}
		default:
			msg.Text = "I don't know that command"
		}

		return
	}

	args := strings.Split(update.Message.Text, ":")
	if len(args) >= 2 {
		userFrom := db.DB.FindUser(
			ctx,
			[]string{"telegram"},
			[]any{fmt.Sprint(update.FromChat().ID)},
		)
		if userFrom == nil {
			msg.Text = "You don't have an account"
			return
		}

		userTo := db.DB.FindUser(ctx, []string{"username"}, []any{args[0]})
		if userTo == nil {
			msg.Text = "There is no such user " + args[0]
			return
		}

		*msgChan <- &Message{
			Msg:  strings.Join(args[1:], ":"),
			From: userFrom,
			To:   userTo,
		}
	} else {
		msg.Text = "You need to specify 2 args: <user>:<message>"
	}

	return
}

func SendMessages(ctx context.Context, wg *sync.WaitGroup, msgChan *chan *Message) {
	defer func() {
		logger.Info("Closed telegram bot messages reciever")
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-*msgChan:
			id, err := strconv.ParseInt(msg.To.Telegram, 10, 64)
			if err != nil {
				logger.Error(err.Error())
			}

			_, err = bot.Send(tgbotapi.NewMessage(id, msg.From.Username+":"+msg.Msg))
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
