package raspivid

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/sebastianwestberg/motiondetector/domain"
)

type raspiVid struct {
	settings Settings
}

type codec string

const (
	CodecMJPEG codec = "MJPEG"
	CodecH264  codec = "H264"
)

type Settings struct {
	Width          int
	Height         int
	FPS            int
	Codec          codec
	VerticalFlip   bool
	HorizontalFlip bool
	DetectOnFrame  int
	Debug          bool
}

func (r *raspiVid) Capture() io.ReadCloser {
	params := r.prepareParameters()
	cmd := exec.Command("raspivid", params...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	return stdout
}

func (r *raspiVid) prepareParameters() []string {
	params := strings.Fields(fmt.Sprintf("-o - -w %d -h %d -n -t 0 -cd %s -fps %d", r.settings.Width, r.settings.Height, r.settings.Codec, r.settings.FPS))
	if r.settings.HorizontalFlip == true {
		params = append(params, "-hf")
	}
	if r.settings.VerticalFlip == true {
		params = append(params, "-vf")
	}

	return params
}

func NewCapture(settings Settings) domain.Capture {
	return &raspiVid{
		settings: settings,
	}
}
