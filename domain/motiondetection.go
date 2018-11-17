package domain

import "image"

type MotionDetection interface {
	Detect(prev image.Image, curr image.Image) bool
}
