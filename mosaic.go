package mosaic

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
)

// Pixel represents a single pixel with RGB values
type Pixel struct {
	R, G, B float64
}

// Region defines the area to apply mosaic effect
type Region struct {
	X      int // x-coordinate of top-left corner
	Y      int // y-coordinate of top-left corner
	Width  int // width of the region
	Height int // height of the region
}

// MosaicOptions contains configuration for mosaic generation
type MosaicOptions struct {
	K          int     // number of colors for k-means
	BlockSize  int     // size of mosaic blocks
	Iterations int     // number of k-means iterations
	Tolerance  float64 // convergence tolerance
	Region     *Region // region to apply mosaic effect (nil for entire image)
}

// DefaultOptions returns default mosaic options
func DefaultOptions() *MosaicOptions {
	return &MosaicOptions{
		K:          8,
		BlockSize:  10,
		Iterations: 50,
		Tolerance:  0.001,
		Region:     nil,
	}
}

// CreateMosaic creates a mosaic image from the input image using k-means clustering
func CreateMosaic(img image.Image, opts *MosaicOptions) image.Image {
	if opts == nil {
		opts = DefaultOptions()
	}

	bounds := img.Bounds()
	region := opts.Region
	if region == nil {
		region = &Region{
			X:      bounds.Min.X,
			Y:      bounds.Min.Y,
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
		}
	}

	// Validate region
	if region.X < bounds.Min.X || region.Y < bounds.Min.Y ||
		region.X+region.Width > bounds.Max.X || region.Y+region.Height > bounds.Max.Y {
		// If invalid region, use entire image
		region = &Region{
			X:      bounds.Min.X,
			Y:      bounds.Min.Y,
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
		}
	}

	// Create output image (copy of original)
	mosaic := image.NewRGBA(bounds)
	draw.Draw(mosaic, bounds, img, bounds.Min, draw.Src)

	// Convert specified region to pixels
	pixels := imageToPixels(img, region)

	// Perform k-means clustering
	centroids := kmeans(pixels, opts.K, opts.Iterations, opts.Tolerance)

	// Process each block within the specified region
	for y := region.Y; y < region.Y+region.Height; y += opts.BlockSize {
		for x := region.X; x < region.X+region.Width; x += opts.BlockSize {
			// Calculate average color for the block
			blockPixels := make([]Pixel, 0)
			for by := 0; by < opts.BlockSize && y+by < region.Y+region.Height; by++ {
				for bx := 0; bx < opts.BlockSize && x+bx < region.X+region.Width; bx++ {
					r, g, b, _ := img.At(x+bx, y+by).RGBA()
					blockPixels = append(blockPixels, Pixel{
						R: float64(r) / 65535,
						G: float64(g) / 65535,
						B: float64(b) / 65535,
					})
				}
			}

			// Find nearest centroid
			avgColor := findNearestCentroid(averagePixels(blockPixels), centroids)

			// Fill block with average color
			blockColor := color.RGBA{
				R: uint8(avgColor.R * 255),
				G: uint8(avgColor.G * 255),
				B: uint8(avgColor.B * 255),
				A: 255,
			}

			fillBlock(mosaic, x, y, opts.BlockSize, blockColor)
		}
	}

	return mosaic
}

// imageToPixels converts a region of an image to a slice of Pixels
func imageToPixels(img image.Image, region *Region) []Pixel {
	pixels := make([]Pixel, 0, region.Width*region.Height)

	for y := region.Y; y < region.Y+region.Height; y++ {
		for x := region.X; x < region.X+region.Width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, Pixel{
				R: float64(r) / 65535,
				G: float64(g) / 65535,
				B: float64(b) / 65535,
			})
		}
	}

	return pixels
}

// kmeans performs k-means clustering on pixels
func kmeans(pixels []Pixel, k, maxIterations int, tolerance float64) []Pixel {
	// Initialize random centroids
	centroids := make([]Pixel, k)
	for i := range centroids {
		idx := rand.Intn(len(pixels))
		centroids[i] = pixels[idx]
	}

	for iteration := 0; iteration < maxIterations; iteration++ {
		// Assign pixels to clusters
		clusters := make([][]Pixel, k)
		for _, p := range pixels {
			nearest := findNearestCentroidIndex(p, centroids)
			clusters[nearest] = append(clusters[nearest], p)
		}

		// Update centroids
		newCentroids := make([]Pixel, k)
		maxDiff := 0.0

		for i := range centroids {
			if len(clusters[i]) > 0 {
				newCentroids[i] = averagePixels(clusters[i])
				diff := distance(centroids[i], newCentroids[i])
				if diff > maxDiff {
					maxDiff = diff
				}
			} else {
				newCentroids[i] = centroids[i]
			}
		}

		centroids = newCentroids

		// Check for convergence
		if maxDiff < tolerance {
			break
		}
	}

	return centroids
}

// findNearestCentroidIndex finds the index of the nearest centroid to a pixel
func findNearestCentroidIndex(p Pixel, centroids []Pixel) int {
	minDist := math.MaxFloat64
	nearest := 0

	for i, c := range centroids {
		dist := distance(p, c)
		if dist < minDist {
			minDist = dist
			nearest = i
		}
	}

	return nearest
}

// findNearestCentroid finds the nearest centroid to a pixel
func findNearestCentroid(p Pixel, centroids []Pixel) Pixel {
	return centroids[findNearestCentroidIndex(p, centroids)]
}

// distance calculates Euclidean distance between two pixels
func distance(p1, p2 Pixel) float64 {
	dr := p1.R - p2.R
	dg := p1.G - p2.G
	db := p1.B - p2.B
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// averagePixels calculates the average color of a slice of pixels
func averagePixels(pixels []Pixel) Pixel {
	if len(pixels) == 0 {
		return Pixel{}
	}

	var sumR, sumG, sumB float64
	for _, p := range pixels {
		sumR += p.R
		sumG += p.G
		sumB += p.B
	}

	n := float64(len(pixels))
	return Pixel{
		R: sumR / n,
		G: sumG / n,
		B: sumB / n,
	}
}

// fillBlock fills a block in the image with a single color
func fillBlock(img *image.RGBA, x, y, size int, c color.Color) {
	bounds := img.Bounds()
	for by := 0; by < size && y+by < bounds.Max.Y; by++ {
		for bx := 0; bx < size && x+bx < bounds.Max.X; bx++ {
			img.Set(x+bx, y+by, c)
		}
	}
}
