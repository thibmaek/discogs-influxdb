package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	discogsAPI "github.com/thibmaek/influxdb-discogs/discogs"
)

func influxClient(c *Config) api.WriteAPIBlocking {
	i := influxdb2.NewClient(c.InfluxHost, c.DiscogsToken)
	return i.WriteAPIBlocking(c.InfluxBucket.Org, c.InfluxBucket.Bucket)
}

func discogsClient(c *Config) discogsAPI.Discogs {
	client, err := discogsAPI.New(&discogsAPI.Options{
		Token:     c.DiscogsToken,
		Currency:  "EUR",
		UserAgent: "InfluxDB Discogs",
	})

	if err != nil {
		log.Fatalf("Could not create Discogs client %v", err)
	}

	return client
}

func createPoints(r []interface{}, discogs discogsAPI.Discogs) []Point {
	var releases []Point
	for _, v := range r {
		p := CreatePoint(discogs, int(v.(int64)))
		releases = append(releases, p)
	}

	return releases
}

func getFlags() (*string, *bool) {
	cwd, _ := os.Executable()
	cfgPath := flag.String("config", path.Join(cwd, "..", "config.toml"), "Path to the config.toml file")
	verbose := flag.Bool("verbose", false, "Show verbose output about records being written")
	flag.Parse()

	if *verbose {
		p, _ := filepath.Abs(*cfgPath)
		fmt.Printf("Using config file %s\n", p)
	}

	return cfgPath, verbose
}

func main() {
	f, v := getFlags()
	c := GetConfig(*f)

	discogs := discogsClient(c)
	influx := influxClient(c)

	if *v {
		fmt.Printf("Writing to InfluxDB host: %s on bucket: %+v\n", c.InfluxHost, c.InfluxBucket)
	}

	releases := createPoints(c.MonitoredReleases, discogs)

	for _, r := range releases {
		point := influxdb2.NewPointWithMeasurement("listing").
			AddTag("id", r.Id).
			AddTag("title", r.Title).
			AddTag("artist", r.Artist).
			AddTag("uri", r.URI).
			AddTag("catno", r.CatNo).
			AddTag("currency", r.Currency).
			AddField("price", r.Price).
			AddField("num_for_sale", r.ForSale).
			SetTime(time.Now())

		if *v {
			fmt.Printf("Writing measurement 'listing' for release %s\n", r.Id)
		}

		err := influx.WritePoint(context.Background(), point)
		if err != nil {
			log.Printf("Failed to write point for release '%s': %v\n", r.Id, err)
		}
	}
}
