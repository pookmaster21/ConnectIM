package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	dotenv "github.com/joho/godotenv"

	"github.com/pookmaster21/ConnectIM/db"
	"github.com/pookmaster21/ConnectIM/telegrambot"
	. "github.com/pookmaster21/ConnectIM/types"
)

var (
	sendMsgChan         = new(chan Message)
	telegramRecvMsgChan = new(chan Message)
	WhatsappRecvMsgChan = new(chan Message)
	discordRecvMsgChan  = new(chan Message)
)

var (
	wg     = new(sync.WaitGroup)
	logger = new(Logger)
)

func main() {
	logger.Init()

	e := dotenv.Load()
	if e != nil {
		logger.Fatal("Error loading .env file")
	}

	ctx, cancelCtx := context.WithCancel(context.TODO())

	db.Init_db(ctx, os.Getenv("MONGO_URI"), logger)
	defer db.DB.Close(ctx)

	TELEGRAM_TOKEN := os.Getenv("TELEGRAM_TOKEN")
	// WHATSAPP_TOKEN := os.Getenv("WHATSAPP_TOKEN")
	// DISCORD_TOKEN := os.Getenv("DISCORD_TOKEN")

	wg.Add(2)
	go handle_msgs(ctx)
	go telegram.RunTelegramBot(ctx, logger, wg, telegramRecvMsgChan, sendMsgChan, TELEGRAM_TOKEN)

	input := ""
	for input != "exit" {
		fmt.Scanln(&input)
	}

	cancelCtx()
	wg.Wait()
}

func handle_msgs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case msg := <-*sendMsgChan:
			logger.Info("Got: %s[%s] From: %s", msg.To, msg.Msg, msg.From)
		}
	}
}
