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
	// コマンドライン引数の解析
	input := flag.String("input", "", "入力画像のパス（必須）")
	output := flag.String("output", "", "出力画像のパス（必須）")
	k := flag.Int("k", 8, "使用する色数")
	blockSize := flag.Int("block", 10, "モザイクブロックのサイズ（ピクセル）")
	iterations := flag.Int("iterations", 50, "k-meansの最大イテレーション回数")
	tolerance := flag.Float64("tolerance", 0.001, "k-meansの収束判定の閾値")

	flag.Parse()

	// 必須パラメータのチェック
	if *input == "" || *output == "" {
		fmt.Println("Error: input と output は必須パラメータです")
		flag.Usage()
		os.Exit(1)
	}

	// 入力画像の読み込み
	file, err := os.Open(*input)
	if err != nil {
		fmt.Printf("Error: 入力画像を開けません: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// 画像のデコード
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error: 画像のデコードに失敗しました: %v\n", err)
		os.Exit(1)
	}

	// モザイクオプションの設定
	opts := mosaic.DefaultOptions()
	opts.K = *k
	opts.BlockSize = *blockSize
	opts.Iterations = *iterations
	opts.Tolerance = *tolerance

	// モザイク画像の生成
	mosaicImg := mosaic.CreateMosaic(img, opts)

	// 出力ディレクトリの作成
	if err := os.MkdirAll(filepath.Dir(*output), 0755); err != nil {
		fmt.Printf("Error: 出力ディレクトリの作成に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// 結果の保存
	outFile, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Error: 出力ファイルの作成に失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, mosaicImg); err != nil {
		fmt.Printf("Error: 画像の保存に失敗しました: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("モザイク画像を生成しました:", *output)
}
