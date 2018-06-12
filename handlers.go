package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"code.cloudfoundry.org/go-log-cache"
	"github.com/alphagov/paas-log-cache-adapter/pkg/metric"
	"github.com/sirupsen/logrus"
)

type myClient struct {
	token string
}

func (mC *myClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", mC.token)

	c := http.Client{}

	return c.Do(req)
}

func (s *server) handleMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		m := []*metric.Metric{}
		token := r.Header.Get("Authorization")
		client := logcache.NewClient(s.logCacheAPI, logcache.WithHTTPClient(&myClient{
			token: token,
		}))

		ctx := context.Background()

		meta, err := client.Meta(ctx)
		if err != nil {
			s.logger.Error(err)
			s.error(w, http.StatusInternalServerError, "we're experiencing issues, sorry")
			return
		}

		// TODO: Consider optimising that. It runs fairly slow (almost 4s) with
		// 745 instances running on my admin account :troll:
		for sourceID := range meta {
			s.logger.WithFields(logrus.Fields{
				"instance_id": sourceID,
			}).Debug("Obtaining metrics for resource")

			wg.Add(1)
			go func(sourceID string) {
				defer wg.Done()
				r, err := client.Read(ctx, sourceID, time.Now().Add(-10*time.Minute))
				if err != nil {
					s.logger.Error(err)
					return
				}
				m = append(m, convertToMetrics(r)...)
			}(sourceID)
		}
		wg.Wait()

		data, err := s.responder.converter(m)
		if err != nil {
			s.logger.Error(err)
			s.error(w, http.StatusInternalServerError, "we're experiencing issues, sorry")
			return
		}
		w.Header().Add("Content-Type", s.responder.contentType)

		w.Write(data)
	}
}

func (s *server) error(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
