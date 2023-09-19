package imageencode

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"log"

	"github.com/chai2010/webp"
	"github.com/ultimate-guitar/go-imagequant"
)

type JPEGEncoder struct{}

func NewJPEGEncoder() *JPEGEncoder {
	return &JPEGEncoder{}
}

func (je *JPEGEncoder) Encode(imageData []byte, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	var encodedImage bytes.Buffer
	err = jpeg.Encode(&encodedImage, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}

	return encodedImage.Bytes(), nil
}

type PNGEncoder struct{}

func NewPNGEncoder() *PNGEncoder {
	return &PNGEncoder{}
}

func (pe *PNGEncoder) Encode(imageData []byte, minPosterization, minQuality, maxQuality, speed int) ([]byte, error) {
	buf := new(bytes.Buffer)
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	attr, err := imagequant.NewAttributes()
	if err != nil {
		log.Printf("failed to create image attributes: %v", err)
		return nil, err
	}
	defer attr.Release()

	// Ignores given number of least significant bits in all channels, posterizing image to 2^bits levels. 0 gives full quality. Use 2 for VGA or 16-bit RGB565 displays. 4 if image is going to be output on a RGB444/RGBA4444 display (e.g. low-quality textures on Android).
	attr.SetMinPosterization(minPosterization)
	if err != nil {
		log.Printf("failed to set min. posterization: %v", err)
		return nil, err
	}

	// Quality is in range 0 (worst) to 100 (best) and values are analoguous to JPEG quality (i.e. 80 is usually good enough). Quantization will attempt to use the lowest number of colors needed to achieve maximum quality. maximum value of 100 is the default and means conversion as good as possible. If it's not possible to convert the image with at least minimum quality (i.e. 256 colors is not enough to meet the minimum quality), then Image.Quantize() will fail. The default minimum is 0 (proceeds regardless of quality).
	// Features dependent on speed: speed 1-5: Noise-sensitive dithering speed 8-10 or if image has more than million colors: Forced posterization speed 1-7 or if minimum quality is set: Quantization error known seed 1-6: Additional quantization techniques
	err = attr.SetQuality(minQuality, maxQuality)
	if err != nil {
		log.Printf("failed to set quality: %v", err)
		return nil, err
	}

	// Higher speed levels disable expensive algorithms and reduce quantization precision. The default speed is 3. Speed 1 gives marginally better quality at significant CPU cost. Speed 10 has usually 5% lower quality, but is 8 times faster than the default. High speeds combined with Attributes.SetQuality() will use more colors than necessary and will be less likely to meet minimum required quality.
	err = attr.SetSpeed(speed)
	if err != nil {
		log.Printf("failed to set speed: %v", err)
		return nil, err
	}

	rgba32data := imageToRGBA32(img)

	imq, err := imagequant.NewImage(attr, string(rgba32data), width, height, 0)
	if err != nil {
		log.Printf("failed to create image quantization: %v", err)
		return nil, err
	}
	defer imq.Release()

	res, err := imq.Quantize(attr)
	if err != nil {
		log.Printf("failed to perform image quantization: %v", err)
		return nil, err
	}
	defer res.Release()

	rgb8data, err := res.WriteRemappedImage()
	if err != nil {
		log.Printf("failed to write remapped image: %v", err)
		return nil, err
	}

	prepImage := RGB8ToImage(res.GetImageWidth(), res.GetImageHeight(), rgb8data, res.GetPalette())

	err = png.Encode(buf, prepImage)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type WebPEncoder struct{}

func NewWebPEncoder() *WebPEncoder {
	return &WebPEncoder{}
}

func (we *WebPEncoder) Encode(imageData []byte, quality int, isLossless bool) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}

	var encodedImage bytes.Buffer
	err = webp.Encode(&encodedImage, img, &webp.Options{Lossless: isLossless, Quality: float32(quality)})
	if err != nil {
		return nil, err
	}

	return encodedImage.Bytes(), nil
}
