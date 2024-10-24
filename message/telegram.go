package message

import (
	"fmt"
	"log"
	"time"

	"github.com/48club/rpc-watchdog/config"
	"github.com/48club/rpc-watchdog/types"
	tele "gopkg.in/telebot.v4"
)

var (
	bot     *tele.Bot
	notifys types.Notifys
)

func init() {
	b, err := tele.NewBot(tele.Settings{Token: config.Config.Token})
	if err != nil {
		panic(err)
	}
	bot = b
	notifys = types.Notifys{}
}

func Notify(c types.Chan) {
	notify, ok := notifys[c.Rpc]
	tt := time.Now()
	if ok {
		if tt.Sub(notify) < config.Config.NotifyInterval*time.Second {
			return
		}
	}

	_, err := bot.Send(
		&tele.User{ID: config.Config.ChatID},
		fmt.Sprintf("[%s]: %s", c.Rpc, longErrChcek(c.Err)),
		&tele.SendOptions{},
	)
	if err != nil {
		log.Printf("Failed to send message: %s", err)
	}
	notifys[c.Rpc] = tt
}

func longErrChcek(err error) (txt string) {
	txt = err.Error()
	if len(txt) > 150 {
		return txt[:150]
	}
	return txt
}
