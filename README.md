# prometheus-exporter-scraper

[![Go Report Card](https://goreportcard.com/badge/github.com/starttoaster/prometheus-exporter-scraper)](https://goreportcard.com/report/github.com/starttoaster/prometheus-exporter-scraper) [![Go Reference](https://pkg.go.dev/badge/github.com/starttoaster/prometheus-exporter-scraper.svg)](https://pkg.go.dev/github.com/starttoaster/prometheus-exporter-scraper)

This package is for parsing Prometheus-compliant metrics libraries to Go data types. It relies on regexp to parse line content into structs.

Currently supports counter and gauge metrics. Histogram and summary support coming soon.

This is mostly useful if all of the following conditions are true:
- You have a metrics exporter that you don't care to view as a timeseries.
- The application you're writing will be deployed to a location that can route to that metrics endpoint.

An example use case for this package is for extracting data from exporters with labels that have high levels of cardinality. Since scraping the exporter could negatively affect Prometheus server performance, it might sometimes be ideal to get that data at a single-point-in-time from the exporter directly.

## Usage

```go
import scraper "github.com/starttoaster/prometheus-exporter-scraper"

// Create scraper -- replace with your metrics URL -- ignores errors
scrp, _ := scraper.NewWebScraper("http://localhost:8080/metrics")

// Scrape metrics -- ignores errors
data, _ := scrp.ScrapeWeb()

// Loop through gauge metrics -- you can do the same with counters by using data.Counters
for _, gauge := range data.Gauges {
    fmt.Println(gauge.Key) // Print the metric name (value of type string)
    fmt.Println(gauge.Value) // Print the metric value (value of type float64)
    fmt.Println(gauge.Labels) // Print out the metric labels (value of type map[string]string)
    fmt.Println(data.GetHelp(gauge.Key)) // Print out the help/info message for a particular metric name
    fmt.Println(data.Type(gauge.Key)) // Print out the type for a particular metric name
}
```

## TODO

- Tests
- Break up the scanLine function into parts
