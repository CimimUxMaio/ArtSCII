package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CimimUxMaio/artscii/artscii"
	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("ArtSCII", "display or create .artscii files")
	read := parser.NewCommand("read", "displays a given artscii file in the command line")
	create := parser.NewCommand("create", "creates an .artscii file from a given image file")

	artsciiPath := read.String("p", "path", &argparse.Options{Required: true, Help: "the given image's path"})

	imagePath := create.String("p", "image", &argparse.Options{Required: true, Help: "the given .artscii file path"})
	defaultScale := " `-:~*r+=xhwAD9MWB@"
	scale := create.String("s", "scale", &argparse.Options{Required: false, Default: defaultScale, Help: "Ascii characters to use ordered by brightness (from low to high)"})
	outputPath := create.String("o", "name", &argparse.Options{Required: false, Help: "the generated .artscii file path"})

	err := parser.Parse(os.Args)
	if err != nil {
		log.Fatal(parser.Usage(err))
		return
	}

	if read.Happened() {
		runRead(*artsciiPath)
	} else if create.Happened() {
		runCreate(*imagePath, []rune(*scale), *outputPath)
	}
}

func runRead(artsciiPath string) {
	ascii, err := artscii.FromArtSCIIFile(artsciiPath)
	if err != nil {
		log.Fatal(err)
	}

	ascii.Print()
}

func runCreate(imagePath string, scale []rune, outputPath string) {
	fmt.Println("Using scale: [" + string(scale) + "]")
	ascii, err := artscii.FromImageFile(imagePath, scale)
	if err != nil {
		log.Fatal(err)
	}

	finalOutputPath := generateOutputfileName(imagePath)
	if outputPath != "" {
		finalOutputPath = outputPath
	}

	ascii.ToFile(finalOutputPath)
	fmt.Println("Created " + finalOutputPath + ".artscii from " + imagePath)
}

func generateOutputfileName(inputPath string) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
	return "photo_" + string(hash[:8])
}
