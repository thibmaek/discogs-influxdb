package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

type Database struct {
	Org    string
	Bucket string
}

type Config struct {
	DiscogsToken      string
	InfluxToken       string
	InfluxHost        string
	InfluxBucket      Database
	MonitoredReleases []interface{}
}

func GetConfig(cfgPath string) *Config {
	f, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	config, err := toml.Load(string(f))
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	influxCfg := config.Get("influxdb").(*toml.Tree).ToMap()

	return &Config{
		DiscogsToken: config.Get("discogs.token").(string),

		InfluxHost:  fmt.Sprintf("%s:%d", influxCfg["host"].(string), influxCfg["port"].(int64)),
		InfluxToken: influxCfg["token"].(string),
		InfluxBucket: Database{
			Org:    influxCfg["org"].(string),
			Bucket: influxCfg["bucket"].(string),
		},

		MonitoredReleases: config.Get("releases.monitored").([]interface{}),
	}
}
