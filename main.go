package main

import (
	"EDFCBank/handlers"
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
	handlers.RegisterResourceRoutes(r)

	// connect, and stay connected until a system interrupt takes place
	defer client.StayConnectedUntilInterrupted(context.Background())
	client.Ready(onReady)
}

func onReady() {
	//if count, err := db.Resources.Count(); err != nil {
	//	fmt.Println("no known resources")
	//} else {
	//	fmt.Printf("%d known resources\n", count)
	//
	//	elems, err := db.Resources.GetAll()
	//	utils.LogOnError(err)
	//	for _, elem := range elems {
	//		fmt.Println("- " + elem.Name)
	//	}
	//}

	fmt.Println("READY NOW!")
}
