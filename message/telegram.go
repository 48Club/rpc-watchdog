package message

import (
	"fmt"
	"log"

	"github.com/48club/rpc-watchdog/config"
	"github.com/48club/rpc-watchdog/types"
	tele "gopkg.in/telebot.v4"
)

var bot *tele.Bot

func init() {
	b, err := tele.NewBot(tele.Settings{Token: config.Config.Token})
	if err != nil {
		panic(err)
	}
	bot = b
}

func Notify(c types.Chan) {
	_, err := bot.Send(
		&tele.User{ID: config.Config.ChatID},
		fmt.Sprintf("[%s]: %s", c.Rpc, c.Err),
		&tele.SendOptions{},
	)
	if err != nil {
		log.Printf("Failed to send message: %s", err)
	}
}
