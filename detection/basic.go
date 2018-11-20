package detection

import (
	"fmt"
	"image"
	"log"
	"math"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/histogram"
	"github.com/anthonynsimon/bild/transform"
	"github.com/sebastianwestberg/motiondetector/domain"
)

type motionDetectionBasic struct {
}

func (m *motionDetectionBasic) Detect(ref image.Image, curr image.Image) bool {
	diff, err := fastCompare(
		histogram.NewRGBAHistogram(applyFilters(ref)).Image(),
		histogram.NewRGBAHistogram(applyFilters(curr)).Image(),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("diff:", diff)
	if diff >= 18000 {
		return true
	}

	return false
}

func applyFilters(img image.Image) *image.Gray {
	resized := transform.Resize(img, 250, 200, transform.Linear)
	blurred := blur.Gaussian(resized, 1)
	greyscaled := effect.Grayscale(blurred)

	return greyscaled
}

func fastCompare(img1, img2 *image.RGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func NewBasicMotionDetection() domain.MotionDetection {
	return &motionDetectionBasic{}
}
