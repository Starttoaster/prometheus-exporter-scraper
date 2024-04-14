package scraper

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ScrapeFile scrapes the file owned by a FileScraper instance and returns some ScrapeData
func (s *FileScraper) ScrapeFile() (*ScrapeData, error) {
	var data *ScrapeData = new(ScrapeData)
	data.helps = make(map[string]string)
	data.types = make(map[string]MetricType)

	// Read lines one by one
	var err error
	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		line := scanner.Text()
		data, err = scanLine(data, line)
		if err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning scrape file: %w", err)
	}

	return data, nil
}

func (s *WebScraper) ScrapeWeb() (*ScrapeData, error) {
	var data *ScrapeData = new(ScrapeData)
	data.helps = make(map[string]string)
	data.types = make(map[string]MetricType)

	// Make a GET request
	resp, err := http.Get(s.url.String())
	if err != nil {
		return nil, fmt.Errorf("error making GET to given endpoint \"%s\": %w", s.url.String(), err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("scraper failed to close response body. If you're seeing this in your logs then something unexpected is happening when closing the response body from your prometheus metrics exporter scrape requests. See %v", err)
		}
	}()

	// Read response body lines one by one
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		data, err = scanLine(data, line)
		if err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning scrape response body: %w", err)
	}

	return data, nil
}

func scanLine(data *ScrapeData, line string) (*ScrapeData, error) {
	switch {
	case strings.HasPrefix(line, "# HELP"):
		parts := strings.SplitN(line, " ", 4)
		metric := parts[2]
		message := parts[3]
		data.helps[metric] = message
	case strings.HasPrefix(line, "# TYPE"):
		parts := strings.SplitN(line, " ", 4)
		metric := parts[2]
		metricType := parts[3]
		data.types[metric] = MetricType(metricType)
	default:
		// Define a regular expression pattern to match the components
		re := regexp.MustCompile(`(\w+)(?:\{([^}]+)\})?\s+([\d.]+)`)

		// Find submatches
		matches := re.FindStringSubmatch(line)

		// For a valid metric line, we need at least 3 matches
		if len(matches) < 3 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		// Extract line components
		metricName := matches[1]
		var keyValuePairs string
		var metricValue string

		switch len(matches) {
		case 3:
			metricValue = matches[2]
		case 4:
			keyValuePairs = matches[2]
			metricValue = matches[3]
		default:
			return nil, fmt.Errorf("find string submatches returned invalid line: %s", line)
		}

		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return data, nil
		}

		// Parse key-value pairs into a map
		labelMap := make(map[string]string)
		kvPairs := regexp.MustCompile(`(\w+)="([^"]+)"`).FindAllStringSubmatch(keyValuePairs, -1)
		for _, kv := range kvPairs {
			if len(kv) < 3 {
				return nil, fmt.Errorf("invalid labels in line (line=%s) (labels=%s)", line, keyValuePairs)
			}
			key := kv[1]
			value := kv[2]
			labelMap[key] = value
		}

		// Get metric type
		var metricType = data.GetType(metricName)
		// If this is a summary or histogram, there are a few suffixes that could be added at the end
		if metricType == "" {
			suffixes := []string{"_bucket", "_count", "_sum"}
			for _, suffix := range suffixes {
				metricTypeMatch := strings.TrimSuffix(metricName, suffix)
				metricType = data.GetType(metricTypeMatch)

				if metricType != "" {
					break
				}
			}
		}

		// Try to find the metric's type to parse it a certain way
		switch MetricType(metricType) {
		case PrometheusCounter:
			data.Counters = append(data.Counters, PrometheusCounterMetric{
				Key:    metricName,
				Labels: labelMap,
				Value:  int(value),
			})
		case PrometheusGauge:
			data.Gauges = append(data.Gauges, PrometheusGaugeMetric{
				Key:    metricName,
				Labels: labelMap,
				Value:  value,
			})
		case PrometheusHistogram:
		case PrometheusSummary:
		default:
			return nil, fmt.Errorf("invalid metric, no type found: %s", line)
		}
	}
	return data, nil
}
