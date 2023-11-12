package internal

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"scheduler/config"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

var (
	// execution time
	startTime = 13
	endTime   = 16

	// api
	botApi    = os.Getenv("TG_BOT_API_KEY")
	channelID = config.GetenvInt64("TG_CHANNEL_ID")

	// cron props
	cronInterval = "@every 10m"

	// parser vars
	parseTarget = "mtli_doc"
	source      = "https://sustec.ru/raspisanie-ptk-ul-gagarina-7/"
)

func StartBot() {
	// Создаем бота с токеном, полученным от BotFather
	bot, err := tgbotapi.NewBotAPI(botApi)
	if err != nil {
		log.Fatal(err)
	}

	// Включаем режим отладки, чтобы видеть логи
	bot.Debug = true

	// get location
	loc, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		log.Fatal(err)
	}

	// scheduler vars
	schedule := make(map[string]string, 4)
	isUpdated := false
	updatedAt := time.Now().In(loc).AddDate(0, 0, -1).Day()

	// init cron
	location := cron.WithLocation(loc)
	cronJob := cron.New(location)

	// add function to cron
	_, err = cronJob.AddFunc(cronInterval, func() {

		now := time.Now().In(loc)
		weekDay := now.Weekday()
		currentDay := now.Day()

		if weekDay == time.Saturday || weekDay == time.Sunday || now.Hour() < startTime || now.Hour() > endTime || (isUpdated && currentDay != updatedAt) {
			fmt.Println("Sleeping")
			return
		}

		links, err := FindLinks(source, parseTarget)
		if err != nil {
			log.Println(err)
		}

		if reflect.DeepEqual(links, schedule) {
			log.Println("No new links")
			return
		} else {
			for k, v := range links {
				msg := tgbotapi.NewMessage(channelID, fmt.Sprintf("<a href='%v'>%v</a>", v, k))
				msg.ParseMode = tgbotapi.ModeHTML
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
			// save schedule
			schedule = make(map[string]string, 4)
			schedule = links

			// save updated
			updatedAt = currentDay
			isUpdated = true
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	cronJob.Start()

	defer cronJob.Stop()

	select {}
}
