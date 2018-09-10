package prometheus

import "sync"

type MetricsCollection struct {
	MetricTypes   map[string]string
	MetricsByType map[string][]string
	mux           sync.Mutex
}

func CreateMetricsCollection() MetricsCollection {
	return MetricsCollection{
		MetricTypes:   make(map[string]string),
		MetricsByType: make(map[string][]string),
	}
}
