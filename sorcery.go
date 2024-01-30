package main

import (
	"image"
	"slices"

	"golang.org/x/image/draw"
)

func resize(imgs []image.Image, width int) []image.Image {

	var resizedImages []image.Image

	for _, img := range imgs {
		// skip resize if already of the width
		if img.Bounds().Max.X == width {
			resizedImages = append(resizedImages, img)
			continue
		}
		height := int(float64(img.Bounds().Max.Y) / float64(img.Bounds().Max.X) * float64(width))
		var newImg draw.Image = image.NewRGBA(image.Rect(0, 0, width, height))

		draw.CatmullRom.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Over, nil)
		resizedImages = append(resizedImages, newImg)
	}

	return resizedImages
}

func merge(imgs []image.Image) image.Image {
	var width, height int = 0, 0
	var mergedImage draw.Image = image.NewRGBA(image.Rect(0, 0, width, height))

	for _, img := range imgs {
		height += img.Bounds().Max.Y
		if img.Bounds().Max.X < width || width == 0 {
			width = img.Bounds().Max.X
		}

		var newImg draw.Image = image.NewRGBA(image.Rect(0, 0, width, height))

		draw.Draw(newImg, mergedImage.Bounds(), mergedImage, image.Pt(0, 0), draw.Src)
		draw.Draw(newImg, image.Rect(0, height-img.Bounds().Max.Y, width, height), img, image.Pt(0, 0), draw.Src)

		mergedImage = newImg
	}

	return mergedImage
}

func slice(img image.Image, roughHeight int, threshold int, skipStep int, neighborCount int, absoluteHeight bool) []image.Image {
	var images []image.Image
	var height, width int = img.Bounds().Max.Y, img.Bounds().Max.X

	// skip slicing
	if roughHeight > height {
		return append(images, img)
	}

	if absoluteHeight {
		// Slice the image into equal-height sections
		for i := 1; i*roughHeight < height; i++ {
			var outImg draw.Image = image.NewRGBA(image.Rect(0, 0, width, roughHeight))
			draw.Draw(outImg, outImg.Bounds(), img, image.Pt(0, (i-1)*roughHeight), draw.Src)
			images = append(images, outImg)
		}
		return images
	} else {
		var skip bool = false
		var prevStep int
		col := roughHeight

		for prevStep < height {
			var pixels []int
			for x := 0; x < img.Bounds().Max.X; x++ {
				pixels = append(pixels, toGray(img.At(x, col)))
			}

			// Check pixel differences in a row to determine slice position
			for r := neighborCount; r < len(pixels); r++ {
				array := pixels[r-neighborCount : r]
				var diff int = slices.Max(array) - slices.Min(array)

				if diff > threshold || diff < -threshold {
					col += skipStep
					skip = true
					break
				}
			}

			if skip {
				skip = false
				continue
			}

			var h int

			// remove blank space & merge bitty slices
			if height-col < roughHeight/4 || col > height {
				h = height
			} else {
				h = col
			}

			// generate the image
			var outImg draw.Image = image.NewRGBA(image.Rect(0, 0, width, h-prevStep))
			draw.Draw(outImg, outImg.Bounds(), img, image.Pt(0, prevStep), draw.Src)
			images = append(images, outImg)

			prevStep = h
			col += roughHeight
		}

		return images
	}
}
