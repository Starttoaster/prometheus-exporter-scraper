package scraper

// GetHelps returns a map[string]string of key metric names and value metric help messages
func (d *ScrapeData) GetHelps() map[string]string {
	return d.helps
}

// GetHelp returns the help message for a given metric name. Returns a blank string if the metric had no help message
func (d *ScrapeData) GetHelp(metric string) string {
	return d.helps[metric]
}

// GetTypes returns a map[string]MetricType of key metric names and value metric types
func (d *ScrapeData) GetTypes() map[string]MetricType {
	return d.types
}

// GetType returns the type for a given metric name
func (d *ScrapeData) GetType(metric string) MetricType {
	return d.types[metric]
}
