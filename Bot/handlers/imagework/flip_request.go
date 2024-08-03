package imagework

import (
	"bytes"
	"io"

	"github.com/bwmarrin/discordgo"
)

type FlipRequest struct {
	Message        *discordgo.MessageCreate
	ImageReader    io.Reader
	IsGif          bool
	RequestType    string
	ResponseBuffer *bytes.Buffer
	ResponseChan   chan error
}
