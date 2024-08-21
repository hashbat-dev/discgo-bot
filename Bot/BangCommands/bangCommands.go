package bangCommands

import (
	imgbank "github.com/dabi-ngin/discgo-bot/Bot/Handlers/ImgBank"
	testhandler "github.com/dabi-ngin/discgo-bot/Bot/Handlers/TestHandler"
	structs "github.com/dabi-ngin/discgo-bot/Structs"
)

var (
	CommandTable = make(map[string]structs.BangCommand)
)

func Init() bool {
	CommandTable["test"] = structs.BangCommand{
		Execute: testhandler.HandleNewMessage,
	}

	CommandTable["speech"] = structs.BangCommand{
		Execute: imgbank.GetImg,
	}

	CommandTable["addspeech"] = structs.BangCommand{
		Execute:     imgbank.AddImg,
		ImgCategory: "speech",
	}

	return true
}
