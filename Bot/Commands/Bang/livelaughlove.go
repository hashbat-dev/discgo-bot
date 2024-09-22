package bang

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

type LiveLaughLove struct{}

func (s LiveLaughLove) Name() string {
	return "livelaughlove"
}

func (s LiveLaughLove) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s LiveLaughLove) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s LiveLaughLove) Execute(message *discordgo.MessageCreate, command string) error {
	var (
		dpi      = flag.Float64("dpi", 72, "screen resolution in dots per inch")
		fontfile = flag.String("fontfile", "Fonts/Kalam-Light.ttf", "filename of the ttf font")
		hinting  = flag.String("hinting", "none", "none | full")
		size     = flag.Float64("size", 24, "font size in points")
		wonb     = flag.Bool("whiteonblack", false, "white text on black background")
	)

	progressMessage := discord.SendUserMessageReply(message, false, "LiveLaughLove: Capturing message...")
	text := message.Content

	//get font
	fontBytes, fontErr := os.ReadFile(*fontfile)
	if fontErr != nil {
		logger.Error(message.GuildID, fontErr)
		return fontErr
	}

	parsedFont, fontParseErr := freetype.ParseFont(fontBytes)
	if fontParseErr != nil {
		logger.Error(message.GuildID, fontParseErr)
		return fontParseErr
	}

	//open freetype context
	fore, back := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}

	if *wonb {
		fore, back = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}

	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), back, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(parsedFont)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fore)

	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	//Draw guidelines
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, i, ruler)
	}

	//render text
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))

	_, stringDrawErr := c.DrawString(text, pt)
	if stringDrawErr != nil {
		logger.Error(message.GuildID, stringDrawErr)
		return stringDrawErr
	}

	//save rgba to buffer
	var newImageBuffer bytes.Buffer
	encodeErr := png.Encode(&newImageBuffer, rgba)
	if encodeErr != nil {
		logger.Error(message.GuildID, encodeErr)
		return encodeErr
	}

	outputImageName := "test.png"

	// 5. Send the new Image back to the User
	replyErr := discord.ReplyToMessageWithImageBuffer(message, true, outputImageName, &newImageBuffer)
	if replyErr != nil {
		logger.Error(message.GuildID, replyErr)
		return replyErr
	}

	discord.DeleteMessageObject(progressMessage)
	return nil
}
