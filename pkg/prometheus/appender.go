package prometheus

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func (metrics *MetricsCollection) Append(
	metricFams *map[string]*dto.MetricFamily,
) error {
	unfilteredOutput := new(bytes.Buffer)
	encoder := expfmt.NewEncoder(unfilteredOutput, expfmt.FmtText)

	for _, fam := range *metricFams {
		err := encoder.Encode(fam)
		if err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(unfilteredOutput)
	scanner.Split(bufio.ScanLines)
	metricTypeRegex, _ := regexp.Compile("^# TYPE ([^ ]+) .*$")
	metricRegex, _ := regexp.Compile("^([^{ ]*).*$")
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# TYPE") {
			metricTypeName := metricTypeRegex.FindStringSubmatch(line)[1]
			metrics.MetricTypes[metricTypeName] = line
		} else {
			metricName := metricRegex.FindStringSubmatch(line)[1]
			metrics.MetricsByType[metricName] = append(
				metrics.MetricsByType[metricName],
				line,
			)
		}
	}

	return nil
}
