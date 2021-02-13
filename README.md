# Discogs InfluxDB

<!-- Banner & Badges. Badges should have newlines -->
[![Go Report](https://goreportcard.com/badge/github.com/thibmaek/discogs-influxdb)](https://github.com/thibmaek/discogs-influxdb)

Export price data from Discogs listings to InfluxDB

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
  - [Config file](#config-file)
  - [CLI](#cli)
  - [Systemd](#systemd)
  - [Others](#others)
- [API](#api)
- [License](#license)

## Background

I wanted a way to visually track the prices of expensive records in Grafana but Telegraf and Discogs API had no way of working together. So I created this small tool written in Go to pull prices from Discogs API (using a fork of [`go-discogs`](https://github.com/irlndts/go-discogs)) and push them to InfluxDB using influxdbv2 client.

The fork of is locally embedded in the project because the changes for Discogs' Marketplace API are not available in upstream.

## Install

### From releases

Download one of the releases from the Releases page.

### Building with Go

You can very easily build the binaries yourself. A Makefile is included for easier building for all different platforms:

```shell
# This will build ARM and x64
$ make build

# This will simply build using the standard go toolchain
$ go build .
```

## Usage

### Config file

You will need to create a `config.toml` file with your credentials for InfluxDB and Discogs. By default, `discogs-influxb` will look in its own directory for a `config.toml` file, but an alternative path can be provided using `--config <path>`

A sample config file is included in `config.sample.toml`. Edit that one with your correct credentials and rename it to `config.toml`

### CLI

You can invoke the binary to just trigger a one time push to InfluxDB.

### Systemd

A sample systemd service and a timer are included to run `discogs-influxdb` at regular intervals.
Copy the files for systemd over to `/etc/systemd/system/` and then copy the binary and config file for `discogs-influxdb` over to `/opt/discogs-influxdb` (or change the path in the systemd service), afterwards reload and restart:

```shell
$ cp discogs-influxdb.service /etc/system/systemd/
$ cp discogs-influxdb.timer /etc/system/systemd/
$ systemctl daemon-reload
$ systemctl enable discogs-influxdb.timer
$ systemctl start discogs-influxdb.service
```

### Others

pm2, cronjobs, docker, bash script...

## API

- `--config <path>`, path to the config.toml file to use
- `--verbose`, show verbose output

## License

MIT

For more info, see [license file](./LICENSE)
