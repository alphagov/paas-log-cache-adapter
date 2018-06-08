package prometheus

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alphagov/paas-log-cache-adapter/pkg/metric"
)

// Converter is responsible for generating Prometheus friendly metric entity in
// a form of slice of bytes.
func Converter(metrics []*metric.Metric) ([]byte, error) {
	b := []string{}

	uniqueMetrics := map[string]float64{}
	for _, m := range metrics {
		uniqueMetrics[fmt.Sprintf("%s%s", m.Key, composeMeta(m))] = m.Value
	}

	for key, value := range uniqueMetrics {
		b = append(b, fmt.Sprintf("%s %g", key, value))
	}

	return []byte(strings.Join(b, "\n")), nil
}

func composeMeta(m *metric.Metric) string {
	if m.Meta == nil {
		return ""
	}

	all := []string{}

	var keys []string
	for k := range m.Meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		all = append(all, fmt.Sprintf(`%s="%s"`, k, m.Meta[k]))
	}

	return fmt.Sprintf("{%s}", strings.Join(all, ","))
}
