package prometheus

import (
	"io"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func WriteMetrics(metricFams map[string]*dto.MetricFamily, writer io.Writer) error {
	encoder := expfmt.NewEncoder(writer, expfmt.FmtText)

	for _, fam := range metricFams {
		err := encoder.Encode(fam)
		if err != nil {
			return err
		}
	}

	return nil
}
