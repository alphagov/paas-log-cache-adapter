package prometheus

import (
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/golang/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dto "github.com/prometheus/client_model/go"
)

var _ = Describe("prometheus package", func() {
	Context("Convert functionality", func() {
		It("should be able to convert empty loggregator Metrics into empty Prometheus Metrics", func() {
			metrics := []*loggregator_v2.Envelope{}

			expected := map[string]*dto.MetricFamily{}
			actual := Convert(metrics)
			Expect(actual).To(Equal(expected))
		})
		It("should be able to convert loggregator Metrics into Prometheus Metrics", func() {
			metrics := []*loggregator_v2.Envelope{
				{
					SourceId: "some-instance",
					Tags: map[string]string{
						"tag-a": "some-value",
					},
					Message: &loggregator_v2.Envelope_Counter{
						Counter: &loggregator_v2.Counter{
							Name:  "a-counter",
							Total: 2,
						},
					},
				},
				{
					SourceId: "some-instance",
					Tags: map[string]string{
						"tag-a": "some-value",
					},
					Message: &loggregator_v2.Envelope_Gauge{
						Gauge: &loggregator_v2.Gauge{
							Metrics: map[string]*loggregator_v2.GaugeValue{
								"a-gauge": &loggregator_v2.GaugeValue{
									Unit:  "b",
									Value: 1,
								},
							},
						},
					},
				},
			}

			expected := map[string]*dto.MetricFamily{
				"a-counter": &dto.MetricFamily{
					Name: proto.String("a-counter"),
					Type: dto.MetricType_COUNTER.Enum(),
					Metric: []*dto.Metric{
						{
							Label: []*dto.LabelPair{
								{
									Name:  proto.String("instance_id"),
									Value: proto.String("some-instance"),
								},
								{
									Name:  proto.String("tag-a"),
									Value: proto.String("some-value"),
								},
							},
							Counter: &dto.Counter{
								Value: proto.Float64(2),
							},
						},
					},
				},
				"a-gauge": &dto.MetricFamily{
					Name: proto.String("a-gauge"),
					Type: dto.MetricType_GAUGE.Enum(),
					Metric: []*dto.Metric{
						{
							Label: []*dto.LabelPair{
								{
									Name:  proto.String("instance_id"),
									Value: proto.String("some-instance"),
								},
								{
									Name:  proto.String("tag-a"),
									Value: proto.String("some-value"),
								},
							},
							Gauge: &dto.Gauge{
								Value: proto.Float64(1),
							},
						},
					},
				},
			}

			actual := Convert(metrics)
			Expect(actual).To(Equal(expected))
		})

	})
})
