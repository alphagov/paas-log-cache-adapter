package prometheus

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("prometheus package", func() {
	Context("Write functionality", func() {
		It("write no metrics correctly", func() {
			metrics := CreateMetricsCollection()
			output := new(bytes.Buffer)

			metrics.Write(output)
			Expect(output.String()).Should(Equal(""))
		})

		It("write multiple families of metrics correctly", func() {
			metrics := CreateMetricsCollection()

			metrics.MetricTypes["a_counter"] = "# TYPE a_counter counter"
			metrics.MetricsByType["a_counter"] = append(
				metrics.MetricsByType["a_counter"],
				`a_counter{tag="aVal"} 5`,
				`a_counter{tag="bVal"} 6`,
			)

			metrics.MetricTypes["a_gauge"] = "# TYPE a_gauge gauge"
			metrics.MetricsByType["a_gauge"] = append(
				metrics.MetricsByType["a_gauge"],
				`a_gauge{tag="aVal"} 1`,
				`a_gauge{tag="bVal"} 2`,
			)

			output := new(bytes.Buffer)

			metrics.Write(output)
			Expect(output.String()).Should(ContainSubstring(`# TYPE a_counter counter
a_counter{tag="aVal"} 5
a_counter{tag="bVal"} 6`))

			Expect(output.String()).Should(ContainSubstring(`# TYPE a_gauge gauge
a_gauge{tag="aVal"} 1
a_gauge{tag="bVal"} 2`))
		})
	})
})
