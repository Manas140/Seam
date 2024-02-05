package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"time"

	"manas140/seam/out"
)

func checkArgs(text string, args []string, val *int, i int) int {
	if i+1 < len(args) {
		v, err := strconv.Atoi(args[i+1])
		if err == nil && v >= 1 {
			*val = v
			return i + 1
		}
	}
	out.Warn("Invalid %s value, using default %d", text, val)
	return i
}

func main() {
	start := time.Now()
	args := os.Args[1:]

	if len(args) == 0 {
		help()
	}

	// create arguments
	var (
		inputDir       string
		outputDir      string = "output"
		roughHeight    int    = 5000
		threshold      int    = 50
		absoluteHeight bool   = false
		format         string = "jpeg"
		skipStep       int    = 5
		neighborCount  int    = 5
		quality        int    = 100
	)

	// get arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h":
			help()

		case "-i":
			if i+1 >= len(args) || args[i+1] == "" {
				out.Fatal("Invalid input directory")
			}
			i++
			inputDir = args[i]
		case "-o":
			if i+1 >= len(args) || args[i+1] == "" {
				out.Warn("Invalid Output Directory, using default '%s'", outputDir)
				continue
			}
			i++
			outputDir = args[i]
		case "-f":
			if i+1 >= len(args) || args[i+1] == "" || !slices.Contains([]string{"jpg", "jpeg", "png"}, args[i+1]) {
				out.Warn("Invalid file format, using default '%s'", format)
				continue
			}
			i++
			format = args[i]
		case "-q":
			i = checkArgs("quality", args, &quality, i)
		case "-s":
			i = checkArgs("skip step", args, &skipStep, i)
		case "-n":
			i = checkArgs("neighbor count", args, &neighborCount, i)
		case "-r":
			i = checkArgs("rough height", args, &roughHeight, i)
		case "-t":
			i = checkArgs("threshold", args, &threshold, i)
		case "-a":
			absoluteHeight = true
		}
	}

	// check if arguments are valid
	if !dirExists(inputDir) {
		out.Fatal("Invalid input directory")
	}

	if filepath.Base(outputDir) == filepath.Base(inputDir) {
		out.Fatal("Input and output directory can't be the same")
	}

	if dirExists(outputDir) {
		files, _ := os.ReadDir(outputDir)
		if len(files) > 0 {
			out.Fatal("Directory '%s' already exists and is not empty", outputDir)
		}
	}

	// start processing the images
	out.Debug("Looking for images in '%s'", inputDir)

	files, err := os.ReadDir(inputDir)
	if err != nil {
		out.Fatal("Unable to read directory: %s", err)
	}

	var imageFiles []string
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".jpg" || ext == ".png" || ext == ".jpeg" {
			imageFiles = append(imageFiles, file.Name())
		}
	}

	if err != nil {
		out.Fatal("Unable to read directory '%s': %s", inputDir, err)
	} else if len(imageFiles) <= 0 {
		out.Fatal("No images found in directory '%s'", inputDir)
	} else {
		out.Debug("%d Images detected", len(imageFiles))
	}

	out.Debug("Creating output directory named '%s'", outputDir)
	err1 := os.MkdirAll(outputDir, 0755)
	if err != nil {
		out.Fatal("Unable to create output directory: %s", err1)
	}

	var images []image.Image
	var w int = 0

	out.Print("Loading images in memory")
	for _, file := range imageFiles {
		img, err := readImage(filepath.Join(inputDir, file))
		if err != nil {
			out.Error("Unable to read image (Skipping): %s", err)
			continue
		}

		if img.Bounds().Max.X < w || w == 0 {
			w = img.Bounds().Max.X
		}

		images = append(images, img)
	}

	out.Print("Resizing images")
	var resizedImg []image.Image = resize(images, w)

	// merge all
	out.Print("Merging images")
	var mergedImg image.Image = merge(resizedImg)

	//slice
	out.Print("Slicing images")
	var slicedImg []image.Image = slice(mergedImg, roughHeight, threshold, skipStep, neighborCount, absoluteHeight)

	//save
	out.Print("Saving %d images", len(slicedImg))

	var wg sync.WaitGroup
	wg.Add(len(slicedImg))
	for i, img := range slicedImg {
		go saveImage(img, fmt.Sprintf("%s/%d.%s", outputDir, i, format), format, quality, &wg)
	}
	wg.Wait()
	out.Debug("Took %.2fs", time.Since(start).Seconds())
}
