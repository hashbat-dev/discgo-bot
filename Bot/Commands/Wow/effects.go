package wow

import "github.com/bwmarrin/discordgo"

type Effect struct {
	Name        string
	Description string
}

type EffectList func(*discordgo.MessageCreate) (int, *Effect)

var effectList = []EffectList{
	effectSixetyNine,
	effectOhRight,
}

func effectSixetyNine(m *discordgo.MessageCreate) (int, *Effect) {
	if len(m.ID) <= 2 || m.ID[len(m.ID)-2:] != "69" {
		return 0, nil
	}
	i := getRandomNumber(6, 9) * 2
	return i, &Effect{
		Name:        "Niceeee",
		Description: "Message ID ended in 69, get a random roll between 6 and 9 doubled.",
	}
}

func effectOhRight(m *discordgo.MessageCreate) (int, *Effect) {
	if len(m.ID) <= 1 || m.ID[len(m.ID)-1:] != "0" {
		return 0, nil
	}
	i := getRandomNumber(1, 9)
	return i, &Effect{
		Name:        "0h Right",
		Description: "Message ID ended in 0, get a random roll between 1 and 9.",
	}
}
