package domain

import "io"

// Capture is an interface for capturing video.
type Capture interface {
	Capture() io.ReadCloser
}
