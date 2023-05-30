package lbxd

import (
	"fmt"
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
	URI         string `csv:"Letterboxd URI"`
}

const lbxdBaseURL = "https://letterboxd.com"
const WatchlistBaseURL = "https://letterboxd.com/junior1z1337/watchlist"

const (
	NODE_POSTER_CONTAINER = ".poster-container"
	NODE_POSTER           = ".film-poster"
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

func (p *WatchlistPoster) MetadataURL() string {
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
		ID:   p.ID,
		Name: p.Name,
	}
}
