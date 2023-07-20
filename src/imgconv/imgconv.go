package imgconv

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"pp/src/imgtype"

	"github.com/chai2010/webp"
)

// ImgConverter is an interface for converting images between imgtypes.
type ImgConverter interface {
	ConvertImg(imgData []byte, srcimgtype, destimgtype imgtype.ImageFormat, quality int) ([]byte, error)
}

// ImgConverterImpl implements the ImgConverter interface for image imgtype conversion.
type ImgConvert struct{}

// NewImgConvert creates a new instance of ImgConvert that implements the ImgConverter interface.
func NewImgConvert() ImgConverter {
	return &ImgConvert{}
}

// ConvertImg converts the image data between imgtypes.
func (ic *ImgConvert) ConvertImg(imgData []byte, srcimgtype, destimgtype imgtype.ImageFormat, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	switch destimgtype {
	case imgtype.JPEG:
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
	case imgtype.PNG:
		err = png.Encode(buf, img)
	case imgtype.WebP:
		err = webp.Encode(buf, img, &webp.Options{Quality: float32(quality)})
	default:
		return nil, fmt.Errorf("unsupported destination imgtype: %s", imgtype.GetImageFormatName(destimgtype))
	}

	if err != nil {
		log.Printf("error converting image: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}
