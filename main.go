package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dotenv "github.com/joho/godotenv"

	"github.com/pookmaster21/ConnectIM/telegrambot"
	"github.com/pookmaster21/ConnectIM/types"
)

var logger = new(types.Logger)

func main() {
	e := dotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}

	logger.Init()

	ctx, cancelCtx := context.WithCancel(context.TODO())
	defer cancelCtx()

	msgChan := make(chan types.Message)

	TELEGRAM_TOKEN := os.Getenv("TELEGRAM_TOKEN")

	go handle_msgs(ctx, msgChan)
	go telegram.RunTelegramBot(ctx, logger, msgChan, TELEGRAM_TOKEN)

	input := ""
	for input != "exit" {
		fmt.Scanln(&input)
	}
}

func handle_msgs(ctx context.Context, msgChan <-chan types.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			logger.Info("Got: %s[%s] From: %s", msg.To, msg.Msg, msg.From)
		}
	}
}
