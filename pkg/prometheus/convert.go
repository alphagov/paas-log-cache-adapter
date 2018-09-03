package prometheus

import (
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/gogo/protobuf/proto"
	dto "github.com/prometheus/client_model/go"
)

func Convert(r []*loggregator_v2.Envelope) map[string]*dto.MetricFamily {
	metricFamilies := map[string]*dto.MetricFamily{}

	for _, envelope := range r {
		if envelope.GetCounter() != nil {
			counter := envelope.GetCounter()
			name := Sanitize(counter.GetName())

			labels := makeLabels(
				name,
				envelope.GetTags(),
				envelope.GetSourceId(),
			)

			metric := &dto.Metric{
				Label: labels,
				Counter: &dto.Counter{
					Value: proto.Float64(float64(counter.GetTotal())),
				},
			}

			addMetricToFamily(
				name,
				metric,
				metricFamilies,
				dto.MetricType_COUNTER.Enum(),
			)
		}

		if envelope.GetGauge() != nil {
			gauge := envelope.GetGauge()
			for name, g := range gauge.GetMetrics() {
				name = Sanitize(name)

				labels := makeLabels(
					name,
					envelope.GetTags(),
					envelope.GetSourceId(),
				)

				metric := &dto.Metric{
					Label: labels,
					Gauge: &dto.Gauge{
						Value: proto.Float64(float64(g.GetValue())),
					},
				}

				addMetricToFamily(
					name,
					metric,
					metricFamilies,
					dto.MetricType_GAUGE.Enum(),
				)
			}
		}
	}

	return metricFamilies
}

func addMetricToFamily(
	name string,
	metric *dto.Metric,
	metricFamilies map[string]*dto.MetricFamily,
	metricType *dto.MetricType,
) {
	metricFamily, ok := metricFamilies[name]
	if !ok {
		metricFamily = &dto.MetricFamily{
			Name:   proto.String(name),
			Metric: make([]*dto.Metric, 0),
			Type:   metricType,
		}
		metricFamilies[name] = metricFamily
	}

	metricFamily.Metric = append(metricFamily.Metric, metric)
}

func makeLabels(
	name string,
	tags map[string]string,
	sourceID string,
) []*dto.LabelPair {

	labels := make([]*dto.LabelPair, 0)

	labels = append(labels, &dto.LabelPair{
		Name:  proto.String("instance_id"),
		Value: proto.String(sourceID),
	})

	for k, v := range tags {
		label := &dto.LabelPair{
			Name:  proto.String(Sanitize(k)),
			Value: proto.String(v),
		}
		labels = append(labels, label)
	}

	return labels
}
