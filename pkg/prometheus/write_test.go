package prometheus

import (
	"bytes"
	"io"

	"github.com/golang/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dto "github.com/prometheus/client_model/go"
)

var _ = Describe("prometheus package", func() {
	Context("Write functionality", func() {
		It("write no metrics correctly", func() {
			metricFamilies := map[string]*dto.MetricFamily{}
			expected := ""
			buffer := new(bytes.Buffer)
			WriteMetrics(metricFamilies, io.Writer(buffer))
			Expect(buffer.String()).To(Equal(expected))
		})

		It("write multiple families of metrics correctly", func() {
			metricFamilies := map[string]*dto.MetricFamily{
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
						{
							Label: []*dto.LabelPair{
								{
									Name:  proto.String("instance_id"),
									Value: proto.String("some-instance"),
								},
								{
									Name:  proto.String("tag-b"),
									Value: proto.String("some-value"),
								},
							},
							Gauge: &dto.Gauge{
								Value: proto.Float64(99),
							},
						},
					},
				},
			}

			buffer := new(bytes.Buffer)
			WriteMetrics(metricFamilies, io.Writer(buffer))
			Expect(buffer.String()).To(ContainSubstring(`# TYPE a-counter counter
a-counter{instance_id="some-instance",tag-a="some-value"} 2`))
			Expect(buffer.String()).To(ContainSubstring(`# TYPE a-gauge gauge
a-gauge{instance_id="some-instance",tag-a="some-value"} 1
a-gauge{instance_id="some-instance",tag-b="some-value"} 99`))
		})
	})
})
