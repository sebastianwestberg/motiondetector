deploy:
	scp ./bin/motiondector/ pi@192.168.1.20:/home/pi/
build:
	GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/motiondector ./cmd/motiondetector/main.go
build/deploy: build deploy