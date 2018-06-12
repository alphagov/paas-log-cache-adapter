package main

import (
	"fmt"
	"net/http"

	"github.com/alphagov/paas-log-cache-adapter/pkg/metric"
	"github.com/sirupsen/logrus"
)

type responder struct {
	accept      string
	contentType string
	converter   metric.Converter
}

type server struct {
	acceptedFormats []responder
	logger          *logrus.Logger

	logCacheAPI string
	router      *http.ServeMux
	responder   *responder
}

func newServer(logger *logrus.Logger, acceptedFormats []responder, logCacheAPI string) (*server, error) {
	if logger == nil {
		return nil, fmt.Errorf("server: logger is required")
	}

	return &server{
		acceptedFormats: acceptedFormats,
		logger:          logger,

		logCacheAPI: logCacheAPI,
		router:      http.NewServeMux(),
	}, nil
}

func (s *server) routes() {
	s.router.HandleFunc("/metrics", s.recordRequest(s.requireAuth(s.chooseFormat(s.handleMetrics()))))
}

func (s *server) run(port int64) {
	s.logger.WithFields(logrus.Fields{
		"url": fmt.Sprintf("http://localhost:%d/", port),
	}).Info("Starting server")
	s.logger.Panic(http.ListenAndServe(fmt.Sprintf(":%d", port), s.router))
}
