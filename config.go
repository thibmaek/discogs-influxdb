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

func getInfluxDBConfig(c *toml.Tree) (h string, t string, b Database) {
	influxCfg := c.Get("influxdb").(*toml.Tree).ToMap()

	var influxToken string
	var influxBucket Database

	host := fmt.Sprintf("%s:%d", influxCfg["host"].(string), influxCfg["port"].(int64))

	token, tokenExists := influxCfg["token"]
	if tokenExists {
		// InfluxDB 2 requires org + bucket
		influxToken = token.(string)
	} else {
		// InfluxDB 1.8 token is represented by user + pass joined with ':'
		influxToken = fmt.Sprintf("%s:%s", influxCfg["user"].(string), influxCfg["password"].(string))
	}

	bucket, bucketExists := influxCfg["bucket"]
	org, orgExists := influxCfg["org"]
	if bucketExists && orgExists {
		// InfluxDB 2 uses org + bucket combination
		influxBucket = Database{
			Org:    org.(string),
			Bucket: bucket.(string),
		}
	} else {
		// InfluxDB 1.8 uses database + retention policy as the bucket and empty org name
		rp, rpExists := influxCfg["retention_policy"]
		if !rpExists {
			rp = ""
		}
		influxBucket = Database{
			Org:    "",
			Bucket: fmt.Sprintf("%s/%s", influxCfg["database"].(string), rp.(string)),
		}
	}

	return host, influxToken, influxBucket
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

	influxHost, influxToken, influxBucket := getInfluxDBConfig(config)

	return &Config{
		DiscogsToken: config.Get("discogs.token").(string),

		InfluxHost:   influxHost,
		InfluxToken:  influxToken,
		InfluxBucket: influxBucket,

		MonitoredReleases: config.Get("releases.monitored").([]interface{}),
	}
}
