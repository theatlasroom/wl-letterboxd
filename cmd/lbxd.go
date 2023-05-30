package main

import (
	"fmt"
	"log"

	lbxd "github.com/theatlasroom/wl-letterboxd"

	"github.com/gocolly/colly/v2"
)

func main() {
	movies := map[int]*lbxd.LetterboxdMovie{}

	collector := colly.NewCollector(
		colly.CacheDir("./tmp"),
	)

	posterMetadataCollector := collector.Clone()

	posterMetadataCollector.OnHTML(lbxd.NODE_POSTER, func(el *colly.HTMLElement) {
		ID := lbxd.AttrInt(el.Attr(lbxd.SELECTOR_FILM_ID))
		ry := lbxd.AttrInt(el.Attr(lbxd.SELECTOR_FILM_RELEASE_YEAR))

		m := movies[ID]
		m.ReleaseYear = ry

		fmt.Println(m)
	})

	collector.OnHTML(lbxd.NODE_POSTER_CONTAINER, func(e *colly.HTMLElement) {
		e.ForEach(lbxd.NODE_POSTER, func(_ int, el *colly.HTMLElement) {
			p := lbxd.NewWatchlistPoster(el)

			m := lbxd.NewLetterboxdMovie(p)
			movies[p.ID] = &m

			err := posterMetadataCollector.Visit(p.MetadataURL())
			if err != nil {
				log.Fatal(err)
			}
		})
	})

	counter := 0
	collector.OnHTML(lbxd.NODE_PAGINATION_NEXT_PAGE, func(e *colly.HTMLElement) {
		counter++
		nextPageURL := e.Attr("href")

		fmt.Println("Next page", nextPageURL, counter)
		if counter < 3 {
			e.Request.Visit(nextPageURL)
		}
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting ", r.URL)
	})

	collector.Visit(lbxd.WatchlistBaseURL)

	// TODO: allow revisiting URLs
	// TODO: persist data as we crawl
	// TODO: keep a key/value pair of pages visited
}
