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
	// TODO: remove this and write server func to read out existing channel names
	CHANNEL_BOT_TESTING     string = "1269632678081073224" // #bot-general /
	CHANNEL_BOT_ERRORS_LIVE string = "1269632717084168267" // #bot-log-live /
	CHANNEL_BOT_ERRORS_DEV  string = "1269632749560529011" // #bot-log-dev /
	CHANNEL_BOT_FEEDBACK    string = "1269632775099645993" // #bot-feedback /
	CHANNEL_BOT_SPAM        string = "1269632800664060055" // #bot-spam /
	CHANNEL_BOT_TEST        string = "1269632838555402302" // #bot-dev-only /
	CHANNEL_GENERAL         string = "1269346433408958579" // #general /

	// Image Resources ===============================================
	IMG_ERROR_TOP    string = "https://i.imgur.com/qdbZadf.png"
	IMG_ERROR_BOTTOM string = "https://i.imgur.com/jsQV8dx.png"

	// Emojis ================================================
	EMOTE_THUMB_UP   string = "👍🏻"
	EMOTE_THUMB_DOWN string = "👎🏻"

	// Gifs ==========================================================
	GIF_AE_CRY string = "https://cdn.discordapp.com/emojis/610092066080292874.gif?size=96&quality=lossless"
	GIF_UPDATE string = "https://cdn.discordapp.com/attachments/1124062364195627128/1228112087650406462/monkey-orangutan.gif?ex=662adb82&is=66186682&hm=d2b8adc34b1946b2f844993e061d09353fdce91a7e229b99cdbbb224647fb132&"

	// Server Roles ==================================================
	ROLE_BOT_DEVELOPER string = "1269347353010114682"
)

// We cannot make arrays/slices constants but we keep it here as we don't intend to change this
var ERROR_RAND_TEXT = []string{
	"Action failed due to new Brexit regulations",
	"On smoke break..",
	"Absolutely fucked it",
}
