package artscii

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
)

type AsciiImage struct {
	scale  []byte
	pixels [][]byte
}

func FromImage(img image.Image, scale []byte) *AsciiImage {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var pixels [][]byte
	for y := 0; y < height; y++ {
		var row []byte
		for x := 0; x < width; x++ {
			gray := color.GrayModel.Convert(img.At(x, y))
			r, _, _, _ := gray.RGBA()
			index := uint8(math.Round(float64(int(r)*(len(scale)-1)) / 0xffff))
			row = append(row, scale[index])
		}
		pixels = append(pixels, row)
	}

	return &AsciiImage{scale: scale, pixels: pixels}
}

func FromFile(filePath string, scale []byte) (*AsciiImage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return FromImage(image, scale), nil
}

func (ascii *AsciiImage) Dim() (width, height int) {
	width = 0
	height = len(ascii.pixels)
	if height > 0 {
		width = len(ascii.pixels[0])
	}

	return width, height
}

func (ascii *AsciiImage) Print() {
	for _, row := range ascii.pixels {
		fmt.Println(string(row))
	}
}
