package imageencode

import (
	"image"
	"image/color"
)

// imageToRGBA32 converts an image to an RGBA32 byte slice.
// It takes an image and returns a byte slice in RGBA32 format representing the image.
// The input image is assumed to be in the sRGB color space.
func imageToRGBA32(img image.Image) []byte {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	ret := make([]byte, width*height*4)

	index := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			ret[index] = uint8(r >> 8)
			ret[index+1] = uint8(g >> 8)
			ret[index+2] = uint8(b >> 8)
			ret[index+3] = uint8(a >> 8)
			index += 4
		}
	}

	return ret
}

// RGB8ToImage converts RGB8 data with a color palette to an image.Image.
// It takes the width and height of the image, an RGB8 data slice, and a color palette.
// The function returns an image.Image created from the RGB8 data with the provided palette.
// The input RGB8 data is assumed to be in the sRGB color space.
func RGB8ToImage(width, height int, rgb8data []byte, palette color.Palette) image.Image {
	rect := image.Rect(0, 0, width, height)
	ret := image.NewPaletted(rect, palette)

	index := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ret.SetColorIndex(x, y, rgb8data[index])
			index++
		}
	}

	return ret
}
