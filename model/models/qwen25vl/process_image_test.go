package qwen25vl

import (
	"image"
	_ "image/jpeg" // Register JPEG decoder
	"testing"
)

func TestSmartResize(t *testing.T) {
	type smartResizeCase struct {
		TestImage image.Image
		Expected  image.Point
	}

	// Create an image processor with default values
	processor := ImageProcessor{
		imageSize:   224, // Example value
		numChannels: 3,
		factor:      28,
		minPixels:   56 * 56,
		maxPixels:   14 * 14 * 4 * 1280,
	}

	cases := []smartResizeCase{
		{
			TestImage: image.NewRGBA(image.Rect(0, 0, 1024, 1024)),
			Expected:  image.Point{980, 980},
		},
		{
			TestImage: image.NewRGBA(image.Rect(0, 0, 1024, 768)),
			Expected:  image.Point{1036, 756},
		},
		{
			TestImage: image.NewRGBA(image.Rect(0, 0, 2000, 2000)),
			Expected:  image.Point{980, 980},
		},
	}

	for _, c := range cases {
		b := c.TestImage.Bounds().Max
		actual := processor.smartResize(b)
		if actual != c.Expected {
			t.Errorf("expected: %v, actual: %v", c.Expected, actual)
		}
	}
}

func TestProcessImage(t *testing.T) {
	type processImageCase struct {
		TestImage   image.Image
		ExpectedLen int
	}

	// Create an image processor with default values
	processor := ImageProcessor{
		imageSize:   224, // Example value
		numChannels: 3,
		factor:      28,
		minPixels:   56 * 56,
		maxPixels:   14 * 14 * 4 * 1280,
	}

	cases := []processImageCase{
		{
			TestImage:   image.NewRGBA(image.Rect(0, 0, 256, 256)),
			ExpectedLen: 252 * 252 * 3,
		},
		{
			TestImage:   image.NewRGBA(image.Rect(0, 0, 2000, 2000)),
			ExpectedLen: 980 * 980 * 3,
		},
	}

	for _, c := range cases {
		imgData, err := processor.ProcessImage(c.TestImage)
		if err != nil {
			t.Fatalf("error processing: %q", err)
		}

		switch len(imgData) {
		case 0:
			t.Errorf("no image data returned")
		case c.ExpectedLen:
			// ok
		default:
			t.Errorf("unexpected image data length: %d, expected: %d", len(imgData), c.ExpectedLen)
		}
	}
}
