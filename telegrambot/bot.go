package telegram

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/pookmaster21/ConnectIM/types"
)

var commands = [2]string{"start", "help"}

func RunTelegramBot(
	ctx context.Context,
	logger *types.Logger,
	msgChan chan types.Message,
	TOKEN string,
) {
	bot, err := initBot(TOKEN)
	if err != nil {
		logger.Error("Error creating bot")
	}

	logger.Info("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		select {
		case <-ctx.Done():
			break
		default:
			msg := handleUpdates(&update, msgChan)
			if msg.Text == "" {
				continue
			}

			_, err := bot.Send(msg)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}

	logger.Info("Closing the Telegram bot")
}

func initBot(TOKEN string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		return nil, err
	}

	bot.Debug = false

	return bot, nil
}

func handleUpdates(
	update *tgbotapi.Update,
	msgChan chan types.Message,
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
				msg.Text += "\t/" + command + "\n"
			}
		default:
			msg.Text = "I don't know that command"
		}

		return
	}

	args := strings.Split(update.Message.Text, ":")
	if len(args) >= 2 {
		msgChan <- types.Message{
			Msg:  strings.Join(args[1:], ":"),
			From: update.SentFrom().String(),
			To:   args[0],
		}
	} else {
		msg.Text = "You need to specify 2 args: User:message"
	}

	return
}
