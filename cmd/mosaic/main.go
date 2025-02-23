package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/kohge4/mosaic-image"
)

func main() {
	// Parse command line arguments
	input := flag.String("input", "", "Path to input image (required)")
	output := flag.String("output", "", "Path to output image (required)")
	k := flag.Int("k", 8, "Number of colors to use")
	blockSize := flag.Int("block", 10, "Size of mosaic blocks in pixels")
	iterations := flag.Int("iterations", 50, "Maximum number of k-means iterations")
	tolerance := flag.Float64("tolerance", 0.001, "Convergence tolerance for k-means")

	// Region options
	x := flag.Int("x", -1, "X-coordinate of top-left corner for mosaic region (-1 for entire width)")
	y := flag.Int("y", -1, "Y-coordinate of top-left corner for mosaic region (-1 for entire height)")
	width := flag.Int("width", -1, "Width of mosaic region (-1 for remaining width)")
	height := flag.Int("height", -1, "Height of mosaic region (-1 for remaining height)")

	flag.Parse()

	// Check required parameters
	if *input == "" || *output == "" {
		fmt.Println("Error: input and output paths are required")
		flag.Usage()
		os.Exit(1)
	}

	// Open input image
	file, err := os.Open(*input)
	if err != nil {
		fmt.Printf("Error: could not open input image: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error: could not decode image: %v\n", err)
		os.Exit(1)
	}

	// Configure mosaic options
	opts := mosaic.DefaultOptions()
	opts.K = *k
	opts.BlockSize = *blockSize
	opts.Iterations = *iterations
	opts.Tolerance = *tolerance

	// Configure region if specified
	bounds := img.Bounds()
	if *x >= 0 && *y >= 0 {
		regionX := *x
		regionY := *y
		regionWidth := *width
		regionHeight := *height

		// Use remaining width/height if not specified
		if regionWidth < 0 {
			regionWidth = bounds.Max.X - regionX
		}
		if regionHeight < 0 {
			regionHeight = bounds.Max.Y - regionY
		}

		opts.Region = &mosaic.Region{
			X:      regionX,
			Y:      regionY,
			Width:  regionWidth,
			Height: regionHeight,
		}
	}

	// Create output directory if needed
	if err := os.MkdirAll(filepath.Dir(*output), 0755); err != nil {
		fmt.Printf("Error: could not create output directory: %v\n", err)
		os.Exit(1)
	}

	// Generate mosaic image
	mosaicImg := mosaic.CreateMosaic(img, opts)

	// Save result
	outFile, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Error: could not create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, mosaicImg); err != nil {
		fmt.Printf("Error: could not save image: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Mosaic image created successfully:", *output)
}
