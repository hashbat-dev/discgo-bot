package meme

import (
	"errors"
	"net/url"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

func GetRandomMeme(interaction *discordgo.InteractionCreate, loopCount int) {

	// Put a Loading Interaction in
	embedStart := embed.NewEmbed()
	embedStart.SetDescription("Searching for the dankest meme...")

	if loopCount == 0 {
		err := config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embedStart.MessageEmbed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			audit.Error(err)
			logging.SendErrorInteraction(interaction)
			return
		}
	}

	// Get the Search term
	optionMap := helpers.GetOptionMap(interaction)
	inSearch := helpers.GetOptionStringValue(optionMap, "search")
	inStills := helpers.GetOptionBoolValue(optionMap, "allow-stills")
	searchTerm := GetSearchTerm(inSearch)

	fileUrl, threadUrl, threadTitle := GetMeme(searchTerm, inStills, "")
	if fileUrl == "" {
		if inSearch == "" && loopCount < len(GenericSearch) {
			GetRandomMeme(interaction, loopCount+1)
		} else {
			config.Session.InteractionResponseEdit(interaction.Interaction, ErrorWebHook())
			return
		}
	}

	// Post the Meme ===================================================================================
	if fileUrl == "" {
		config.Session.InteractionResponseEdit(interaction.Interaction, ErrorWebHook())
		audit.Error(errors.New("fileUrl was blank"))
		return
	}
	msg, err := config.Session.ChannelMessageSend(interaction.ChannelID, FormatMessage(fileUrl, threadUrl, threadTitle))
	if err != nil {
		config.Session.InteractionResponseEdit(interaction.Interaction, ErrorWebHook())
		audit.Error(err)
		return
	}

	// Post the Success interaction ====================================================================
	passStills := "false"
	if inStills {
		passStills = "true"
	}
	successEmbed := embed.NewEmbed()
	successEmbed.SetDescription("All done! Enjoy your meme (Searched for: " + searchTerm + ")...")
	var successHook discordgo.WebhookEdit
	successHook.Embeds = &[]*discordgo.MessageEmbed{successEmbed.MessageEmbed}
	successHook.Components = &[]discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Regenerate",
					CustomID: "regenerate-" + interaction.ID + "-" + msg.ID + "-" + url.QueryEscape(inSearch) + "-" + passStills,
					Style:    discordgo.SecondaryButton,
					Emoji: discordgo.ComponentEmoji{
						Name: "🔄",
					},
				},
				discordgo.Button{
					Label:    "Delete",
					CustomID: "delete-" + interaction.ID + "-" + msg.ID + "-" + url.QueryEscape(inSearch),
					Style:    discordgo.DangerButton,
					Emoji: discordgo.ComponentEmoji{
						Name: "🗑️",
					},
				},
			},
		},
	}

	_, err = config.Session.InteractionResponseEdit(interaction.Interaction, &successHook)
	if err != nil {
		audit.Error(err)
	} else {
		AddToInteractionCache(interaction)
	}

}

func HandleRegenerate(interaction *discordgo.InteractionCreate, interactionId string, messageId string, inSearch string, inStills string) {

	// This stops the "(!) This interaction failed" error, do this first as it's Time Imperative
	err := FakeInteractionResponse(interaction)
	if err != nil {
		audit.Error(err)
	}

	originalInteraction, err := GetFromInteractionCache(interactionId)
	if err != nil {
		audit.Error(err)
		return
	}

	allowStills := false
	if inStills != "" {
		if inStills == "true" {
			allowStills = true
		}
	}

	// Set the Original Interaction to the loading state
	e := embed.NewEmbed()
	e.SetDescription("Getting you another Meme...")
	_, err = config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{e.MessageEmbed},
	})
	if err != nil {
		audit.Error(err)
		return
	}

	// Get the Search Term which was used
	unEscaped, err := url.QueryUnescape(inSearch)
	if err != nil {
		config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{ErrorEmbed()},
		})
		return
	}

	searchTerm := GetSearchTerm(unEscaped)
	fileUrl, threadUrl, threadTitle := GetMeme(searchTerm, allowStills, "")

	// Attempt 2? - If for whatever reason we wanted a random category and there were no threads
	if fileUrl == "" {
		searchTerm = GetSearchTerm(inSearch)
		fileUrl, threadUrl, threadTitle = GetMeme(searchTerm, allowStills, "")
	}

	if fileUrl == "" {
		config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{ErrorEmbed()},
		})
		return
	} else {
		_, err := config.Session.ChannelMessageEdit(interaction.ChannelID, messageId, FormatMessage(fileUrl, threadUrl, threadTitle))
		if err != nil {
			audit.Error(err)
			config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{ErrorEmbed()},
			})
		} else {
			successEmbed := embed.NewEmbed()
			successEmbed.SetDescription("All done! Enjoy your meme (Searched for: " + searchTerm + ")...")
			_, err := config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{successEmbed.MessageEmbed},
			})
			if err != nil {
				audit.Error(err)
			}
		}
	}
}

func HandleDelete(interaction *discordgo.InteractionCreate, interactionId string, messageId string) {

	// This stops the "(!) This interaction failed" error, do this first as it's Time Imperative
	err := FakeInteractionResponse(interaction)
	if err != nil {
		audit.Error(err)
	}

	originalInteraction, err := GetFromInteractionCache(interactionId)
	if err != nil {
		audit.Error(err)
		return
	}

	// Set the Original Interaction to the loading state
	e := embed.NewEmbed()
	e.SetDescription("Deleting Meme...")
	_, err = config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{e.MessageEmbed},
	})
	if err != nil {
		audit.Error(err)
		return
	}

	err = config.Session.ChannelMessageDelete(originalInteraction.ChannelID, messageId)
	if err != nil {
		audit.Error(err)
		config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{ErrorEmbed()},
		})
	} else {
		e := embed.NewEmbed()
		e.SetDescription("Meme Deleted, forget this ever happened >_>")

		_, err := config.Session.InteractionResponseEdit(originalInteraction.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{e.MessageEmbed},
		})
		if err != nil {
			audit.Error(err)
		}
	}
}
