package scraper

// MetricType represents the type of prometheus metric for a timeseries
type MetricType string

const (
	// PrometheusGauge a gauge metric type name
	PrometheusGauge MetricType = "gauge"
	// PrometheusCounter a counter metric type name
	PrometheusCounter MetricType = "counter"
	// PrometheusHistogram a histogram metric type name
	PrometheusHistogram MetricType = "histogram"
	// PrometheusSummary a summary metric type name
	PrometheusSummary MetricType = "summary"
)

// PrometheusGaugeMetric represents the data contained in an individual gauge metric timeseries
type PrometheusGaugeMetric struct {
	Key    string
	Labels map[string]string
	Value  float64
}

// PrometheusCounterMetric represents the data contained in an individual counter metric timeseries
type PrometheusCounterMetric struct {
	Key    string
	Labels map[string]string
	Value  int
}

// PrometheusHistogramMetric represents the data contained in an individual histogram metric timeseries
// TODO -- fix to represent an actual histogram metric's data
type PrometheusHistogramMetric struct {
	Key    string
	Labels map[string]string
	Value  int
}

// PrometheusSummaryMetric represents the data contained in an individual summary metric timeseries
// TODO -- fix to represent an actual summary metric's data
type PrometheusSummaryMetric struct {
	Key    string
	Labels map[string]string
	Value  int
}

// ScrapeData contains all the serialized data from a prometheus compatible metrics endpoint
type ScrapeData struct {
	helps map[string]string
	types map[string]MetricType

	Gauges   []PrometheusGaugeMetric
	Counters []PrometheusCounterMetric
}
