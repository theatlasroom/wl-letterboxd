package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

type WatchlistPoster struct {
	ID              int    `csv:"ID"`
	Name            string `csv:"Name"`
	Slug            string `csv:"Slug"`
	Height          int    `csv:"Height"`
	Width           int    `csv:"Width"`
	CacheBustingKey string `csv:"CacheBustingKey"`
}

type LetterboxdMovie struct {
	ID          int    `csv:"ID"`
	Name        string `csv:"Name"`
	ReleaseYear int    `csv:"Year"`
	// CreatedAt   time.Time `csv:"Date"`
	URI         string `csv:"Letterboxd URI"`
	MetadataURL string `csv:"Metadata URL`
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
func onPageHTML(e *colly.HTMLElement)    {}
func onRequest(r *colly.Request)         {}

func onResponse(r *colly.Response) {
	fmt.Println("Top Response: ", r.Request.URL)
	// fmt.Println("Response: ", string(r.Body))
}

func onError() {}

func onScraped(ms map[int]*LetterboxdMovie) {
	// called at the end of the whole scraping process
	for k, v := range ms {
		fmt.Println(k, ": ", v)
	}
}

func (p *WatchlistPoster) metadataURL() string {
	url := "%s/ajax/poster%sstd/%dx%d?k=%s"
	return fmt.Sprintf(url, lbxdBaseURL, p.Slug, p.Width, p.Height, p.CacheBustingKey)
}

func AttrInt(str string) int {
	val, strErr := strconv.Atoi(str)
	if strErr != nil {
		return -1
	}
	return val
}

func NewWatchlistPoster(el *colly.HTMLElement) *WatchlistPoster {
	return &WatchlistPoster{
		ID:              AttrInt(el.Attr(SELECTOR_FILM_ID)),
		Name:            el.DOM.Find("img").AttrOr("alt", ""),
		Slug:            el.Attr(SELECTOR_FILM_SLUG),
		Height:          AttrInt(el.Attr(SELECTOR_POSTER_HEIGHT)),
		Width:           AttrInt(el.Attr(SELECTOR_POSTER_WIDTH)),
		CacheBustingKey: el.Attr(SELECTOR_POSTER_CACHE_BUSTING_KEY),
	}
}

func NewLetterboxdMovie(p *WatchlistPoster) LetterboxdMovie {
	return LetterboxdMovie{
		ID:          p.ID,
		Name:        p.Name,
		MetadataURL: p.metadataURL(),
	}
}

func main() {
	// watchlistPosters := []WatchlistPoster{}
	movies := map[int]*LetterboxdMovie{}

	collector := colly.NewCollector(
		colly.CacheDir("./tmp"),
	)

	// TODO: parallelize collectors
	posterMetadataCollector := collector.Clone()
	// posterMetadataCollector.OnResponse(func(r *colly.Response) {
	// 	fmt.Println("metadata response: ", r.Request.URL)

	// })

	posterMetadataCollector.OnHTML(".film-poster", func(el *colly.HTMLElement) {
		fmt.Println("metadata html release year: ", AttrInt(el.Attr(SELECTOR_FILM_RELEASE_YEAR)))

		ID := AttrInt(el.Attr(SELECTOR_FILM_ID))
		ry := AttrInt(el.Attr(SELECTOR_FILM_RELEASE_YEAR))

		m := movies[ID]
		m.ReleaseYear = ry
	})

	collector.OnHTML(NODE_POSTER_CONTAINER, func(e *colly.HTMLElement) {
		e.ForEach(NODE_POSTER, func(_ int, el *colly.HTMLElement) {
			p := NewWatchlistPoster(el)

			m := NewLetterboxdMovie(p)
			movies[p.ID] = &m

			err := posterMetadataCollector.Visit(p.metadataURL())
			if err != nil {
				log.Fatal(err)
			}
		})
	})

	counter := 0
	collector.OnHTML(NODE_PAGINATION_NEXT_PAGE, func(e *colly.HTMLElement) {
		counter++
		nextPageURL := e.Attr("href")

		fmt.Println("Next page", nextPageURL, counter)
		if counter < 3 {
			e.Request.Visit(nextPageURL)
		}
	})

	collector.OnResponse(onResponse)

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting ", r.URL)
	})

	collector.OnScraped(func(r *colly.Response) {
		onScraped(movies)
	})

	collector.Visit(watchlistBaseURL)

	// TODO: allow revisiting URLs
	// TODO: persist data as we crawl
	// TODO: keep a key/value pair of pages visited
}
