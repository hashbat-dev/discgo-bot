package constants

const (
	MAX_MESSAGE_LENGTH    int    = 2000
	LOOP_INTERVAL_SECONDS int    = 10
	TEMP_DIRECTORY        string = "tmp"

	// TTS API settings
	TTS_CHECK_ATTEMPTS int = 150
	TTS_CHECK_DELAY    int = 2000

	// TO DO Settings
	TODO_FIRST_ID int = 1000

	// Channels
	CHANNEL_BOT_TESTING     string = "1222571691851518024" // #bot-general
	CHANNEL_BOT_ERRORS_LIVE string = "1226998385999872030" // #bot-log-live
	CHANNEL_BOT_ERRORS_DEV  string = "1228049795034251284" // #bot-log-dev
	CHANNEL_BOT_FEEDBACK    string = "1225470121900773408" // #bot-feedback
	CHANNEL_BOT_SPAM        string = "1226264005056794674" // #bot-spam
	CHANNEL_BOT_TEST        string = "1244628773433905172" // #bot-dev-only
	CHANNEL_GENERAL         string = "1114324765973434410" // #general
	CHANNEL_AUTO_MEME       string = "1244646268010102917" // #auto-meme

	// Image Resources ===============================================
	IMG_ERROR_TOP    string = "https://i.imgur.com/qdbZadf.png"
	IMG_ERROR_BOTTOM string = "https://i.imgur.com/jsQV8dx.png"

	// Emojis ================================================
	EMOTE_THUMB_UP   string = "👍🏻"
	EMOTE_THUMB_DOWN string = "👎🏻"

	// Discord User IDs ==============================================
	USER_ID_POG    string = "711416363150737508"
	USER_ID_ZEST   string = "192015008039698432"
	USER_ID_CALLUM string = "261978245212143626"

	// Gifs ==========================================================
	GIF_AE_CRY string = "https://cdn.discordapp.com/emojis/610092066080292874.gif?size=96&quality=lossless"
	GIF_UPDATE string = "https://cdn.discordapp.com/attachments/1124062364195627128/1228112087650406462/monkey-orangutan.gif?ex=662adb82&is=66186682&hm=d2b8adc34b1946b2f844993e061d09353fdce91a7e229b99cdbbb224647fb132&"

	// Server Roles ==================================================
	ROLE_BOT_DEVELOPER string = "1226626550850392064"
)

// We cannot make arrays/slices constants but we keep it here as we don't intend to change this
var ERROR_RAND_TEXT = []string{
	"I had a fucky wucky >w<! Pwease don't hate me",
	"Please link your PSN to Bottom Bot",
	"Something went wrong >w< It was probably Callum's fault",
	"Waaaaaaah",
	"KLSADSALKDSAHLKDH >w<",
	"Action failed due to new Brexit regulations",
	"I was on Second Hook >w<",
	"Pog bottomed out so much I was unable to handle the submissiveness",
	"Fuck this, Me and Levi are gonna go smoko",
	"I'M ON SMOKO, SO LEAVE ME ALONE",
	"Absolutely fucked it",
}
