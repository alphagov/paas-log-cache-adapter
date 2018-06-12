package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/alphagov/paas-log-cache-adapter/pkg/metric"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func faultyConverter(m []*metric.Metric) ([]byte, error) {
	return nil, fmt.Errorf("__EVERYTHING_IS_UNDER_CONTROL__")
}

var _ = Describe("main package", func() {
	var log *logrus.Logger
	var acceptedFormats []responder
	api := "https://example.com"

	BeforeEach(func() {
		log = logrus.New()
		log.Out = GinkgoWriter
		acceptedFormats = []responder{
			responder{
				accept:      "application/json",
				contentType: "application/json",
				converter:   metric.JSONConverter,
			},
			responder{
				accept:      "apocalypse",
				contentType: "application/json",
				converter:   faultyConverter,
			},
		}
	})

	Context("server", func() {
		It("should fail to create server due to missing logger", func() {
			_, err := newServer(nil, acceptedFormats, api)

			Expect(err).To(HaveOccurred())
		})

		It("should fail to serve the route due to missing auth", func() {
			srv, err := newServer(log, acceptedFormats, api)

			Expect(err).NotTo(HaveOccurred())
			srv.routes()
			r, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			w := httptest.NewRecorder()
			srv.requireAuth(srv.chooseFormat(srv.handleMetrics()))(w, r)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
		})

		It("should fail to serve the route due unsupported format", func() {
			srv, err := newServer(log, acceptedFormats, api)

			Expect(err).NotTo(HaveOccurred())
			srv.routes()
			r, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			r.Header.Add("Authorization", "__JWT_ACCESS_TOKEN__")
			r.Header.Add("Accept", "text/html")

			w := httptest.NewRecorder()
			srv.requireAuth(srv.chooseFormat(srv.handleMetrics()))(w, r)

			Expect(w.Code).To(Equal(http.StatusNotAcceptable))
		})

		It("should fail to serve the route due to faulty converter", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/v1/meta", api),
				httpmock.NewStringResponder(http.StatusOK, `{"meta":{"test":{}}}`))
			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/v1/read/test", api),
				httpmock.NewStringResponder(http.StatusOK, `{"envelopes":{"batch":[{}]}}`))

			srv, err := newServer(log, acceptedFormats, api)

			Expect(err).NotTo(HaveOccurred())
			srv.routes()
			r, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			r.Header.Add("Authorization", "__JWT_ACCESS_TOKEN__")
			r.Header.Add("Accept", "apocalypse")

			w := httptest.NewRecorder()
			srv.requireAuth(srv.chooseFormat(srv.handleMetrics()))(w, r)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
		})

		It("should serve the requested data correctly", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/v1/meta", api),
				httpmock.NewStringResponder(http.StatusOK, `{"meta":{"test":{}}}`))

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/v1/read/test", api),
				httpmock.NewStringResponder(http.StatusOK, `{"envelopes":{"batch":[{}]}}`))

			srv, err := newServer(log, acceptedFormats, api)

			Expect(err).NotTo(HaveOccurred())
			srv.routes()
			r, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			r.Header.Add("Authorization", "__JWT_ACCESS_TOKEN__")
			r.Header.Add("Accept", "application/json")

			w := httptest.NewRecorder()
			srv.recordRequest(srv.requireAuth(srv.chooseFormat(srv.handleMetrics())))(w, r)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(w.Body).To(ContainSubstring(`[]`))
		})
	})
})
