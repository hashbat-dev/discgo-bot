// Called when a New Message is entered into any channel in the Server
package handlers

import (
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/commands"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	dbhelpers "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/gifbank"
	"github.com/ZestHusky/femboy-control/Bot/handlers/adventures"
	"github.com/ZestHusky/femboy-control/Bot/handlers/animegif"
	"github.com/ZestHusky/femboy-control/Bot/handlers/fakeyou"
	"github.com/ZestHusky/femboy-control/Bot/handlers/fart"
	"github.com/ZestHusky/femboy-control/Bot/handlers/imagework"
	"github.com/ZestHusky/femboy-control/Bot/handlers/jason"
	"github.com/ZestHusky/femboy-control/Bot/handlers/pogcorrection"
	"github.com/ZestHusky/femboy-control/Bot/handlers/reactions"
	"github.com/ZestHusky/femboy-control/Bot/handlers/returnspecificmedia"
	"github.com/ZestHusky/femboy-control/Bot/handlers/sentience"
	"github.com/ZestHusky/femboy-control/Bot/handlers/slurs"
	"github.com/ZestHusky/femboy-control/Bot/handlers/toothjack"
	"github.com/ZestHusky/femboy-control/Bot/handlers/trackword"
	"github.com/ZestHusky/femboy-control/Bot/handlers/translate"
	"github.com/ZestHusky/femboy-control/Bot/handlers/wow"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	logger "github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
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
				} else if action.DoPokeSlur {
					deleteOriginalMsg = slurs.PokeSlur(message)
				} else if action.DoFart {
					fart.LogFart(message)
				} else if action.DoTranslate {
					translate.TranslateBottomSpeech(message)
				} else if action.DoRacism {
					slurs.SendASlur(message, botCommand)
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

	// Now we do Pog Correction checking, ALL Messages need to enter here
	if !matched {
		pogCorr := pogcorrection.Detect(message.Content, message.Author.ID, message.ReferencedMessage)
		switch pogCorr {
		case "tooth":
			dbhelpers.CountCommand("pog-correction", message.Message.Author.ID)
			gifbank.Post(session, message, "tooth")
		case "correct":
			pogcorrection.SendCorrection(message)
			dbhelpers.CountCommand("pog-correction", message.Message.Author.ID)
		}
	}

	sentience.DetectBotTalk(message, msgLower)

	if !matched {
		sentRandom := false
		if !sentRandom {
			sentRandom = toothjack.RandToothjack(message, session)
		}
		if !sentRandom {
			sentRandom = sentience.RandomReplies(message)
		}
		if !sentRandom {
			sentRandom = slurs.RandomSlur(message)
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
