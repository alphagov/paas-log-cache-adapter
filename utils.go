package main

import (
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/alphagov/paas-log-cache-adapter/pkg/metric"
)

func convertToMetrics(r []*loggregator_v2.Envelope) []*metric.Metric {
	data := []*metric.Metric{}

	for _, envelope := range r {
		if envelope.GetCounter() != nil {
			m := metricFromCounter(envelope, envelope.GetCounter())
			data = append(data, m)
		}
		if envelope.GetGauge() != nil {
			m := metricFromGauge(envelope, envelope.GetGauge())
			data = append(data, m...)
		}
	}

	return data
}

func metricFromCounter(e *loggregator_v2.Envelope, c *loggregator_v2.Counter) *metric.Metric {
	data := &metric.Metric{
		Key:   c.GetName(),
		Value: float64(c.GetTotal()),
		Unit:  "count",
		Meta:  e.Tags,
	}

	data.Meta["instance_id"] = e.GetSourceId()
	data.Meta["type"] = "counter"

	return data
}

func metricFromGauge(e *loggregator_v2.Envelope, c *loggregator_v2.Gauge) []*metric.Metric {
	data := []*metric.Metric{}

	for name, g := range c.GetMetrics() {
		m := &metric.Metric{
			Meta: e.Tags,
		}
		m.Meta["instance_id"] = e.GetSourceId()
		m.Meta["type"] = "gauge"
		m.Key = name
		m.Value = g.GetValue()
		m.Unit = g.GetUnit()

		data = append(data, m)
	}

	return data
}
