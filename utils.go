package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"sync"

	"manas140/seam/out"
)

func col(text string, col int) string {
	return fmt.Sprintf("\033[0;%dm%s\033[1;0m", col, text)
}

func toGray(c color.Color) int {
	r, g, b, _ := c.RGBA()
	return int((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
}

func dirExists(dir string) bool {
	if _, err := os.Open(dir); err == nil {
		return true
	} else {
		return false
	}
}

func readImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, nil
}

func saveImage(img image.Image, path string, ext string, wg *sync.WaitGroup) {
	file, err := os.Create(path)
	if err != nil {
		out.Error("Unable to create file '%s'", path)
	}
	defer file.Close()

	switch ext {
	case "png":
		png.Encode(file, img)
	default:
		jpeg.Encode(file, img, nil)
	}
	wg.Done()
}
