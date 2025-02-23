# Mosaic Image Generator
※Clineにより生成したコード

このGoライブラリは、k-means法を用いて画像をモザイク化するためのツールです。入力画像をブロックに分割し、k-means法による色の量子化を適用してモザイク効果を生成します。

ライブラリとしての利用とCLIツールとしての利用の両方に対応しています。

## 特徴

- k-means法による色の量子化
- モザイクブロックサイズの設定が可能
- 使用する色数（kの値）の調整が可能
- PNG、JPEG形式の画像をサポート

## インストール

### ライブラリとして
```bash
go get github.com/kohge4/mosaic-image
```

### CLIツールとして
```bash
go install github.com/kohge4/mosaic-image/cmd/mosaic@latest
```

## 使用方法

### CLIツールとして

```bash
# 基本的な使用方法
mosaic -input input.png -output output.png

# オプションを指定
mosaic -input input.png -output output.png -k 16 -block 20

# 利用可能なオプション
mosaic -help
```

利用可能なオプション：
- `-input`: 入力画像のパス（必須）
- `-output`: 出力画像のパス（必須）
- `-k`: 使用する色数（デフォルト: 8）
- `-block`: モザイクブロックのサイズ（デフォルト: 10）
- `-iterations`: k-meansの最大イテレーション回数（デフォルト: 50）
- `-tolerance`: k-meansの収束判定の閾値（デフォルト: 0.001）

### ライブラリとして

基本的な使用例：

```go
package main

import (
    "image"
    "image/png"
    "os"

    "github.com/kohge4/mosaic-image"
)

func main() {
    // 画像を開いてデコード
    file, _ := os.Open("input.png")
    img, _, _ := image.Decode(file)
    file.Close()

    // オプションを設定
    opts := mosaic.DefaultOptions()
    opts.K = 8          // 使用する色数
    opts.BlockSize = 10 // モザイクブロックのサイズ（ピクセル）

    // モザイク画像を生成
    mosaicImg := mosaic.CreateMosaic(img, opts)

    // 結果を保存
    outFile, _ := os.Create("output.png")
    png.Encode(outFile, mosaicImg)
    outFile.Close()
}
```

## 設定オプション

`MosaicOptions`構造体で以下の設定をカスタマイズできます：

- `K`: 使用する色数（デフォルト: 8）
- `BlockSize`: モザイクブロックのサイズ（ピクセル単位、デフォルト: 10）
- `Iterations`: k-meansの最大イテレーション回数（デフォルト: 50）
- `Tolerance`: k-meansの収束判定の閾値（デフォルト: 0.001）

`DefaultOptions()`を使用してデフォルト設定を取得し、必要に応じて変更することができます。

## ライセンス

MIT License

## 貢献

プルリクエストや提案は大歓迎です！
