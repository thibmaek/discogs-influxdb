package main

import (
	"fmt"
	"math"

	discogsAPI "github.com/thibmaek/influxdb-discogs/discogs"
)

type Point struct {
	Artist   string
	Title    string
	Id       string
	Currency string
	Price    float64
	ForSale  int
}

// CreatePoint makes a new Point struct containing record price info
func CreatePoint(discogs discogsAPI.Discogs, releaseID int) Point {
	release, _ := discogs.Release(releaseID)
	stats, _ := discogs.ReleaseStats(releaseID)

	return Point{
		Title:    release.Title,
		Id:       fmt.Sprintf("%d", releaseID),
		Currency: stats.LowestPrice.Currency,
		Price:    math.Round(stats.LowestPrice.Value*100) / 100,
		ForSale:  stats.ForSale,
	}
}
