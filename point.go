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
	CatNo    string
	URI      string
}

// CreatePoint makes a new Point struct containing record price info
func CreatePoint(discogs discogsAPI.Discogs, releaseID int) Point {
	release, _ := discogs.Release(releaseID)
	stats, _ := discogs.ReleaseStats(releaseID)

	var catNums []string
	for _, l := range release.Labels {
		catNums = append(catNums, l.Catno)
	}

	return Point{
		Artist:   release.ArtistsSort,
		Title:    release.Title,
		Id:       fmt.Sprintf("%d", releaseID),
		Currency: stats.LowestPrice.Currency,
		Price:    math.Round(stats.LowestPrice.Value*100) / 100,
		ForSale:  stats.ForSale,
		CatNo:    catNums[0], // This is optimistic since we will just grab the first match here
		URI:      release.URI,
	}
}
