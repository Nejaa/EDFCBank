package main

import (
	"EDFCBank/handlers"
	"EDFCBank/params"
	"EDFCBank/router"
	"context"
	"fmt"
	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
		Logger:   logrus.StandardLogger(),
	})

	botPrefix := os.Getenv("BOT_PREFIX")

	r := router.NewRouter(client, botPrefix)
	r.DeleteOnAnswer(*params.DeleteOnAnswer)
	handlers.RegisterResourceRoutes(r)
	handlers.RegisterBankRoutes(r)

	// connect, and stay connected until a system interrupt takes place
	defer func() { _ = client.StayConnectedUntilInterrupted(context.Background()) }()
	client.Ready(onReady)
}

func onReady() {
	fmt.Println("READY NOW!")
}
