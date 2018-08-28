package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/alphagov/paas-log-cache-adapter/pkg/prometheus"
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
		w.Header().Add("Content-Type", "text/plain")

		token := r.Header.Get("Authorization")
		client := logcache.NewClient(s.logCacheAPI, logcache.WithHTTPClient(&myClient{
			token: token,
		}))

		ctx := context.Background()

		meta, err := client.Meta(ctx)
		if err != nil {
			s.logger.Error(err)
			s.error(
				w,
				http.StatusInternalServerError,
				"Cannot connect to log-cache",
			)
			return
		}

		var logGetters, appender sync.WaitGroup
		sourceIDs := make(chan string)
		allEnvelopes := make([]*loggregator_v2.Envelope, 0)
		envelopeChan := make(chan []*loggregator_v2.Envelope)

		for i := 0; i < 10; i++ {
			logGetters.Add(1)
			go func() {
				defer logGetters.Done()

				for sourceID := range sourceIDs {
					s.logger.WithFields(logrus.Fields{
						"instance_id": sourceID,
					}).Debug("Obtaining metrics for resource")

					envelopes, err := client.Read(
						ctx,
						sourceID,
						time.Now().Add(-10*time.Minute),
					)

					if err != nil {
						s.logger.Error(err)
						continue
					}

					envelopeChan <- envelopes
				}
			}()
		}

		appender.Add(1)
		go func() {
			defer appender.Done()
			for envelopes := range envelopeChan {
				allEnvelopes = append(allEnvelopes, envelopes...)
			}
		}()

		for sourceID := range meta {
			sourceIDs <- sourceID
		}
		close(sourceIDs)

		logGetters.Wait()
		close(envelopeChan)
		appender.Wait()

		metricFams := prometheus.Convert(allEnvelopes)
		err = prometheus.WriteMetrics(metricFams, w)
		if err != nil {
			s.logger.Error(err)
			s.error(
				w,
				http.StatusInternalServerError,
				"Error converting log-cache metrics to prometheus format",
			)
			return
		}

	}
}

func (s *server) error(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
