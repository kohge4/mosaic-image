package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/kohge4/mosaic-image"
)

func main() {
	// Open the input image
	file, err := os.Open("input.png")
	if err != nil {
		fmt.Printf("Error opening image: %v\n", err)
		return
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		return
	}

	// Create mosaic options
	opts := mosaic.DefaultOptions()
	opts.K = 8          // number of colors
	opts.BlockSize = 10 // size of mosaic blocks

	// Create mosaic
	mosaicImg := mosaic.CreateMosaic(img, opts)

	// Save the result
	outFile, err := os.Create("output.png")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outFile.Close()

	if err := png.Encode(outFile, mosaicImg); err != nil {
		fmt.Printf("Error encoding output image: %v\n", err)
		return
	}

	fmt.Println("Mosaic image created successfully!")
}
