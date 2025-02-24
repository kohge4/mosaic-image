# Mosaic Image Generator

[![Test](https://github.com/kohge4/mosaic-image/actions/workflows/test.yml/badge.svg)](https://github.com/kohge4/mosaic-image/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kohge4/mosaic-image)](https://goreportcard.com/report/github.com/kohge4/mosaic-image)
[![Go Reference](https://pkg.go.dev/badge/github.com/kohge4/mosaic-image.svg)](https://pkg.go.dev/github.com/kohge4/mosaic-image)

**Created by Cline(Clineにより作成)**

A Go library for creating mosaic images using k-means clustering. This tool processes input images by dividing them into blocks and applying k-means clustering for color quantization to create a mosaic effect.

Can be used both as a library and as a CLI tool.

## Features

- Color quantization using k-means clustering
- Configurable mosaic block size
- Adjustable number of colors (k value)
- Support for PNG and JPEG images
- Selective region mosaic processing

## Installation

### As a Library
```bash
go get github.com/kohge4/mosaic-image
```

### As a CLI Tool
```bash
go install github.com/kohge4/mosaic-image/cmd/mosaic@latest
```

## Usage

### CLI Tool

```bash
# Basic usage
mosaic -input input.png -output output.png

# With custom options
mosaic -input input.png -output output.png -k 16 -block 20

# Show available options
mosaic -help
```

Available options:
- `-input`: Path to input image (required)
- `-output`: Path to output image (required)
- `-k`: Number of colors to use (default: 8)
- `-block`: Size of mosaic blocks in pixels (default: 10)
- `-iterations`: Maximum number of k-means iterations (default: 50)
- `-tolerance`: Convergence tolerance for k-means (default: 0.001)

Region options:
- `-x`: X-coordinate of top-left corner for mosaic region (-1 for entire width)
- `-y`: Y-coordinate of top-left corner for mosaic region (-1 for entire height)
- `-width`: Width of mosaic region (-1 for remaining width)
- `-height`: Height of mosaic region (-1 for remaining height)

Example with region:
```bash
# Apply mosaic effect to a 200x200 region starting at (100,100)
mosaic -input input.png -output output.png -x 100 -y 100 -width 200 -height 200
```

### As a Library

Basic usage:

```go
package main

import (
    "image"
    "image/png"
    "os"

    "github.com/kohge4/mosaic-image"
)

func main() {
    // Open and decode image
    file, _ := os.Open("input.png")
    img, _, _ := image.Decode(file)
    file.Close()

    // Configure options
    opts := mosaic.DefaultOptions()
    opts.K = 8          // number of colors
    opts.BlockSize = 10 // size of mosaic blocks in pixels

    // Create mosaic
    mosaicImg := mosaic.CreateMosaic(img, opts)

    // Save result
    outFile, _ := os.Create("output.png")
    png.Encode(outFile, mosaicImg)
    outFile.Close()
}
```

## Configuration Options

The `MosaicOptions` struct allows customization of the following settings:

- `K`: Number of colors to use (default: 8)
- `BlockSize`: Size of mosaic blocks in pixels (default: 10)
- `Iterations`: Maximum number of k-means iterations (default: 50)
- `Tolerance`: Convergence tolerance for k-means (default: 0.001)
- `Region`: Region to apply mosaic effect (nil for entire image)
  - `X`: X-coordinate of top-left corner
  - `Y`: Y-coordinate of top-left corner
  - `Width`: Width of the region
  - `Height`: Height of the region

Use `DefaultOptions()` to get default settings and modify them as needed.

## Example with Region

```go
package main

import (
    "image"
    "image/png"
    "os"

    "github.com/kohge4/mosaic-image"
)

func main() {
    // Open and decode image
    file, _ := os.Open("input.png")
    img, _, _ := image.Decode(file)
    file.Close()

    // Configure options with region
    opts := mosaic.DefaultOptions()
    opts.Region = &mosaic.Region{
        X:      100,
        Y:      100,
        Width:  200,
        Height: 200,
    }

    // Create mosaic (only affects specified region)
    mosaicImg := mosaic.CreateMosaic(img, opts)

    // Save result
    outFile, _ := os.Create("output.png")
    png.Encode(outFile, mosaicImg)
    outFile.Close()
}
```

## Contributing

Pull requests and suggestions are welcome! This project uses GitHub Actions for continuous integration and Dependabot for dependency management.

### Development Workflow

1. Fork the repository
2. Create a new branch for your feature/fix
3. Make your changes
4. Run tests locally: `go test -v ./...`
5. Run linter: `golangci-lint run`
6. Submit a pull request

### Continuous Integration

- All pull requests are automatically tested using GitHub Actions
- Tests are run on multiple Go versions (1.23)
- Code quality is checked using golangci-lint
- Dependencies are automatically kept up to date using Dependabot

### Code Quality

- Follow Go best practices and style guidelines
- Write tests for new features
- Keep dependencies up to date
- Document public APIs
