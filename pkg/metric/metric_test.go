package metric

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("metric package", func() {
	Context("JSONConverter functionality", func() {
		It("should be able to convert metric to JSON", func() {
			m := []*Metric{
				&Metric{Key: "test.success", Value: 1, Unit: "counter"},
				&Metric{Key: "test.failure", Value: 0, Unit: "counter"},
			}

			res, err := JSONConverter(m)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(ContainSubstring(`"key":"test.success"`))
			Expect(res).To(ContainSubstring(`"key":"test.failure"`))
		})
	})
})
