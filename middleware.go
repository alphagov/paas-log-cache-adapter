package main

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func (s *server) chooseFormat(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accept := r.Header.Get("Accept")

		for _, af := range s.acceptedFormats {
			if strings.Contains(accept, af.accept) {
				s.responder = &af
				break
			}
		}

		if s.responder == nil {
			s.error(w, http.StatusNotAcceptable, "unsupported Accept header provided")
			return
		}

		handler(w, r)
	}
}

func (s *server) recordRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.WithFields(logrus.Fields{
			"accept": r.Header.Get("Accept"),
			"auth":   r.Header.Get("Authorization") != "",
			"path":   r.URL.EscapedPath(),
		}).Debug("Request made")

		handler(w, r)
	}
}

func (s *server) requireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			s.error(w, http.StatusUnauthorized, "you need to provide Authorization header")
			return
		}

		handler(w, r)
	}
}
