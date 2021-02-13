.DEFAULT_GOAL := build

clean:
	rm -rf bin
	go mod tidy

build: clean build_arm build_x64

build_arm:
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/discogs-influxdb_arm .

build_x64:
	go build -o bin/discogs-influxdb .

install:
	mkdir -p /opt/discogs-influxdb
	cp bin/discogs-influxdb /opt/discogs-influxdb
	cp discogs-influxdb.service /etc/system/systemd/
	cp discogs-influxdb.timer /etc/system/systemd/