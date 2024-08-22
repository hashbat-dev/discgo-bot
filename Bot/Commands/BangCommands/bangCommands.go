package bangCommands

import (
	commands "github.com/dabi-ngin/discgo-bot/Bot/Commands"
	imgbank "github.com/dabi-ngin/discgo-bot/Bot/Handlers/ImgBank"
)

var (
	CommandTable = make(map[string]commands.Command)
)

func init() {
	CommandTable["reverse"] = Reverse{}

	CommandTable["speech"] = commands.Command{
		Begin: imgbank.GetImg,
	}

	CommandTable["addspeech"] = commands.Command{
		Begin:       imgbank.AddImg,
		ImgCategory: "speech",
	}
}
