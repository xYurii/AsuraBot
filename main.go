/*
    __   ____  _  _  ____   __
   / _\ / ___)/ )( \(  _ \ / _\
  /    \\___ \) \/ ( )   //    \
  \_/\_/(____/\____/(__\_)\_/\_/

  This is the source code of the bot
  "Asura" made in Golang and manteined by
  Chiyoku and Acnologia
  Copyright 2020

*/

package main

import (
	_ "asura/src/commands" // Initialize all commands and put them into an array
	"asura/src/database"
	"asura/src/handler"
	"asura/src/telemetry"
	"context"
	"fmt"
	"time"
	"os"
	"strconv"
	"github.com/andersfylling/disgord"
	"github.com/joho/godotenv"
)

func onReady(session disgord.Session, evt *disgord.Ready) {
	telemetry.Info(fmt.Sprintf("%s Started", evt.User.Username), map[string]string{
		"eventType": "ready",
	})
	go telemetry.MetricUpdate(handler.Client)

}

func onGuildDelete(session disgord.Session, evt *disgord.GuildDelete) {
	return
	guild ,err := handler.Client.GetGuild(context.Background(),evt.UnavailableGuild.ID)
	if err != nil {
		fmt.Println(err)
	}
	telemetry.Warn(fmt.Sprintf("Leaved from %s", guild.Name), map[string]string{
		"id": strconv.FormatUint(uint64(evt.UnavailableGuild.ID), 10),
		"eventType": "leave",
	})
	
}

func onGuildCreate(session disgord.Session, evt *disgord.GuildCreate) {
	return
	guild := evt.Guild
	if guild.JoinedAt == nil{
		telemetry.Warn(fmt.Sprintf("Joined in  %s", guild.Name), map[string]string{
			"id":  strconv.FormatUint(uint64(guild.ID), 10),
			"eventType": "join",
		})
		return
	}
	if (20 * time.Second) > (time.Since(guild.JoinedAt.Time) * time.Second){
		telemetry.Warn(fmt.Sprintf("Joined in  %s", guild.Name), map[string]string{
			"id":  strconv.FormatUint(uint64(guild.ID), 10),
			"eventType": "join",
		})
	}
}

func main() {

	// If it's not in production so it's good to read a ".env" file
	if os.Getenv("PRODUCTION") == "" {
		err := godotenv.Load()
		if err != nil {
			panic("Cannot read the motherfucking envfile")
		}
	}

	// Initialize datalog services for telemetry of the application
	telemetry.Init()
	database.Init()

	fmt.Println("Starting bot...")

	client := disgord.New(disgord.Config{
		BotToken: os.Getenv("TOKEN"),
	})

	handler.Client = client

	client.On(disgord.EvtMessageCreate, handler.OnMessage)
	client.On(disgord.EvtMessageUpdate, handler.OnMessageUpdate)
	client.On(disgord.EvtReady, onReady)
	client.On(disgord.EvtGuildCreate, onGuildCreate)
	client.On(disgord.EvtGuildDelete, onGuildDelete)
	client.StayConnectedUntilInterrupted(context.Background())

	fmt.Println("Good bye!")
}
