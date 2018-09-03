package prometheus

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("prometheus package", func() {
	Context("Sanitize functionality", func() {
		It("does not mangle good things", func() {

			input := "a_good_metric"
			Expect(Sanitize(input)).Should(Equal(input))
		})

		It("replaces periods", func() {
			input := "a.metric.with_periods"
			Expect(Sanitize(input)).Should(Equal("a_metric_with_periods"))
		})

		It("replaces dashes", func() {
			input := "a-metric-with_dashes"
			Expect(Sanitize(input)).Should(Equal("a_metric_with_dashes"))
		})

		It("does not replace colons", func() {
			input := ":a_metric:with:colons"
			Expect(Sanitize(input)).Should(Equal(":a_metric:with:colons"))
		})

		It("inserts a leading character in front of leading numbers", func() {
			input := "123metric_456"
			Expect(Sanitize(input)).Should(Equal("_123metric_456"))
		})
	})
})
