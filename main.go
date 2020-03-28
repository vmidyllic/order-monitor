package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"order-monitor/models"
	"strings"
	"time"
)

type Config struct {
	BotAPIKey  string `mapstructure:"bot_api_key"`
	TargetURL  string `mapstructure:"target_url"`
	ChatID     int64  `mapstructure:"chat_id"`
	MessageURL string `mapstructure:"message_url"`
}

func main() {

	config := ParseConfig()

	FetchPeriodically(FetchDeliverySlots, config.TargetURL, config.BotAPIKey, config.ChatID, config.MessageURL)
}

type fetchFunc func(string) ([]models.DeliverySlot, error)

func FetchDeliverySlots(targetURL string) ([]models.DeliverySlot, error) {

	var slots []models.DeliverySlot

	//resolve claim by link
	resp, err := http.Get(targetURL)

	if err != nil {
		return nil, fmt.Errorf("couldn't fetch delivery slots % s %v ", targetURL, err)
	}
	err = json.NewDecoder(resp.Body).Decode(&slots)
	if err != nil {
		return nil, fmt.Errorf("unexpected model %s ", err)
	}
	return slots, nil

}
func FetchPeriodically(f fetchFunc, targetURL string, key string, chatId int64, messageURL string) {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	for {
		res, err := f(targetURL)
		if err != nil {
			fmt.Errorf(err.Error())
		} else {
			//show open slots
			slotFound := false
			for _, day := range res {
				for _, s := range day.Items {
					if s.IsOpen {
						p := fmt.Sprintf("time : %s free slot : %+v ", time.Now().Format("15:04:05"), s)
						fmt.Println(p)
						slotFound = true
						msg := tgbotapi.NewMessage(chatId, p)
						bot.Send(msg)
						msg2 := tgbotapi.NewMessage(chatId, fmt.Sprintf("go go go : %s", messageURL))
						bot.Send(msg2)
					}
				}
			}
			if !slotFound {
				fmt.Printf("time : %s no free slots found \n", time.Now().Format("15:04:05"))
			}
		}
		time.Sleep(10 * time.Minute)
	}
}
func ParseConfig() Config {
	var config Config

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error while loading config file: %s\n", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("Fatal error while unmarshaling config: %s\n", err))
	}
	return config
}
