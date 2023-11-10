package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"log"
	"reflect"
	"time"
)

func StartBot() {
	// Создаем бота с токеном, полученным от BotFather
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Fatal(err)
	}

	// Включаем режим отладки, чтобы видеть логи
	bot.Debug = true

	schedule := make(map[string]string, 4)

	// Задаем ID канала, в который хотим отправлять сообщения
	// Вместо 123456789 нужно подставить ваш ID канала
	channelID := int64(123456789)

	loc, err := time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		log.Fatal(err)
	}
	location := cron.WithLocation(loc)
	cronJob := cron.New(location)

	_, err = cronJob.AddFunc("@every 10m", func() {
		now := time.Now().In(loc)
		weekDay := now.Weekday()
		if weekDay == time.Saturday || weekDay == time.Sunday {
			return
		}
		if now.Hour() < 20 || now.Hour() > 21 {
			fmt.Println("Sleeping")
			return
		}
		links, err := FindLinks("https://sustec.ru/raspisanie-ptk-ul-gagarina-7/", "mtli_doc")
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
			schedule = make(map[string]string, 4)
			schedule = links
		}
	})

	if err != nil {
		log.Fatal(err)
	}

	cronJob.Start()

	defer cronJob.Stop()

	select {}
}
