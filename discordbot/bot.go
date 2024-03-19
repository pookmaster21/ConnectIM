package discord

import (
	"context"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/pookmaster21/ConnectIM/db"
	. "github.com/pookmaster21/ConnectIM/types"
)

type discordBot struct {
	session *discordgo.Session
	ID      string
}

var (
	ctx             context.Context
	sendMessageChan *chan *Message
	recvMessageChan *chan *Message
	bot             *discordBot
	logger          *Logger
	wg              *sync.WaitGroup
)

func Run(
	Ctx context.Context,
	WG *sync.WaitGroup,
	RecvMsgChan *chan *Message,
	SendMsgChan *chan *Message,
	token string,
) {
	ctx = Ctx
	sendMessageChan = SendMsgChan
	recvMessageChan = RecvMsgChan
	logger = NewLogger()
	wg = WG

	defer func() {
		<-ctx.Done()
		logger.Info("Closed discord bot")
		wg.Done()
	}()

	bot = new(discordBot)

	var err error
	bot.session, err = discordgo.New("Bot " + token)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	u, err := bot.session.User("@me")
	if err != nil {
		logger.Error("Failed getting current User: %s", err.Error())
		return
	}

	bot.ID = u.ID

	bot.session.AddHandler(handleMessage)

	err = bot.session.Open()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	wg.Add(1)
	go sendMessages()

	logger.Info("Discord bot connected: %s", u.String())
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == bot.ID {
		return
	}

	args := strings.Split(m.Content, ":")
	if len(args) != 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "You need to specify 2 args: <user>:<message>")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	fromUser := db.DB.FindUser(ctx, []string{"discord"}, []any{m.ChannelID})
	if fromUser == nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "You don't have an account")
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	toUser := db.DB.FindUser(ctx, []string{"username"}, []any{args[0]})
	if toUser == nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "There is no such user "+args[0])
		if err != nil {
			logger.Error(err.Error())
		}
		return
	}

	*sendMessageChan <- &Message{
		To:   toUser,
		From: fromUser,
		Msg:  strings.Join(args[1:], ":"),
	}
}

func sendMessages() {
	defer func() {
		logger.Info("Closed discord bot message reciever")
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-*recvMessageChan:
			_, err := bot.session.ChannelMessageSend(msg.To.Discord, msg.From.Username+":"+msg.Msg)
			if err != nil {
				logger.Error(err.Error())
			}
		default:
		}
	}
}
