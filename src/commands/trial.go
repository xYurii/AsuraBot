package commands

import (
	"asura/src/database"
	"asura/src/entities"
	"asura/src/handler"
	"asura/src/rinha"
	"asura/src/rinha/engine"
	"asura/src/utils"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "trial",
		Category:    handler.Rinha,
		Description: "Faça desafios pro seu galo ficar mais forte",
		Run:         runTrial,
		Cooldown:    7,
		Options: utils.GenerateOptions(
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "status",
				Description: "Veja seu status atual",
			},
			&disgord.ApplicationCommandOption{
				Type:        disgord.OptionTypeSubCommand,
				Name:        "battle",
				Description: "Va para o proximo desafio",
			},
		),
	})
}

var trialDamageMultiplier = 4

const MAX_TRIALS = 5
const TRIAL_MIN_LEVEL = 15

func getTrial(trials []*entities.Trial, rooster int) *entities.Trial {
	for _, trial := range trials {
		if trial.Rooster == rooster {
			return trial
		}
	}
	return nil
}

func runTrial(ctx context.Context, itc *disgord.InteractionCreate) *disgord.CreateInteractionResponse {
	user := database.User.GetUser(ctx, itc.Member.UserID, "Items", "Galos", "Trials")
	command := itc.Data.Options[0].Name
	rooster := rinha.GetEquippedGalo(&user)
	trial := getTrial(user.Trials, rooster.Type)
	if trial == nil {
		database.User.InsertTrial(ctx, itc.Member.UserID, &entities.Trial{
			Rooster: rooster.Type,
		})
		user = database.User.GetUser(ctx, itc.Member.UserID, "Items", "Galos", "Trials")
		trial = getTrial(user.Trials, rooster.Type)
	}
	galoSprite := rinha.GetGaloImage(rooster, user.Items)
	switch command {
	case "status":
		return &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Embeds: []*disgord.Embed{
					{
						Title:       "Trials",
						Description: fmt.Sprintf("Seu galo está na batalha **%d/%d** das Trials.\nBônus de dano atual: **%d%%**\n\nNa última batalha da Trial, você ganhará uma lootbox de acordo com a raridade de seu galo.\nUse `/trial battle` para lutar com os bosses das Trials!", trial.Win, MAX_TRIALS, trialDamageMultiplier*trial.Win),
						Thumbnail: &disgord.EmbedThumbnail{
							URL: galoSprite,
						},
						Color: 65535,
					},
				},
			},
		}
	case "battle":
		if trial.Win >= MAX_TRIALS {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: "Você já completou todas as Trials desse galo!",
				},
			}
		}
		roosterLevel := rinha.CalcLevel(rooster.Xp)
		if roosterLevel < TRIAL_MIN_LEVEL {
			return &disgord.CreateInteractionResponse{
				Type: disgord.InteractionCallbackChannelMessageWithSource,
				Data: &disgord.CreateInteractionResponseData{
					Content: fmt.Sprintf("Seu galo precisa ser no mínimo nível **%d** para entrar nas Trials!", TRIAL_MIN_LEVEL),
				},
			}
		}
		discordUser := itc.Member.User
		authorRinha := isInRinha(ctx, discordUser)
		if authorRinha != "" {
			return rinhaMessage(discordUser.Username, authorRinha)
		}
		lockEvent(ctx, discordUser.ID, "Desafio trial")
		defer unlockEvent(ctx, discordUser.ID)
		class := rinha.Classes[trial.Rooster]
		xpMultiplier := 1
		if class.Rarity == rinha.Epic {
			xpMultiplier = 3 + rooster.Resets
		}
		if class.Rarity > rinha.Epic {
			xpMultiplier = (int(class.Rarity) + rooster.Resets) * 4
		}
		xp := rinha.CalcXP(15+(trial.Win*7)) * xpMultiplier
		level := rinha.CalcLevel(xp)
		galoAdv := &entities.Rooster{
			Type:    rinha.GetRandByType(class.Rarity),
			Xp:      xp,
			Evolved: class.Rarity > rinha.Epic || level >= 40,
			Resets:  rooster.Resets * (1 + trial.Win),
		}

		userAdv := entities.User{
			Galos: []*entities.Rooster{galoAdv},
		}

		if class.Rarity > rinha.Epic {
			atbs := (user.UserXp / 150) * ((trial.Win / 2) + 1)
			if class.Rarity == rinha.Legendary {
				atbs = atbs / 2
			}
			healthAtb := float64(atbs) * 0.65
			userAdv.Attributes = [5]int{int(healthAtb) + 600, atbs / 10, atbs, atbs / 6, atbs / 6}
		}

		itc.Reply(ctx, handler.Client, &disgord.CreateInteractionResponse{
			Type: disgord.InteractionCallbackChannelMessageWithSource,
			Data: &disgord.CreateInteractionResponseData{
				Content: "A batalha está iniciando...",
			},
		})

		winner, _ := engine.ExecuteRinha(itc, handler.Client, engine.RinhaOptions{
			GaloAuthor: &user,
			GaloAdv:    &userAdv,
			IDs:        [2]disgord.Snowflake{discordUser.ID},

			AuthorName:  rinha.GetName(discordUser.Username, *rooster),
			AdvName:     "Boss " + rinha.Classes[galoAdv.Type].Name,
			AuthorLevel: rinha.CalcLevel(rooster.Xp),
			AdvLevel:    rinha.CalcLevel(galoAdv.Xp),
			NoItems:     false,
		}, false)

		if winner == -1 {
			return nil
		}

		ch := handler.Client.Channel(disgord.Snowflake(itc.ChannelID))

		if winner == 0 {
			database.User.AddTrialWin(ctx, trial)
			if trial.Win >= MAX_TRIALS {
				lootbox := rinha.GetTrialLootbox(class.Rarity)
				database.User.InsertItem(ctx, itc.Member.UserID, user.Items, lootbox, entities.LootboxType)
				lootboxName := rinha.LootNames[lootbox]
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{{
						Color:       65535,
						Title:       "Trial",
						Description: fmt.Sprintf("Parabéns, você conquistou a última etapa das Trials do galo **%s**, e ganhou uma lootbox **%s**!", class.Name, lootboxName),
					}},
				})
			} else {
				ch.CreateMessage(&disgord.CreateMessage{
					Embeds: []*disgord.Embed{{
						Color:       16776960,
						Title:       "Trial",
						Description: fmt.Sprintf("Você venceu a etapa **%d/%d** das Trials!\nGanhou **%d%%** de dano bônus para seu galo.", trial.Win, MAX_TRIALS, trialDamageMultiplier*trial.Win),
					}},
				})
			}
			return nil
		}
		ch.CreateMessage(&disgord.CreateMessage{
			Embeds: []*disgord.Embed{{
				Color:       16711680,
				Title:       "Trial",
				Description: "Infelizmente, você perdeu. Use /trial battle para tentar novamente!",
			}},
		})
	}

	return nil
}
