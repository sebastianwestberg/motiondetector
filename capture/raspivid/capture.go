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

type Settings struct {
	Width  int
	Height int
	FPS    int
}

func (r *raspiVid) Capture() io.ReadCloser {
	param := strings.Fields(fmt.Sprintf("-o - -w %d -h %d -n -vf -t 0 -cd MJPEG -fps %d", 1024, 800, 3))
	cmd := exec.Command("raspivid", param...)

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

func NewCapture(settings Settings) domain.Capture {
	return &raspiVid{
		settings: settings,
	}
}
