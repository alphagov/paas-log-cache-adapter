package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPaasMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PaasMetrics Suite")
}
