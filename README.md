# prometheus-exporter-scraper

This package is for parsing Prometheus-compliant metrics libraries to Go data types. It relies on regexp to parse line content into structs

Currently supports counter and gauge metrics. Histogram and summary support coming soon.

This is mostly useful if all of the following conditions are true:
- You have a metrics exporter that you don't care to view as a timeseries, and is not scraped by Prometheus.
- The application you're writing will be deployed to a location that can route to that metrics endpoint.

I am currently using this package to scrape single point-in-time metrics from an exporter with high cardinality in labels that could impact my Prometheus server's performance.

## Usage

This API client library currently supports API tokens for authentication.

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