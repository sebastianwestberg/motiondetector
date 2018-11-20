deploy:
	scp ./bin/raspividmotion/ pi@192.168.1.20:/home/pi/
build:
	GOOS=linux GOARCH=arm GOARM=7 go build -o ./bin/raspividmotion ./cmd/raspividmotion/main.go
build/deploy: build deploy