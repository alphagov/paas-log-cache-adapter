package prometheus

type MetricsCollection struct {
	MetricTypes   map[string]string
	MetricsByType map[string][]string
}

func CreateMetricsCollection() MetricsCollection {
	return MetricsCollection{
		MetricTypes:   make(map[string]string),
		MetricsByType: make(map[string][]string),
	}
}
