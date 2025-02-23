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
	// テスト用の画像を作成
	width, height := 100, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 画像を2色で塗り分ける
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x < width/2 {
				img.Set(x, y, color.RGBA{R: 255, A: 255}) // 赤
			} else {
				img.Set(x, y, color.RGBA{B: 255, A: 255}) // 青
			}
		}
	}

	// モザイク化のテスト
	opts := &MosaicOptions{
		K:          2,  // 2色
		BlockSize:  10, // 10x10ブロック
		Iterations: 10, // 10回イテレーション
		Tolerance:  0.001,
	}

	result := CreateMosaic(img, opts)
	bounds := result.Bounds()

	// 結果の検証
	if bounds.Dx() != width || bounds.Dy() != height {
		t.Errorf("画像サイズが異なります。got = %dx%d, want = %dx%d",
			bounds.Dx(), bounds.Dy(), width, height)
	}

	// 左半分と右半分で色が異なることを確認
	leftColor := result.At(width/4, height/2)
	rightColor := result.At(3*width/4, height/2)

	lr, _, _, _ := leftColor.RGBA()
	_, _, rb, _ := rightColor.RGBA()

	if lr == 0 || rb == 0 {
		t.Error("期待される色の分離が行われていません")
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
		{R: 1.0, G: 0.0, B: 0.0}, // 赤
		{R: 0.0, G: 0.0, B: 1.0}, // 青
	}

	testPixel := Pixel{R: 0.9, G: 0.0, B: 0.1} // 赤に近い色

	got := findNearestCentroid(testPixel, centroids)
	expected := centroids[0] // 最初のcentroid（赤）が最も近いはず

	if got != expected {
		t.Errorf("findNearestCentroid() = %v, want %v", got, expected)
	}
}

func TestImageToPixels(t *testing.T) {
	// テスト用の画像を作成
	width, height := 2, 2
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 単色で塗りつぶし
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	pixels := imageToPixels(img)

	// ピクセル数の確認
	expectedLen := width * height
	if len(pixels) != expectedLen {
		t.Errorf("imageToPixels() returned %d pixels, want %d", len(pixels), expectedLen)
	}

	// 各ピクセルの値を確認
	expectedPixel := Pixel{R: 1.0, G: 0.0, B: 0.0}
	for i, pixel := range pixels {
		if pixel != expectedPixel {
			t.Errorf("pixel[%d] = %v, want %v", i, pixel, expectedPixel)
		}
	}
}
