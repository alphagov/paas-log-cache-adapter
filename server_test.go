package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
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
				accept:      "text/plain",
				contentType: "text/plain",
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

		It("should return nothing if there are no metrics", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/v1/meta", api),
				httpmock.NewStringResponder(http.StatusOK, `{"meta":{}}`))

			srv, err := newServer(log, acceptedFormats, api)

			Expect(err).NotTo(HaveOccurred())
			srv.routes()
			r, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			r.Header.Add("Authorization", "__JWT_ACCESS_TOKEN__")
			r.Header.Add("Accept", "text/plain")

			w := httptest.NewRecorder()
			srv.recordRequest(srv.requireAuth(srv.chooseFormat(srv.handleMetrics())))(w, r)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain"))
			Expect(w.Body.String()).To(Equal(""))
		})

		It("should serve the requested data correctly", func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/v1/meta", api),
				httpmock.NewStringResponder(http.StatusOK, `{"meta":{"test":{}}}`))

			httpmock.RegisterResponder("GET", fmt.Sprintf("%s/v1/read/test", api),
				httpmock.NewStringResponder(http.StatusOK, `{
  "envelopes": {
    "batch": [
      {
        "source_id": "instance-a",
        "tags": {
          "tag-a": "val-a"
        },
        "gauge": {
          "metrics": {
            "a-gauge": {
              "value": 3.14
            }
          }
        }
      },
      {
        "source_id": "instance-a",
        "tags": {
          "tag-a": "val-a"
        },
        "counter": {
          "name": "counter",
          "total": 8
        }
      },
      {
        "source_id": "instance-b",
        "tags": {
          "tag-b": "val-b"
        },
        "counter": {
          "name": "counter",
          "total": 10
        }
      }
    ]
  }
}`))

			srv, err := newServer(log, acceptedFormats, api)

			Expect(err).NotTo(HaveOccurred())
			srv.routes()
			r, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			r.Header.Add("Authorization", "__JWT_ACCESS_TOKEN__")
			r.Header.Add("Accept", "text/plain")

			w := httptest.NewRecorder()
			srv.recordRequest(srv.requireAuth(srv.chooseFormat(srv.handleMetrics())))(w, r)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Header().Get("Content-Type")).To(Equal("text/plain"))
			Expect(w.Body.String()).To(ContainSubstring(`# TYPE counter counter
counter{instance_id="instance-a",tag_a="val-a"} 8
counter{instance_id="instance-b",tag_b="val-b"} 10`))

			Expect(w.Body.String()).To(ContainSubstring(`# TYPE a_gauge gauge
a_gauge{instance_id="instance-a",tag_a="val-a"} 3.14`))
		})
	})
})
