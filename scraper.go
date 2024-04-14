package scraper

import (
	"fmt"
	"net/url"
	"os"
)

// FileScraper custom type for scraping a file that contains the data from a prometheus metrics endpoint
type FileScraper struct {
	file *os.File
}

// NewFileScraper returns a new file scraper instance from a given filepath
func NewFileScraper(path string) (*FileScraper, error) {
	s := new(FileScraper)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("making new file scraper: %w", err)
	}

	s.file = file

	return s, nil
}

// CloseFileScraper closes the file scraper gracefully
func (s *FileScraper) CloseFileScraper() error {
	return s.file.Close()
}

// WebScraper custom type for scraping an http endpoint that contains the data from a prometheus metrics endpoint
type WebScraper struct {
	url *url.URL
}

// NewWebScraper returns a new web scraper instance from a given url
func NewWebScraper(urlStr string) (*WebScraper, error) {
	s := new(WebScraper)

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("making new web scraper: %w", err)
	}
	s.url = u

	return s, nil
}
