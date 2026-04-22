package converter

import (
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

func init() {
	// png conversions
	Register("png", "jpg", convertImage)
	Register("png", "gif", convertImage)
	Register("png", "webp", convertToWebp)

	// jpg conversions
	Register("jpg", "png", convertImage)
	Register("jpg", "gif", convertImage)
	Register("jpg", "webp", convertToWebp)
	Register("jpeg", "png", convertImage)
	Register("jpeg", "gif", convertImage)
	Register("jpeg", "webp", convertToWebp)

	// gif conversions
	Register("gif", "png", convertImage)
	Register("gif", "jpg", convertImage)

	// webp conversions
	Register("webp", "png", convertFromWebp)
	Register("webp", "jpg", convertFromWebp)
}

// convertImage handles conversions between png, jpg, and gif using the imaging library.
func convertImage(input, output string) error {
	img, err := imaging.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	return imaging.Save(img, output)
}

// convertToWebp encodes an image to webp format.
func convertToWebp(input, output string) error {
	img, err := imaging.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	return webp.Encode(f, img, &webp.Options{Quality: 90})
}

// convertFromWebp decodes a webp image and saves it in another format.
func convertFromWebp(input, output string) error {
	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open webp file: %w", err)
	}
	defer f.Close()

	img, err := webp.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode webp: %w", err)
	}

	outExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(output)), ".")

	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	switch outExt {
	case "png":
		return png.Encode(out, img)
	case "jpg", "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
	case "gif":
		return gif.Encode(out, img, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", outExt)
	}
}

// Ensure standard image decoders are registered via their init() functions.
var (
	_ = png.Decode
	_ = jpeg.Decode
	_ = gif.Decode
)
