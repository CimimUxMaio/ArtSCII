package artscii

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"io"
	"math"
	"os"
)

type AsciiImage struct {
	scale  []rune
	pixels [][]rune
}

func FromImage(img image.Image, scale []rune) *AsciiImage {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var pixels [][]rune
	for y := 0; y < height; y++ {
		var row []rune
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

func FromImageFile(filePath string, scale []rune) (*AsciiImage, error) {
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

var ErrBadFormat = errors.New("the given file is not of type .ascii or its corrupted")

func FromArtSCIIFile(filePath string) (*AsciiImage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buff, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	content := string(buff)
	var rows [][]rune
	var row []rune
	for _, c := range content {
		if c == '\n' {
			rows = append(rows, row)
			row = []rune{}
			continue
		}
		row = append(row, c)
	}

	if len(rows) < 1 {
		return nil, ErrBadFormat
	}

	scale := rows[0]
	return &AsciiImage{scale: scale, pixels: rows[1:]}, nil
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

func (ascii *AsciiImage) ToFile(path string) (*os.File, error) {
	filePath := path + ".artscii"

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	var content []byte = []byte(string(ascii.scale) + "\n")
	for _, row := range ascii.pixels {
		content = append(content, []byte(string(row)+"\n")...)
	}
	_, err = file.Write(content)
	if err != nil {
		return nil, err
	}

	return file, nil
}
