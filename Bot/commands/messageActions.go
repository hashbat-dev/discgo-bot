package commands

type MessageAction struct {
	MessageAliases        []string
	MessageKeepOrigin     bool
	MessageDontNeedExMark bool
	DoTrackWord           bool
	DoAnimeGif            bool
	DoGifBank             bool
	DoAddGif              bool
	DoDelGif              bool
	DoFart                bool
	Category              string
	Description           string
	AdminOnly             bool
	DoFakeYouSearch       bool
	DoFakeYouTTS          bool
	DoAsk                 bool
	DoTranslate           bool
	DoRacism              bool
	DoWowStat             bool
	DoPokeSlur            bool
	DoImageWork           bool
	DoJason               bool
	DoImageWorkType       string
	ReturnSpecificMedia   bool
	Adventure             bool
}

var MessageActions = map[string]MessageAction{
	"fart": {
		MessageAliases: []string{
			"addfart", "farted", "justfarted", "brap", "addbrap",
		},
		DoFart: true,
	},
	"addfartgif": {
		DoAddGif:          true,
		MessageKeepOrigin: true,
		Category:          "fart",
		AdminOnly:         true,
	},
	"delfartgif": {
		DoDelGif:  true,
		Category:  "fart",
		AdminOnly: true,
	},
	"tooth": {
		DoGifBank: true,
	},
	"speech": {
		DoGifBank: true,
	},
	"addspeech": {
		DoAddGif: true,
		Category: "speech",
	},
	"delspeech": {
		DoDelGif:  true,
		Category:  "speech",
		AdminOnly: true,
	},
	"makespeech": {
		Category:        "speech",
		AdminOnly:       true,
		DoImageWork:     true,
		DoImageWorkType: "makespeech",
	},
	"speen": {
		MessageAliases: []string{
			"spin",
		},
		DoGifBank: true,
	},
	"addspeen": {
		MessageAliases: []string{
			"addspin",
		},
		DoAddGif: true,
		Category: "speen",
	},
	"delspeen": {
		MessageAliases: []string{
			"delspin",
		},
		DoDelGif:  true,
		Category:  "speen",
		AdminOnly: true,
	},
	"blush": {
		DoAnimeGif: true,
	},
	"cry": {
		DoAnimeGif: true,
	},
	"cuddle": {
		MessageAliases: []string{
			"snuggle",
		},
		DoAnimeGif: true,
	},
	"dance": {
		DoAnimeGif: true,
	},
	"facepalm": {
		DoAnimeGif: true,
	},
	"happy": {
		DoAnimeGif: true,
	},
	"highfive": {
		DoAnimeGif: true,
	},
	"hug": {
		DoAnimeGif: true,
	},
	"kick": {
		MessageAliases: []string{
			"boot", "punt",
		},
		DoAnimeGif: true,
	},
	"kiss": {
		MessageAliases: []string{
			"smooch",
		},
		DoAnimeGif: true,
	},
	"laugh": {
		MessageAliases: []string{
			"lol",
		},
		DoAnimeGif: true,
	},
	"nod": {
		MessageAliases: []string{
			"yup", "mhmm",
		},
		DoAnimeGif: true,
	},
	"nom": {
		MessageAliases: []string{
			"omnom",
		},
		DoAnimeGif: true,
	},
	"nope": {
		MessageAliases: []string{
			"nuhuh", "nuh-uh",
		},
		DoAnimeGif: true,
	},
	"pat": {
		MessageAliases: []string{
			"pet", "patpat",
		},
		DoAnimeGif: true,
	},
	"poke": {
		DoAnimeGif: true,
	},
	"pout": {
		MessageAliases: []string{
			"sulk", "hmpff",
		},
		DoAnimeGif: true,
	},
	"punch": {
		MessageAliases: []string{
			"thump", "whack",
		},
		DoAnimeGif: true,
	},
	"shoot": {
		DoAnimeGif: true,
	},
	"shrug": {
		MessageAliases: []string{
			"meh",
		},
		DoAnimeGif: true,
	},
	"slap": {
		DoAnimeGif: true,
	},
	"sleep": {
		MessageAliases: []string{
			"sleepy", "eepy",
		},
		DoAnimeGif: true,
	},
	"smile": {
		DoAnimeGif: true,
	},
	"smug": {
		MessageAliases: []string{
			"ronyo",
		},
		DoAnimeGif: true,
	},
	"stare": {
		DoAnimeGif: true,
	},
	"think": {
		MessageAliases: []string{
			"hmm",
		},
		DoAnimeGif: true,
	},
	"thumbsup": {
		DoAnimeGif: true,
	},
	"wave": {
		MessageAliases: []string{
			"hello", "hey",
		},
		DoAnimeGif: true,
	},
	"wink": {
		DoAnimeGif: true,
	},
	"yawn": {
		DoAnimeGif: true,
	},
	"yeet": {
		DoAnimeGif: true,
	},
	"wow": {
		DoWowStat: true,
	},
	"tts": {
		DoFakeYouTTS:      true,
		MessageKeepOrigin: true,
	},
	"flipleft": {
		DoImageWork:     true,
		DoImageWorkType: "flipleft",
	},
	"flipright": {
		DoImageWork:     true,
		DoImageWorkType: "flipright",
	},
	"flipup": {
		DoImageWork:     true,
		DoImageWorkType: "flipup",
	},
	"flipdown": {
		DoImageWork:     true,
		DoImageWorkType: "flipdown",
	},
	"flipboth": {
		DoImageWork:     true,
		DoImageWorkType: "flipboth",
	},
	"flipall": {
		DoImageWork:     true,
		DoImageWorkType: "flipall",
	},
	"pokeslur": {
		DoPokeSlur:        true,
		MessageKeepOrigin: true,
	},
	"addgoon": {
		MessageKeepOrigin: true,
		DoAddGif:          true,
		Category:          "goon",
	},
	"delgoon": {
		MessageKeepOrigin: true,
		DoDelGif:          true,
		Category:          "goon",
		AdminOnly:         true,
	},
	"addbruh": {
		MessageKeepOrigin: true,
		DoAddGif:          true,
		Category:          "bruh",
	},
	"delbruh": {
		MessageKeepOrigin: true,
		DoDelGif:          true,
		Category:          "bruh",
		AdminOnly:         true,
	},
	"addpog": {
		MessageKeepOrigin: true,
		DoAddGif:          true,
		Category:          "pog",
	},
	"delpog": {
		MessageKeepOrigin: true,
		DoDelGif:          true,
		Category:          "pog",
		AdminOnly:         true,
	},
	"pog": {
		MessageAliases: []string{
			"poggers", "pogchamp",
		},
		DoGifBank: true,
		Category:  "pog",
	},
	"bruh": {
		DoGifBank: true,
		Category:  "bruh",
	},
	"addgoodbot": {
		MessageKeepOrigin: true,
		DoAddGif:          true,
		Category:          "goodbot",
		AdminOnly:         true,
	},
	"delgoodbot": {
		MessageKeepOrigin: true,
		DoDelGif:          true,
		Category:          "goodbot",
		AdminOnly:         true,
	},
	"addbadbot": {
		MessageKeepOrigin: true,
		DoAddGif:          true,
		Category:          "badbot",
		AdminOnly:         true,
	},
	"delbadbot": {
		MessageKeepOrigin: true,
		DoDelGif:          true,
		Category:          "badbot",
		AdminOnly:         true,
	},
	"reverse": {
		DoImageWork:     true,
		DoImageWorkType: "reverse",
	},
	"translate": {
		DoTranslate: true,
	},
	"slur": {
		MessageAliases: []string{
			"racism",
			"abuse",
		},
		DoRacism: true,
	},
	"deepfry": {
		DoImageWork:     true,
		DoImageWorkType: "deepfry",
		MessageAliases: []string{
			"cook",
			"fry",
			"frythis",
		},
	},
	"job": {
		MessageAliases: []string{
			"jobslur",
		},
		DoRacism: true,
	},
	"jason": {
		MessageAliases: []string{
			"jasonstatham",
			"statham",
		},
		DoJason: true,
	},

	"adventure": {
		MessageAliases: []string{
			"adv",
			"quest",
		},
		Adventure: true,
	},

	// Keep at the bottom!
	"sound": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"penis": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"shart": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"piss": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"edge": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"cum": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"goon": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"mischief": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"cigarette": {
		MessageKeepOrigin:     true,
		MessageDontNeedExMark: true,
		DoTrackWord:           true,
	},
	"infidel": {
		ReturnSpecificMedia: true,
	},
	"speedup": {
		DoImageWork:     true,
		DoImageWorkType: "speedup",
	},
	"slowdown": {
		DoImageWork:     true,
		DoImageWorkType: "slowdown",
	},
}
