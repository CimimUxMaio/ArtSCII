package main

import (
	"log"

	"github.com/CimimUxMaio/artscii/artscii"
)

func main() {
	ascii, err := artscii.FromArtSCIIFile("generated-example.artscii")
	if err != nil {
		log.Fatal(err)
	}

	ascii.Print()
}
