package prometheus

import (
	"fmt"
	"io"
)

func (metrics *MetricsCollection) Write(writer io.Writer) {
	for metricTypeName, metricLines := range metrics.MetricsByType {
		if metricTypeLine, ok := metrics.MetricTypes[metricTypeName]; ok {
			fmt.Fprintln(writer, metricTypeLine)
		}
		for _, line := range metricLines {
			fmt.Fprintln(writer, line)
		}
	}
}
