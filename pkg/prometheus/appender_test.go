package prometheus

import (
	"github.com/golang/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dto "github.com/prometheus/client_model/go"
)

var _ = Describe("prometheus package", func() {
	Context("Appender functionality", func() {
		It("appends no metrics correctly", func() {
			metricFamilies := map[string]*dto.MetricFamily{}
			metrics := CreateMetricsCollection()

			metrics.Append(&metricFamilies)

			Expect(metrics.MetricTypes).Should(BeEmpty())
			Expect(metrics.MetricsByType).Should(BeEmpty())
		})

		It("does not contain duplicate typeswcorrectly", func() {
			dupFam1 := map[string]*dto.MetricFamily{
				"a_counter": &dto.MetricFamily{
					Name: proto.String("a_counter"),
					Type: dto.MetricType_COUNTER.Enum(),
					Metric: []*dto.Metric{
						{
							Label: []*dto.LabelPair{},
							Counter: &dto.Counter{
								Value: proto.Float64(0),
							},
						},
					},
				},
			}
			dupFam2 := map[string]*dto.MetricFamily{
				"a_counter": &dto.MetricFamily{
					Name: proto.String("a_counter"),
					Type: dto.MetricType_COUNTER.Enum(),
					Metric: []*dto.Metric{
						{
							Label: []*dto.LabelPair{},
							Counter: &dto.Counter{
								Value: proto.Float64(0),
							},
						},
					},
				},
			}
			metrics := CreateMetricsCollection()
			metrics.Append(&dupFam1)
			metrics.Append(&dupFam2)

			Expect(metrics.MetricTypes).Should(HaveLen(1))
			Expect(metrics.MetricsByType).Should(HaveLen(1))
			Expect(metrics.MetricsByType["a_counter"]).Should(HaveLen(2))
		})

		It("appends multiple families of metrics correctly", func() {
			metricFamilies := map[string]*dto.MetricFamily{
				"a_counter": &dto.MetricFamily{
					Name: proto.String("a_counter"),
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
				"a_gauge": &dto.MetricFamily{
					Name: proto.String("a_gauge"),
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

			metrics := CreateMetricsCollection()
			metrics.Append(&metricFamilies)

			Expect(metrics.MetricTypes).Should(HaveKey("a_counter"))
			Expect(metrics.MetricsByType["a_counter"]).Should(SatisfyAll(
				ContainElement(
					`a_counter{instance_id="some-instance",tag-a="some-value"} 2`,
				),
			))

			Expect(metrics.MetricTypes).Should(HaveKey("a_gauge"))
			Expect(metrics.MetricsByType["a_gauge"]).Should(SatisfyAll(
				ContainElement(
					`a_gauge{instance_id="some-instance",tag-a="some-value"} 1`,
				),
				ContainElement(
					`a_gauge{instance_id="some-instance",tag-b="some-value"} 99`,
				),
			))
		})
	})
})
