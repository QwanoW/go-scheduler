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
	startTime = 19
	endTime   = 22

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
	// prepare db
	scheduleMsgIds, err := GetMsgIds()
	if err != nil {
		log.Fatal(err)
	}

	schedule, err := GetSchedule()
	if err != nil {
		log.Fatal(err)
	}

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

	// init cron
	location := cron.WithLocation(loc)
	cronJob := cron.New(location)

	// add function to cron
	_, err = cronJob.AddFunc(cronInterval, func() {

		now := time.Now().In(loc)
		weekDay := now.Weekday()

		if weekDay == time.Saturday || weekDay == time.Sunday || now.Hour() < startTime || now.Hour() > endTime {
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
			if len(scheduleMsgIds) > 0 {
				deletePreviousSchedule(bot, channelID, scheduleMsgIds)
			}
			scheduleMsgIds = []int{}
			for k, v := range links {
				msg := tgbotapi.NewMessage(channelID, fmt.Sprintf("<a href='%v'>%v</a>", v, k))
				scheduleMsgIds = append(scheduleMsgIds, msg.ReplyToMessageID)
				msg.ParseMode = tgbotapi.ModeHTML
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
			// save schedule
			err := SaveSchedule(links)
			if err != nil {
				return
			}
			schedule = make(map[string]string, 4)
			schedule = links

			err = SaveMsgIds(scheduleMsgIds)
			if err != nil {
				return
			}
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	cronJob.Start()

	defer cronJob.Stop()

	select {}
}

func deletePreviousSchedule(bot *tgbotapi.BotAPI, channelID int64, msgIds []int) {
	for _, messageID := range msgIds {
		msgToDelete := tgbotapi.NewDeleteMessage(channelID, messageID)
		_, err := bot.Request(msgToDelete)
		if err != nil {
			log.Println(err)
		}
	}
}
