package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/sebastianwestberg/motiondetector/capture/raspivid"
	"github.com/sebastianwestberg/motiondetector/detection"
)

func main() {
	errorChannel := make(chan error)

	settings := raspivid.Settings{
		Width:  1024,
		Height: 800,
		FPS:    3,
	}
	r := raspivid.NewCapture(settings).Capture()

	motion := detection.NewBasicMotionDetection()

	go func() {
		var lastFrame image.Image
		var counter uint
		frameWorker(splitOnFrame(r), func(img image.Image) {
			timeNow := time.Now()

			defer func() {
				counter++
			}()
			if counter%3 != 0 {
				return
			}

			if lastFrame != nil {
				motion := motion.Detect(lastFrame, img)
				if motion == true {
					log.Println("Motion detected!!")

					if err := imgio.Save(fmt.Sprintf("imgs/motion-%d.jpg", counter), img, imgio.JPEG); err != nil {
						panic(err)
					}

					var buf bytes.Buffer
					jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
					os.Stdout.Write(buf.Bytes())
				}
			}

			lastFrame = img
			log.Printf("counter %d took %v\n", counter, time.Since(timeNow))
		})
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errorChannel <- fmt.Errorf("Got signal: %s", <-c)
	}()

	// go metrics()

	if err := <-errorChannel; err != nil {
		log.Fatalf("Got error: %+v\n", err)
	}
}

func frameWorker(framesScanner *bufio.Scanner, imageHandler func(img image.Image)) {
	for framesScanner.Scan() {
		frame := framesScanner.Bytes()
		if frame[0] == byte(0xFF) && frame[1] == byte(0xD8) {
			frame = append(frame, []byte{0xFF, 0xD9}...)
			b := bytes.NewReader(frame)

			img, err := jpeg.Decode(b)
			if err != nil {
				log.Fatal(err)
			}

			go imageHandler(img)
		}
	}
}

func splitOnFrame(r io.ReadCloser) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	splitFn := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data)-1; i++ {
			// Scan for end of image marker
			if data[i] == byte(0xFF) && data[i+1] == byte(0xD9) {
				return i + 2, data[:i], nil
			}
		}
		return 0, nil, nil
	}

	buffer := make([]byte, 0, bufio.MaxScanTokenSize)
	scanner.Split(splitFn)
	scanner.Buffer(buffer, 1024*1024)

	return scanner
}

func metrics() {
	for {
		mem := runtime.MemStats{}
		runtime.ReadMemStats(&mem)
		log.Println(
			"mem.alloc", mem.Alloc/1000000,
			"mem.heapSys", mem.HeapSys,
		)
		time.Sleep(2 * time.Second)
	}
}
