package prometheus

import (
	"github.com/alphagov/paas-log-cache-adapter/pkg/metric"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("prometheus package", func() {
	Context("Convert functionality", func() {
		It("should be able to convert metric to prometheus format", func() {
			m := []*metric.Metric{
				&metric.Metric{Key: "test.success", Value: 1, Unit: "bool"},
			}

			res, err := Converter(m)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(ContainSubstring("test.success 1"))
		})

		It("should be able to convert metric to prometheus format with meta", func() {
			m := []*metric.Metric{
				&metric.Metric{
					Key:   "response",
					Value: 23,
					Unit:  "count",
					Meta: map[string]string{
						"status": "200",
						"type":   "json",
					},
				},
			}

			res, err := Converter(m)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(res)).To(ContainSubstring(`response{status="200",type="json"} 23`))
		})

		It("should be able to convert metric to prometheus format with float", func() {
			m := []*metric.Metric{
				&metric.Metric{
					Key:   "cpu_usage",
					Value: .45,
					Unit:  "percent",
				},
				&metric.Metric{
					Key:   "disk_usage",
					Value: .1,
					Unit:  "percent",
				},
			}

			res, err := Converter(m)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(res)).To(ContainSubstring(`cpu_usage 0.45`))
			Expect(string(res)).To(ContainSubstring(`disk_usage 0.1`))
		})
	})
})
