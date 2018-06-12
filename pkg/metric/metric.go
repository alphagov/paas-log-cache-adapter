package metric

import (
	"encoding/json"
)

// Converter is defined interface, for additional/external functionality, that
// would convert Metric struct into different format. For instance, prometheus
// or statsd.
type Converter func([]*Metric) ([]byte, error)

// Metric is the standard entity that would be used for representation to the
// end user.
type Metric struct {
	Key   string            `json:"key"`
	Value float64           `json:"value"`
	Unit  string            `json:"unit"`
	Meta  map[string]string `json:"meta,omitempty"`
}

// JSONConverter is simple generator of JSON version of a metric. It implements
// the Converter interface.
func JSONConverter(m []*Metric) ([]byte, error) {
	return json.Marshal(m)
}
