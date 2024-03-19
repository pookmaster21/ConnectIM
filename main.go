package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	dotenv "github.com/joho/godotenv"
	"github.com/pookmaster21/ConnectIM/db"
	discord "github.com/pookmaster21/ConnectIM/discordbot"
	"github.com/pookmaster21/ConnectIM/router"
	"github.com/pookmaster21/ConnectIM/telegrambot"
	. "github.com/pookmaster21/ConnectIM/types"
)

var (
	sendMsgChan         = make(chan *Message)
	telegramRecvMsgChan = make(chan *Message)
	WhatsappRecvMsgChan = make(chan *Message)
	discordRecvMsgChan  = make(chan *Message)
)

var (
	wg     = new(sync.WaitGroup)
	logger *Logger
)

func main() {
	logger = NewLogger()

	e := dotenv.Load()
	if e != nil {
		logger.Error("Error loading .env file")
	}

	ctx, cancelCtx := context.WithCancel(context.TODO())

	db.InitDB(ctx, os.Getenv("MONGO_URI"))

	TELEGRAM_TOKEN := os.Getenv("TELEGRAM_TOKEN")
	// WHATSAPP_TOKEN := os.Getenv("WHATSAPP_TOKEN")
	DISCORD_TOKEN := os.Getenv("DISCORD_TOKEN")

	wg.Add(4)
	go handle_msgs(ctx)
	go telegram.Run(ctx, wg, &telegramRecvMsgChan, &sendMsgChan, TELEGRAM_TOKEN)
	go discord.Run(ctx, wg, &discordRecvMsgChan, &sendMsgChan, DISCORD_TOKEN)
	go router.Run(ctx, wg)

	user := &User{
		Discord:  "1183753592386093156",
		Username: "Lavi",
		Password: "a",
		Telegram: "5587088441",
		Whatsapp: "",
		Prefered: TELEGRAM | DISCORD,
	}

	db.DB.DeleteUser(ctx, []string{"username"}, []any{"Lavi"})
	db.DB.InsertUser(ctx, user)

	input := ""
	for input != "exit" {
		fmt.Scanln(&input)
		db.DB.Close(ctx)
	}

	cancelCtx()
	wg.Wait()
}

func handle_msgs(ctx context.Context) {
	defer func() {
		logger.Info("Closed message middlewhare")
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-sendMsgChan:
			logger.Info("Got: %s[%s] From: %s", msg.To.Username, msg.Msg, msg.From.Username)

			if msg.To.Prefered&TELEGRAM == TELEGRAM {
				telegramRecvMsgChan <- msg
			}

			if msg.To.Prefered&DISCORD == DISCORD {
				discordRecvMsgChan <- msg
			}
		}
	}
}
