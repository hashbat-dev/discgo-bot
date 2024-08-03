package imgwork

import (
	"bytes"
	"image"
	"image/png"
	"io"
)

func ConvertWebpToPNG(webpReader io.Reader) (io.Reader, error) {
	img, _, err := image.Decode(webpReader)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
