package todo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	dbhelper "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	"github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func ToDoAdd(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	optionMap := helpers.GetOptionMap(interaction)
	inCategory := helpers.GetOptionStringValue(optionMap, "category")
	inItem := helpers.GetOptionStringValue(optionMap, "item")

	if inCategory == "" || inItem == "" {
		audit.Log(fmt.Sprintf("Not all values provided. inCategory=%v, inItem=%v", inCategory, inItem))
	}

	newKey, err := dbhelper.ToDoDBAdd(inCategory, inItem, interaction.Member.User.ID)
	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	} else if newKey == "" {
		audit.Error(errors.New("new Key returned blank"))
		logging.SendErrorInteraction(interaction)
		return
	}

	logging.SendMessageInteraction(interaction, newKey+" created", inItem, "", "", false)

}

func ToDoEdit(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	optionMap := helpers.GetOptionMap(interaction)
	inID := helpers.GetOptionStringValue(optionMap, "id")

	// Is the provided ID Valid?
	if inID == "" {
		audit.Error(errors.New("provided id was blank"))
		logging.SendErrorInteraction(interaction)
		return
	}

	dbCat, dbID := GetCategoryAndAssignedIDFromID(inID)
	if dbCat == "" || dbID == 0 {
		audit.Error(errors.New("could not parse provided id"))
		logging.SendErrorInteraction(interaction)
		return
	}

	valid, err := dbhelper.ToDoDBIsIDValid(dbCat, dbID)
	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	} else if !valid {
		audit.Error(errors.New("could not find provided id"))
		logging.SendErrorInteraction(interaction)
		return
	}

	// Do we have at least 1 other parameter?
	inStarted := helpers.GetOptionUserIDValue(optionMap, "started")
	inFinished := helpers.GetOptionUserIDValue(optionMap, "finished")
	inNewText := helpers.GetOptionStringValue(optionMap, "newtext")
	inNewCategory := helpers.GetOptionStringValue(optionMap, "newcategory")
	inVersion := helpers.GetOptionStringValue(optionMap, "version")

	if inStarted == "" && inFinished == "" && inNewText == "" && inNewCategory == "" && inVersion == "" {
		logging.SendErrorMsgInteraction(interaction, "Couldn't Update", "No new values were provided", false)
		return
	}

	// Do the Update
	err = dbhelper.ToDoDBUpdate(dbCat, dbID, inStarted, inFinished, inNewText, inNewCategory, inVersion)
	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	}

	logging.SendMessageInteraction(interaction, inID+" updated", helpers.GetEmote("yup", true), "", "", false)

}

func ToDoDelete(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	optionMap := helpers.GetOptionMap(interaction)
	inConfirm := helpers.GetOptionStringValue(optionMap, "delete-confirmation")
	if strings.ToLower(inConfirm) != "delete" {
		logging.SendErrorMsgInteraction(interaction, "Error Deleting", "You didn't confirm! "+helpers.GetEmote("ae_cry", true), false)
		return
	}

	inID := helpers.GetOptionStringValue(optionMap, "id")
	if inID == "" {
		audit.Error(errors.New("provided id was blank"))
		logging.SendErrorInteraction(interaction)
		return
	}

	dbCat, dbID := GetCategoryAndAssignedIDFromID(inID)
	if dbCat == "" || dbID == 0 {
		audit.Error(errors.New("could not parse provided id"))
		logging.SendErrorInteraction(interaction)
		return
	}

	valid, err := dbhelper.ToDoDBIsIDValid(dbCat, dbID)
	if !valid {
		audit.Error(errors.New("could not find provided id"))
		logging.SendErrorInteraction(interaction)
		return
	} else if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	}

	err = dbhelper.ToDoDelete(dbCat, dbID)
	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	} else {
		logging.SendMessageInteraction(interaction, dbCat+"-"+strconv.Itoa(dbID)+" Deleted", "Item has been deleted "+helpers.GetEmote("ae_cry", true), "", "", false)
	}
}

func ToDoList(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

	optionMap := helpers.GetOptionMap(interaction)
	inCategory := helpers.GetOptionStringValue(optionMap, "category")

	results, err := dbhelper.ToDoGetList(inCategory)
	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	}

	for _, result := range results {

		e := embed.NewEmbed()
		e.SetTitle(result.Category + "-" + strconv.Itoa(result.AssignedID))
		e.SetDescription(result.ToDoText)

		footerText := "Created By: " + helpers.GetNicknameFromID(interaction.GuildID, result.CreatedBy) + " - " + result.CreatedDateTime.Format("02/01/2006 15:04")
		if result.StartedBy != "" {
			footerText += "\nStarted By: " + helpers.GetNicknameFromID(interaction.GuildID, result.StartedBy) + " - " + result.StartedDateTime.Format("02/01/2006 15:04")
		}
		if result.FinishedBy != "" {
			footerText += "\nFinished By: " + helpers.GetNicknameFromID(interaction.GuildID, result.FinishedBy) + " - " + result.FinishedDateTime.Format("02/01/2006 15:04")
		}
		if result.Version != "" {
			footerText += "\nReleased in Version: " + result.Version
		}
		e.SetFooter(footerText)

		_, embederr := config.Session.ChannelMessageSendEmbed(interaction.ChannelID, e.MessageEmbed)
		if embederr != nil {
			audit.Error(embederr)
		}

	}

}
