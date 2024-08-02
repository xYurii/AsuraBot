package events

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/utils"
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

type GuildInfo struct {
	sync.Mutex
	NewLootBoxTime int64
	LastUser       string
}

var cache = map[string]*GuildInfo{}

func SendLootbox(msg *disgord.Message) {
	rarity, lootbox := rinha.MessageRandomLootbox()
	rand := rinha.GetRandByType(rarity)
	sprite := rinha.Sprites[0][rand-1]
	embeds := []*disgord.Embed{
		{
			Color: rarity.Color(),
			Image: &disgord.EmbedImage{
				URL: sprite,
			},
			Description: fmt.Sprintf("Uma lootbox de raridade **%s** apareceu, clique no botão abaixo para adquirir", rinha.LootNames[lootbox]),
			Title:       "Lootbox",
			Footer: &disgord.EmbedFooter{
				Text: "Para desativar ou mudar de canal essas mensagens use o comando /config (precisa ter a permissão de gerenciar servidor)",
			},
		},
	}
	message := disgord.CreateMessage{
		Embeds: embeds,
		Components: []*disgord.MessageComponent{
			{
				Type: disgord.MessageComponentActionRow,
				Components: []*disgord.MessageComponent{
					{
						Type:     disgord.MessageComponentButton,
						Label:    "Pegar Lootbox",
						CustomID: "GetLoobox",
						Style:    disgord.Primary,
					},
				},
			},
		},
	}
	newMessage, err := msg.Reply(context.Background(), handler.Client, message)
	if err == nil {
		handler.RegisterHandler(newMessage.ID, func(ic *disgord.InteractionCreate) {
			done := false
			if done {
				return
			}
			done = true
			u := database.User.GetUser(context.Background(), ic.Member.UserID, "Items")
			database.User.InsertItem(context.Background(), ic.Member.UserID, u.Items, lootbox, entities.LootboxType)
			embeds[0].Description = fmt.Sprintf("Parabéns  <@%s> Você adquiriu uma lootbox **%s**\nUse **/lootbox open** para abrir", ic.Member.UserID, rinha.LootNames[lootbox])
			embeds[0].Color = 16776960
			message.Components[0].Components[0].Disabled = true
			message.Components[0].Components[0].Label = "Lootbox pega"
			handler.Client.SendInteractionResponse(context.Background(), ic, &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackUpdateMessage,
				Data: &disgord.CreateInteractionResponseData{
					Embeds:     embeds,
					Components: message.Components,
				},
			})

		}, 100)
	}
}

func GetGuildInfo(guildID string) *GuildInfo {
	if cache[guildID] == nil {
		fmt.Println("Creating new cache for guild", guildID)
		cache[guildID] = &GuildInfo{}
	}
	return cache[guildID]
}

func IsFlood(msg *disgord.Message, cache *GuildInfo) bool {
	fmt.Println(cache.LastUser)
	if cache.LastUser == msg.Author.ID.String() {
		return true
	}
	cache.LastUser = msg.Author.ID.String()
	return false
}

func setNewLootboxTime(cache *GuildInfo, now int64) {
	randomMinutes := utils.RandInt(150)
	cache.NewLootBoxTime = now + 60*60*1.5 + int64(randomMinutes)*60
}

const MIN_MEMBERS = 1 //change later to a real value

func RecieveLootbox(msg *disgord.Message) {
	guildDb := database.Guild.GetGuild(context.Background(), msg.GuildID)
	if guildDb.DisableLootbox || (guildDb.LootBoxChannel != 0 && guildDb.LootBoxChannel != msg.ChannelID) {
		return
	}
	cache := GetGuildInfo(msg.GuildID.String())
	cache.Lock()
	defer cache.Unlock()
	now := time.Now().Unix()
	guild, _ := handler.Client.Cache().GetGuild(msg.GuildID)
	members := guild.MemberCount
	randomNumber := rand.Intn(100 + int(members/100))
	if msg.GuildID.String() == "710179373860519997" {
		randomNumber = rand.Intn(20)
	}
	if members > MIN_MEMBERS {
		if randomNumber < 3 && now > cache.NewLootBoxTime && !IsFlood(msg, cache) {
			setNewLootboxTime(cache, now)
			go SendLootbox(msg)
		}
	}
}

var activeGuilds = []string{"710179373860519997", "597089324114116635", "1189997575697338441", "862649888290242570", "1266918944971948225", "1177549280156860427"}

func HandleMessage(s disgord.Session, h *disgord.MessageCreate) {
	msg := h.Message
	appID := os.Getenv("APP_ID")
	if !msg.Author.Bot {
		if msg.GuildID != 0 {
			for _, user := range msg.Mentions {
				if user.ID.String() == appID {
					msg.Reply(context.Background(), s, "Use /help para ver meus comandos\nCaso meus comandos não aparecam me readicione no servidor com este link:\nhttps://discordapp.com/oauth2/authorize?client_id=470684281102925844&scope=applications.commands%%20bot&permissions=8")
					break
				}
			}

			if utils.Includes(activeGuilds, msg.GuildID.String()) {
				RecieveLootbox(msg)
			}
		}
	}
}
