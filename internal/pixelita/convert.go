package pixelita

import (
	"app/internal/imagetype"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/chai2010/webp"
)

func (px *Pixelita) convertBetweenFormats(imgData []byte, srcImageType, destImageType imagetype.ImageFormat, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	switch destImageType {
	case imagetype.JPEG:
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	case imagetype.PNG:
		err = png.Encode(buf, img)
	case imagetype.WebP:
		err = webp.Encode(buf, img, &webp.Options{Quality: float32(quality)})
	default:
		return nil, fmt.Errorf("unsupported destination imgtype: %s", imagetype.New().GetFormatName(destImageType))
	}

	if err != nil {
		log.Printf("error converting image: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
