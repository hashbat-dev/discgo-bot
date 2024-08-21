package bangCommands

import (
	"errors"

	imgbank "github.com/dabi-ngin/discgo-bot/Bot/Handlers/ImgBank"
	testhandler "github.com/dabi-ngin/discgo-bot/Bot/Handlers/TestHandler"
	structs "github.com/dabi-ngin/discgo-bot/Structs"
)

var (
	commandTable = make(map[string]structs.BangCommand)
)

func Init() bool {

	commandTable["test"] = structs.BangCommand{
		Begin: testhandler.HandleNewMessage,
	}

	commandTable["speech"] = structs.BangCommand{
		Begin: imgbank.GetImg,
	}

	commandTable["addspeech"] = structs.BangCommand{
		Begin:       imgbank.AddImg,
		ImgCategory: "speech",
	}

	return true
}

func GetCommand(query string) (structs.BangCommand, error) {
	if val, ok := commandTable[query]; ok {
		return val, nil
	} else {
		return structs.BangCommand{}, errors.New("not found")
	}
}
