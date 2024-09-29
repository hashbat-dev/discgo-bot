package logger

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
)

var InitComplete bool

// Initialise the Logger package
func Init() bool {
	if config.ServiceSettings.LOGTODISCORD && config.ServiceSettings.LOGGINGUSESTHREADS {
		threadTitle := fmt.Sprintf("[%v] - %v", time.Now().Format("15:04"), config.ServiceSettings.HOSTNAME)
		threadId, err := config.Session.ThreadStart(config.ServiceSettings.LOGGINGCHANNELID, threadTitle, discordgo.ChannelTypeGuildText, 60)
		if err != nil {
			fmt.Println(fmt.Printf("Log.Init() - Error creating Logging Thread :: %v", err))
			return false
		}

		_, err = config.Session.ChannelMessageSend(config.ServiceSettings.LOGGINGCHANNELID, fmt.Sprintf("<#%v> - Logging thread created", threadId.ID))
		if err != nil {
			fmt.Println(fmt.Printf("Log.Init() - Error posting Thread Notification :: %v", err))
			return false
		}
		config.ServiceSettings.LOGGINGCHANNELID = threadId.ID
	}

	InitComplete = true
	return true
}
