package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type WatchlistPoster struct {
	ID              int    `csv:"ID"`
	Name            string `csv:"Name"`
	Slug            string `csv:"Slug"`
	Height          int    `csv:"Height"`
	Width           int    `csv:"Width"`
	CacheBustingKey string `csv:"CacheBustingKey"`
	MetadataURL     string `csv:"MetadataURL` // should only ever need to construct this once for a new movie
}

type LetterboxdMovie struct {
	ID          int       `csv:"ID"`
	Name        string    `csv:"Name"`
	ReleaseYear int16     `csv:"Year"`
	CreatedAt   time.Time `csv:"Date"`
	URI         string    `csv:"Letterboxd URI"`
}

const lbxdBaseURL = "https://letterboxd.com"
const watchlistBaseURL = "https://letterboxd.com/junior1z1337/watchlist"

const (
	NODE_POSTER_CONTAINER = ".poster-container"
	NODE_POSTER           = ".poster-container .film-poster"
	NODE_POSTER_IMAGE     = "img"
	// NODE_PAGINATION_LAST_PAGE = ".pagination .paginate-page"
	NODE_PAGINATION_CONTAINER    = ".paginate-nextprev .next"
	NODE_PAGINATION_NEXT_PAGE    = ".pagination .paginate-nextprev .next"
	NODE_PAGINATION_CURRENT_PAGE = ".pagination .paginate-current"
)

const (
	SELECTOR_FILM_NAME                = "data-film-name"
	SELECTOR_FILM_ID                  = "data-film-id"
	SELECTOR_FILM_SLUG                = "data-film-slug"
	SELECTOR_FILM_RELEASE_YEAR        = "data-film-release-year"
	SELECTOR_POSTER_HEIGHT            = "data-image-height"
	SELECTOR_POSTER_WIDTH             = "data-image-width"
	SELECTOR_POSTER_CACHE_BUSTING_KEY = "data-cache-busting-key"
)

func extractPoster(e *colly.HTMLElement) {}
func onScraped()                         {} // called at the end of the whole scraping process
func onPageHTML(e *colly.HTMLElement)    {}
func onRequest(r *colly.Request)         {}

func onResponse(r *colly.Response) {
	fmt.Println("Response: ", r.Request.URL)
	fmt.Println("Response: ", string(r.Body))
}

func onError() {}

func metadataURL(p *WatchlistPoster) string {
	url := "/ajax/poster%sstd/%dx%d?k=%s"
	return fmt.Sprintf(url, p.Slug, p.Width, p.Height, p.CacheBustingKey)
}

func AttrInt(str string) int {
	val, strErr := strconv.Atoi(str)
	if strErr != nil {
		return -1
	}
	return val
}

func NewWatchlistPoster(el *colly.HTMLElement) *WatchlistPoster {
	p := &WatchlistPoster{
		ID:              AttrInt(el.Attr(SELECTOR_FILM_ID)),
		Name:            el.DOM.Find("img").AttrOr("alt", ""),
		Slug:            el.Attr(SELECTOR_FILM_SLUG),
		Height:          AttrInt(el.Attr(SELECTOR_POSTER_HEIGHT)),
		Width:           AttrInt(el.Attr(SELECTOR_POSTER_WIDTH)),
		CacheBustingKey: el.Attr(SELECTOR_POSTER_CACHE_BUSTING_KEY),
	}
	p.MetadataURL = metadataURL(p)

	return p
}

func main() {
	watchlistPosters := []WatchlistPoster{}

	collector := colly.NewCollector(
		colly.CacheDir("./tmp"),
	)

	// collector.Async = true
	// paginationCollector := collector.Clone()

	collector.OnHTML(NODE_POSTER_CONTAINER, func(e *colly.HTMLElement) {
		e.ForEach(NODE_POSTER, func(_ int, el *colly.HTMLElement) {

			// ReleaseYear, ryerr := strconv.Atoi(el.Attr(SELECTOR_FILM_RELEASE_YEAR))
			// if ryerr != nil {
			// 	ReleaseYear = 0
			// }

			p := NewWatchlistPoster(el)

			fmt.Println("Found a poster", p)
			// fmt.Sprintf("Found a poster %+v, %s", p, metadataURL(p))
			// fire up a filmCollector
			// for each slug, we should visit the page and extract the metadata

			e.Request.Visit(p.MetadataURL)
			watchlistPosters = append(watchlistPosters, *p)
		})
	})

	// collector.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("RESPONSE", r.Request.URL)

	// 	if strings.Index(r.Headers.Get("x-letterboxd-type"), "Film") > -1 {
	// 		fmt.Println("r", r.Body)
	// 	}
	// })

	// collector.OnHTML(NODE_POSTER_CONTAINER, func(e *colly.HTMLElement) {
	// 	e.ForEach(NODE_POSTER, func(_ int, el *colly.HTMLElement) {

	// 		ID, iderr := strconv.Atoi(el.Attr(SELECTOR_FILM_ID))

	// 		if iderr != nil {
	// 			ID = -1
	// 		}

	// 		ReleaseYear, ryerr := strconv.Atoi(el.Attr(SELECTOR_FILM_RELEASE_YEAR))
	// 		if ryerr != nil {
	// 			ReleaseYear = 0
	// 		}

	// 		p := LetterboxdMovie{
	// 			ID:          ID,
	// 			ReleaseYear: int16(ReleaseYear),
	// 			Name:        el.Attr(SELECTOR_FILM_NAME),
	// 			URI:         el.Attr(SELECTOR_FILM_SLUG),
	// 		}

	// 		fmt.Println("Found a poster", p)

	// 		posters = append(posters, p)
	// 	})
	// })

	collector.OnHTML(NODE_PAGINATION_NEXT_PAGE, func(e *colly.HTMLElement) {
		nextPageURL := e.Attr("href")

		fmt.Println("Next page", nextPageURL)
		// e.Request.Visit(nextPage)
	})

	collector.OnResponse(onResponse)

	// paginationCollector.OnHTML(NODE_PAGINATION_PAGE, func(e *colly.HTMLElement) {
	// 	nextURL := e.Attr("href")
	// })

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting ", r.URL)
	})

	collector.Visit(watchlistBaseURL)
}
