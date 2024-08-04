// Called when a New Message is entered into any channel in the Server
package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/commands"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/constants"
	dbhelpers "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
	"github.com/dabi-ngin/discgo-bot/Bot/gifbank"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/adventures"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/animegif"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/fakeyou"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/imagework"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/jason"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/reactions"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/returnspecificmedia"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/sentience"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/toothjack"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/trackword"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers/wow"

	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	logger "github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func NewMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if config.IsDev {
		if message.ChannelID != constants.CHANNEL_BOT_TEST {
			return
		}
	} else {
		if message.ChannelID == constants.CHANNEL_BOT_TEST {
			return
		}
	}
	go ProcessMessage(session, message)
}

func ProcessMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Check for empty messages
	if len(message.Content) == 0 {
		return
	}
	var err error
	var skipped bool
	if skipped, err = helpers.SkipProcessing(session, message, nil); skipped {
		return
	}
	if err != nil {
		audit.Error(err)
	}

	// Detect for a !command style message
	botCommand := ""
	msgWithoutCommand := ""
	firstChar := message.Content[0]
	if firstChar == '!' {
		inMsg := strings.TrimSpace(message.Content)

		// Check for spaces
		if strings.Contains(inMsg, " ") {
			parts := strings.SplitN(inMsg, " ", 2)
			if len(parts) > 1 {
				botCommand = strings.ToLower(parts[0])
			}
		} else {
			botCommand = strings.ToLower(inMsg)
		}
		msgWithoutCommand = strings.TrimSpace(strings.Replace(message.Content, botCommand, "", -1))
		botCommand = botCommand[1:]
	}

	ttsCommand := ""
	if botCommand != "" {
		if botCommand[:3] == "tts" {
			ttsCommand = strings.Replace(botCommand, "tts", "", -1)
			botCommand = "tts"
		}
	}

	deleteOriginalMsg := false
	msgLower := strings.ToLower(message.Content)
	matched := false

	// Check through our MessageActions
	for keyword, action := range commands.MessageActions {

		// Check if the message matches a MessageAction
		needExclaimMark := !action.MessageDontNeedExMark
		if (!needExclaimMark && strings.Contains(msgLower, keyword)) || (needExclaimMark && botCommand == keyword) {
			matched = true
		} else if action.MessageAliases != nil && len(action.MessageAliases) > 0 {
			for _, alias := range action.MessageAliases {
				if (!needExclaimMark && strings.Contains(msgLower, alias)) || (needExclaimMark && botCommand == alias) {
					matched = true
					break
				}
			}
		}

		if matched {
			if action.AdminOnly && !helpers.IsAdmin(message.Message.Member) {
				logger.SendErrorMsg(message, "This command is for Bot Developers only! You're not a bot developer!")
			} else {
				deleteOriginalMsg = !action.MessageKeepOrigin
				if action.DoTrackWord && botCommand == "" {
					trackword.TrackWord(session, message, keyword)
				} else if action.DoTrackWord && botCommand != "" {
					matched = false
				} else if action.DoAnimeGif {
					animegif.AnimeGif(session, message, keyword)
				} else if action.DoGifBank {
					gifbank.Post(session, message, keyword)
				} else if action.DoAddGif {
					gifbank.AddGIF(message, action.Category)
				} else if action.DoDelGif {
					gifbank.Delete(message, action.Category)
				} else if action.DoFakeYouTTS {
					deleteOriginalMsg = fakeyou.FakeYouGetTTS(message, ttsCommand, msgWithoutCommand)
				} else if action.DoWowStat {
					deleteOriginalMsg = wow.GetWowStat(message)
				} else if action.DoJason {
					jason.RequestJason(message)
				} else if action.DoImageWork {
					if action.DoImageWorkType == "flipboth" {
						go imagework.HandleMessage(message, "flipleft")
						go imagework.HandleMessage(message, "flipright")
					} else if action.DoImageWorkType == "flipall" {
						go imagework.HandleMessage(message, "flipleft")
						go imagework.HandleMessage(message, "flipright")
						go imagework.HandleMessage(message, "flipup")
						go imagework.HandleMessage(message, "flipdown")
					} else {
						go imagework.HandleMessage(message, action.DoImageWorkType)
					}
				} else if action.ReturnSpecificMedia {
					returnspecificmedia.HandleMessage(message, "infidel")
				} else if action.Adventure {
					adventures.HandleMessage(message)
				}

				dbhelpers.CountCommand(botCommand+ttsCommand, message.Message.Author.ID)

				if matched {
					break
				}
			}
		}

	}

	if !matched {
		sentRandom := false
		if !sentRandom {
			sentRandom = toothjack.RandToothjack(message, session)
		}
		if !sentRandom {
			sentRandom = sentience.RandomReplies(message)
		}
		if !sentRandom {
			sentRandom = jason.DetectJason(message)
		}
		if !sentRandom {
			wow.WowDetection(message)
		}
	}

	reactions.CheckMessage(message)

	if deleteOriginalMsg {
		helpers.DeleteSourceMessage(session, message)
	}
}
