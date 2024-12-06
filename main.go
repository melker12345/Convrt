package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/nfnt/resize"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/image/tiff"
)

type ConversionOptions struct {
	Quality    int  // JPEG quality (1-100)
	Optimize   bool // Whether to optimize the image
	Width      int  // New width (0 means keep original)
	Height     int  // New height (0 means keep original)
	Recursive  bool // Process directories recursively
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse options
	options := parseOptions(os.Args)

	// Get input paths (can be files or directories)
	inputs := getInputPaths(os.Args[1])
	
	// Process each input
	for _, input := range inputs {
		if err := processInput(input, options); err != nil {
			color.Red("Error processing %s: %v", input, err)
		}
	}
}

func processInput(path string, options ConversionOptions) error {
	// Check if it's a directory
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return processDirectory(path, options)
	}
	return processFile(path, options)
}

func processDirectory(dir string, options ConversionOptions) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return processFile(path, options)
		}
		return nil
	})
}

func processFile(inputFile string, options ConversionOptions) error {
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file '%s' does not exist", inputFile)
	}

	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	defer file.Close()

	// Create progress bar
	bar := progressbar.DefaultBytes(
		-1,
		"Processing",
	)

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding image: %v", err)
	}

	// Resize if requested
	if options.Width > 0 || options.Height > 0 {
		color.Blue("Resizing image to %dx%d", options.Width, options.Height)
		img = resize(img, options.Width, options.Height)
	}

	var outputFile string
	var outputFormat string

	if len(os.Args) >= 3 && strings.HasPrefix(os.Args[2], ".") {
		outputFormat = strings.TrimPrefix(os.Args[2], ".")
		outputFile = getConvertedFilename(inputFile, os.Args[2])
		err = saveWithFormat(img, outputFile, outputFormat, options, bar)
	} else if options.Optimize {
		outputFormat = format
		outputFile = getOptimizedFilename(inputFile)
		err = saveWithOptimization(img, outputFile, format, options, bar)
	} else {
		return fmt.Errorf("invalid command. Use -h for help")
	}

	if err != nil {
		return err
	}

	printSuccess(inputFile, outputFile)
	return nil
}

func resize(img image.Image, width, height int) image.Image {
	// Get original dimensions
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// If neither width nor height is specified, return original
	if width == 0 && height == 0 {
		return img
	}

	// Calculate new dimensions while maintaining aspect ratio
	var newWidth, newHeight uint
	if width == 0 {
		// Height specified, calculate width
		ratio := float64(height) / float64(originalHeight)
		newHeight = uint(height)
		newWidth = uint(float64(originalWidth) * ratio)
	} else if height == 0 {
		// Width specified, calculate height
		ratio := float64(width) / float64(originalWidth)
		newWidth = uint(width)
		newHeight = uint(float64(originalHeight) * ratio)
	} else {
		// Both specified
		newWidth = uint(width)
		newHeight = uint(height)
	}

	// Use Lanczos3 resampling for high-quality resizing
	return resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
}

func parseOptions(args []string) ConversionOptions {
	options := ConversionOptions{
		Quality:    90,
		Optimize:   false,
		Width:      0,
		Height:     0,
		Recursive:  false,
	}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-o", "--optimize":
			options.Optimize = true
		case "-q", "--quality":
			if i+1 < len(args) {
				quality := 0
				if _, err := fmt.Sscanf(args[i+1], "%d", &quality); err == nil {
					options.Quality = quality
					i++
				}
			}
		case "-w", "--width":
			if i+1 < len(args) {
				width := 0
				if _, err := fmt.Sscanf(args[i+1], "%d", &width); err == nil {
					options.Width = width
					i++
				}
			}
		case "-h", "--height":
			if i+1 < len(args) {
				height := 0
				if _, err := fmt.Sscanf(args[i+1], "%d", &height); err == nil {
					options.Height = height
					i++
				}
			}
		case "-r", "--recursive":
			options.Recursive = true
		}
	}

	// Validate quality
	if options.Quality < 1 {
		options.Quality = 1
	} else if options.Quality > 100 {
		options.Quality = 100
	}

	// Validate dimensions
	if options.Width < 0 {
		options.Width = 0
	}
	if options.Height < 0 {
		options.Height = 0
	}

	return options
}

func getInputPaths(input string) []string {
	// Handle wildcard patterns
	if strings.Contains(input, "*") {
		matches, err := filepath.Glob(input)
		if err == nil && len(matches) > 0 {
			return matches
		}
	}
	return []string{input}
}

func saveWithFormat(img image.Image, outputPath, format string, options ConversionOptions, bar *progressbar.ProgressBar) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := io.MultiWriter(out, bar)

	switch format {
	case "jpg", "jpeg":
		return jpeg.Encode(writer, img, &jpeg.Options{Quality: options.Quality})
	case "png":
		return png.Encode(writer, img)
	case "gif":
		return gif.Encode(writer, img, nil)
	case "tiff":
		return tiff.Encode(writer, img, &tiff.Options{Compression: tiff.Deflate})
	case "webp":
		return webp.Encode(writer, img, &encoder.Options{
			Lossless: false,
			Quality:  float32(options.Quality),
		})
	default:
		return fmt.Errorf("unsupported format: %s (supported formats: .jpg, .jpeg, .png, .gif, .tiff, .webp)", format)
	}
}

func saveWithOptimization(img image.Image, outputPath, format string, options ConversionOptions, bar *progressbar.ProgressBar) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	writer := io.MultiWriter(out, bar)

	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(writer, img, &jpeg.Options{Quality: 75})
	case "png":
		return png.Encode(writer, img)
	case "gif":
		return gif.Encode(writer, img, nil)
	case "tiff":
		return tiff.Encode(writer, img, &tiff.Options{Compression: tiff.Deflate})
	case "webp":
		return webp.Encode(writer, img, &encoder.Options{
			Lossless: false,
			Quality:  75.0,
		})
	default:
		return fmt.Errorf("unsupported format: %s (supported formats: jpg, jpeg, png, gif, tiff, webp)", format)
	}
}

func printUsage() {
	color.Blue("Usage: convrt [input] [.format] [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -o, --optimize     Optimize the image")
	fmt.Println("  -q, --quality N    Set quality (1-100, default: 90)")
	fmt.Println("  -w, --width N      Resize to width N")
	fmt.Println("  -h, --height N     Resize to height N")
	fmt.Println("  -r, --recursive    Process directories recursively")
	fmt.Println("\nExamples:")
	fmt.Println("  convrt image.jpg .png              # Convert to PNG")
	fmt.Println("  convrt image.png .jpg -q 85        # Convert to JPEG with quality 85")
	fmt.Println("  convrt image.jpg -o                # Optimize the image")
	fmt.Println("  convrt image.jpg .webp -w 800      # Convert to WebP and resize width to 800px")
	fmt.Println("  convrt images/*.jpg .webp          # Convert all JPGs to WebP")
	fmt.Println("  convrt images/ -o -r               # Optimize all images in directory recursively")
}

func getConvertedFilename(inputFile, newExt string) string {
	ext := filepath.Ext(inputFile)
	return strings.TrimSuffix(inputFile, ext) + newExt
}

func getOptimizedFilename(inputFile string) string {
	ext := filepath.Ext(inputFile)
	return strings.TrimSuffix(inputFile, ext) + "_optimized" + ext
}

func printSuccess(inputFile, outputFile string) {
	inputInfo, _ := os.Stat(inputFile)
	outputInfo, _ := os.Stat(outputFile)
	inputSize := float64(inputInfo.Size()) / 1024  // KB
	outputSize := float64(outputInfo.Size()) / 1024 // KB
	reduction := ((inputSize - outputSize) / inputSize) * 100

	color.Green("\n✓ Successfully processed image!")
	fmt.Printf("Input:  %s (%.2f KB)\n", inputFile, inputSize)
	fmt.Printf("Output: %s (%.2f KB)\n", outputFile, outputSize)
	if reduction > 0 {
		color.Green("Size reduced by %.2f%%", reduction)
	}
}
