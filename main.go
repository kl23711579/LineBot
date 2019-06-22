package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	godotenv.Load()

	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")
	userId := os.Getenv("USER_ID")

	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	go Weather(userId, bot)
	go Rubbish(userId, bot)

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
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func Weather(userId string, bot *linebot.Client) {
	fiveMTick := time.Tick(5 * time.Minute)

	local, err := time.LoadLocation("Australia/Brisbane")
	if err != nil {
		fmt.Println(err)
	}

	for t := range fiveMTick {
		tLocal := t.In(local)
		hour, minute, _ := tLocal.Clock()
		if hour == 7 && minute >= 0 && minute <= 5 {
			weatherDetail := GetData()
			fmt.Println(weatherDetail)
			if _, err := bot.PushMessage(userId, linebot.NewTextMessage(weatherDetail)).Do(); err != nil {
				fmt.Println("err    ", err)
				log.Fatal(err)
			}
		} else {
			fmt.Println("Time no OK")
		}
	}
}

func Rubbish(userId string, bot *linebot.Client) {
	hourTick := time.Tick(time.Hour)

	local, err := time.LoadLocation("Australia/Brisbane")  //set time zone of brisbane
	if err != nil {
		fmt.Println(err)
	}

	for t := range hourTick {
		tLocal := t.In(local)  // change to Brisbane time
		weekDay := tLocal.Weekday()
		fmt.Println(weekDay)
		hour := tLocal.Hour()
		if int(weekDay) == 1 && hour >= 7 && hour < 8{
			if _, err := bot.PushMessage(userId, linebot.NewTextMessage("Take out rubbish.")).Do(); err != nil {
				fmt.Println("Rubbish err   ", err)
				log.Fatal(err)
			}
		}
	}
}
