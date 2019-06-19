package main

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
)

func main() {

	channelSecret := "cc8b2d517c550675def314b00543a098"
	channelToken := "nyKC/WBqxvVrx15ftpm7+GGSlaRY1BQjDh/vV712PL3iAvJ0UCyqE3Nz5MIhuaymAvB+DP8v17IDzAlaUNYJ5CpOTa8ByRlYVtYS5Sxusd2EUvuPSYo7zvndX09RTSKPLqmLYtl91X7JT7cLqRC2YAdB04t89/1O/w1cDnyilFU="
	userId := "U2ccabf42f930672b3187c9e21e3cb52a"

	weatherDetail := GetData()
	fmt.Println(weatherDetail)

	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := bot.PushMessage(userId, linebot.NewTextMessage(weatherDetail)).Do(); err != nil {
		fmt.Println("err    ", err)
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/webhook", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					fmt.Println(message.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}

