package internal

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/nfnt/resize"
)

const TILE_WIDTH = 256
const TILE_HEIGHT = TILE_WIDTH

type TileData struct {
	Name string
	Data []byte
}

func TileImage(img image.Image, tileChan chan<- TileData) error {
	fmt.Println("Tiling concurrently...")
	// The number of workers is typically set to the number of CPU cores.
	const numWorkers = 4 // Adjust based on target Lambda/machine specs

	defer close(tileChan)

	maxZoom := calcZoomLevel(img, TILE_WIDTH, TILE_HEIGHT)
	reducedImage := img

	for z := 0; z <= maxZoom; z++ {
		zFolder := maxZoom - z
		width := reducedImage.Bounds().Dx()
		height := reducedImage.Bounds().Dy()

		// 1. Job Creation Channel and WaitGroup
		// A channel to send tile coordinates as jobs
		tileCoordChan := make(chan struct{ x, y int })
		var wg sync.WaitGroup // Waits for all workers to finish

		// 2. Start the Worker Pool
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go tileWorker(tileCoordChan, tileChan, reducedImage, zFolder, &wg)
		}

		// 3. Submit Jobs (Coordinates)
		tilesAlongX := int(math.Ceil(float64(width) / float64(TILE_WIDTH)))
		tilesAlongY := int(math.Ceil(float64(height) / float64(TILE_HEIGHT)))
		tilesAlongAxis := int(math.Max(float64(tilesAlongX), float64(tilesAlongY)))

		for x := 0; x < tilesAlongAxis; x++ {
			for y := 0; y < tilesAlongAxis; y++ {
				// Non-blocking submission to the channel
				tileCoordChan <- struct{ x, y int }{x, y}
			}
		}

		// 4. Close Job Channel and Wait
		close(tileCoordChan) // Signal to workers that no more jobs are coming
		wg.Wait()            // Block until all workers finish all jobs for this zoom level

		// Error handling needs to be robust here (see notes below)

		// 5. Downscale for next zoom level (Sequential)
		if z < maxZoom {
			newWidth := uint(width / 2)
			newHeight := uint(height / 2)
			// This step MUST remain sequential, as it creates the image for the next loop iteration.
			reducedImage = resize.Resize(newWidth, newHeight, reducedImage, resize.Lanczos3)
		}
	}
	return nil
}

// tileWorker is the goroutine function that handles the cropping and encoding job.
func tileWorker(jobChan <-chan struct{ x, y int }, resultChan chan<- TileData, img image.Image, zFolder int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Pre-create the PNG encoder for optimal reuse
	var encoder png.Encoder
	encoder.CompressionLevel = png.BestSpeed // Optimization 2: Faster PNG compression

	for job := range jobChan {
		x, y := job.x, job.y

		x0 := x * TILE_WIDTH
		y0 := y * TILE_HEIGHT
		x1 := x0 + TILE_WIDTH
		y1 := y0 + TILE_HEIGHT

		// 1. Crop the tile (Optimization 1: Use the faster image/draw implementation)
		tileImage := npCropImageOptimized(img, x0, y0, x1, y1)

		// 2. Encode to PNG
		var buf bytes.Buffer

		if err := encoder.Encode(&buf, tileImage); err != nil {
			// Error handling: The worker can't return error directly.
			// It must send it back via a separate error channel or handle it internally.
			fmt.Printf("Error encoding tile %d/%d_%d: %v\n", zFolder, x, y, err)
			continue // Skip this tile
		}

		tileName := filepath.Join(strconv.Itoa(zFolder), fmt.Sprintf("%d_%d.png", x, y))
		resultChan <- TileData{
			Name: tileName,
			Data: buf.Bytes(),
		}
	}
}

// npCropImageOptimized crops the image slice and efficiently pads it with
// transparency if the requested area extends out of bounds.
func npCropImageOptimized(img image.Image, x0, y0, x1, y1 int) *image.RGBA {
	width := x1 - x0
	height := y1 - y0

	// Create the destination image (the tile). image.NewRGBA initializes
	// all pixels to transparent black (color.RGBA{0, 0, 0, 0}), serving as the padding.
	croppedImg := image.NewRGBA(image.Rect(0, 0, width, height))

	// 1. Define the source rectangle *requested* in the original image coordinates.
	srcRect := image.Rect(x0, y0, x1, y1)

	// 2. Clamp the source rectangle to the original image's bounds.
	// This gives us the rectangle of actual data available in 'img'.
	srcClamped := srcRect.Intersect(img.Bounds())

	// 3. Define the destination rectangle. This shifts the clamped source
	// coordinates so they start at (0, 0) within the new 'croppedImg'.
	dstClamped := srcClamped.Sub(srcRect.Min)

	// 4. Use draw.Draw to perform the optimized block copy.
	// - croppedImg: The destination image.
	// - dstClamped: Where in the destination to draw the pixels.
	// - img: The source image.
	// - srcClamped.Min: The top-left point in the source image to start copying from.
	// - draw.Src: The drawing operator.
	draw.Draw(croppedImg, dstClamped, img, srcClamped.Min, draw.Src)

	// Pixels in croppedImg outside dstClamped retain the transparent padding from initialization.
	return croppedImg
}

// calcZoomLevel calculates the max zoom level based on image size and tile size.
func calcZoomLevel(img image.Image, mapTileWidth, mapTileHeight int) int {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	wz := int(math.Ceil(math.Log2(float64(height) / float64(mapTileHeight))))
	hz := int(math.Ceil(math.Log2(float64(width) / float64(mapTileWidth))))

	return int(math.Max(float64(wz), float64(hz)))
}
