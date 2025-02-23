package mosaic

import (
	"image"
	"image/color"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"K value", opts.K, 8},
		{"BlockSize value", opts.BlockSize, 10},
		{"Iterations value", opts.Iterations, 50},
		{"Tolerance value", opts.Tolerance, 0.001},
		{"Region value", opts.Region, (*Region)(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("DefaultOptions() %s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestCreateMosaic(t *testing.T) {
	// Create test image
	width, height := 100, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill image with two colors
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x < width/2 {
				img.Set(x, y, color.RGBA{R: 255, A: 255}) // Red
			} else {
				img.Set(x, y, color.RGBA{B: 255, A: 255}) // Blue
			}
		}
	}

	// Test mosaic creation with specific region
	region := &Region{
		X:      width / 4,
		Y:      height / 4,
		Width:  width / 2,
		Height: height / 2,
	}

	opts := &MosaicOptions{
		K:          2,  // 2 colors
		BlockSize:  10, // 10x10 blocks
		Iterations: 10, // 10 iterations
		Tolerance:  0.001,
		Region:     region,
	}

	result := CreateMosaic(img, opts)
	bounds := result.Bounds()

	// Verify image dimensions
	if bounds.Dx() != width || bounds.Dy() != height {
		t.Errorf("Invalid image dimensions. got = %dx%d, want = %dx%d",
			bounds.Dx(), bounds.Dy(), width, height)
	}

	// Check colors in the mosaic region
	centerX := region.X + region.Width/2
	centerY := region.Y + region.Height/2
	color := result.At(centerX, centerY)
	r, _, _, _ := color.RGBA()

	if r == 0 {
		t.Error("Expected color separation not found in mosaic region")
	}
}

func TestPixelOperations(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"Distance calculation", testDistance},
		{"Average pixels", testAveragePixels},
		{"Find nearest centroid", testFindNearestCentroid},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}

func testDistance(t *testing.T) {
	p1 := Pixel{R: 1.0, G: 0.0, B: 0.0}
	p2 := Pixel{R: 0.0, G: 1.0, B: 0.0}

	expected := 1.4142135623730951 // sqrt(2)
	got := distance(p1, p2)

	if got != expected {
		t.Errorf("distance() = %v, want %v", got, expected)
	}
}

func testAveragePixels(t *testing.T) {
	pixels := []Pixel{
		{R: 1.0, G: 0.0, B: 0.0},
		{R: 0.0, G: 1.0, B: 0.0},
	}

	expected := Pixel{R: 0.5, G: 0.5, B: 0.0}
	got := averagePixels(pixels)

	if got != expected {
		t.Errorf("averagePixels() = %v, want %v", got, expected)
	}
}

func testFindNearestCentroid(t *testing.T) {
	centroids := []Pixel{
		{R: 1.0, G: 0.0, B: 0.0}, // Red
		{R: 0.0, G: 0.0, B: 1.0}, // Blue
	}

	testPixel := Pixel{R: 0.9, G: 0.0, B: 0.1} // Color close to red

	got := findNearestCentroid(testPixel, centroids)
	expected := centroids[0] // Should be closest to first centroid (red)

	if got != expected {
		t.Errorf("findNearestCentroid() = %v, want %v", got, expected)
	}
}

func TestImageToPixels(t *testing.T) {
	// Create test image
	width, height := 2, 2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with solid color
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	region := &Region{
		X:      0,
		Y:      0,
		Width:  width,
		Height: height,
	}

	pixels := imageToPixels(img, region)

	// Verify pixel count
	expectedLen := width * height
	if len(pixels) != expectedLen {
		t.Errorf("imageToPixels() returned %d pixels, want %d", len(pixels), expectedLen)
	}

	// Verify pixel values
	expectedPixel := Pixel{R: 1.0, G: 0.0, B: 0.0}
	for i, pixel := range pixels {
		if pixel != expectedPixel {
			t.Errorf("pixel[%d] = %v, want %v", i, pixel, expectedPixel)
		}
	}
}
