package handler

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

var ButtonsHandlers = map[disgord.Snowflake]func(*disgord.InteractionCreate){}
var ButtonLock = sync.RWMutex{}

func RegisterBHandler(msg *disgord.Message, callback func(*disgord.InteractionCreate), timeout int) {
	ButtonsHandlers[msg.ID] = callback
	if timeout != 0 {
		time.Sleep(time.Duration(timeout) * time.Second)
		DeleteBHandler(msg)
	}
}

func DeleteBHandler(msg *disgord.Message) {
	ButtonLock.Lock()
	delete(ButtonsHandlers, msg.ID)
	ButtonLock.Unlock()
}

func handleButton(interaction *disgord.InteractionCreate) {
	ButtonLock.RLock()
	if cb, found := ButtonsHandlers[interaction.Message.ID]; found {
		ButtonLock.RUnlock()
		cb(interaction)
		return
	}
	ButtonLock.RUnlock()
}

func Interaction(session disgord.Session, evt *disgord.InteractionCreate) {
	if evt.Type == disgord.InteractionMessageComponent && evt.Member != nil {
		go handleButton(evt)
	}
}