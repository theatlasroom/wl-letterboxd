package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type LetterboxdMovie struct {
	ID          int       `csv:"ID"`
	Name        string    `csv:"Name"`
	ReleaseYear int16     `csv:"Year"`
	CreatedAt   time.Time `csv:"Date"`
	URI         string    `csv:"Letterboxd URI"`
}

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
	SELECTOR_FILM_NAME         = "data-film-name"
	SELECTOR_FILM_ID           = "data-film-id"
	SELECTOR_FILM_SLUG         = "data-film-slug"
	SELECTOR_FILM_RELEASE_YEAR = "data-film-release-year"
)

func extractPoster(e *colly.HTMLElement) {}
func onPageHTML(e *colly.HTMLElement)    {}
func onRequest(r *colly.Request)         {}

func main() {
	posters := []LetterboxdMovie{}

	collector := colly.NewCollector(
		colly.CacheDir("./tmp"),
	)
	// paginationCollector := collector.Clone()

	collector.OnHTML(NODE_POSTER_CONTAINER, func(e *colly.HTMLElement) {
		e.ForEach(NODE_POSTER, func(_ int, el *colly.HTMLElement) {

			ID, iderr := strconv.Atoi(el.Attr(SELECTOR_FILM_ID))

			if iderr != nil {
				ID = -1
			}

			ReleaseYear, ryerr := strconv.Atoi(el.Attr(SELECTOR_FILM_RELEASE_YEAR))
			if ryerr != nil {
				ReleaseYear = 0
			}

			p := LetterboxdMovie{
				ID:          ID,
				ReleaseYear: int16(ReleaseYear),
				Name:        el.DOM.Find("img").AttrOr("alt", ""),
				URI:         el.Attr(SELECTOR_FILM_SLUG),
			}

			fmt.Println("Found a poster", p)
			// fire up a filmCollector
			// for each slug, we should visit the page and extract the metadata

			posters = append(posters, p)
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

	// paginationCollector.OnHTML(NODE_PAGINATION_PAGE, func(e *colly.HTMLElement) {
	// 	nextURL := e.Attr("href")
	// })

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting ", r.URL)
	})

	collector.Visit(watchlistBaseURL)
}
