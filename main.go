package main

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	logCacheAPI = kingpin.Flag("log-cache-api", "The log-cache API URL.").Required().OverrideDefaultFromEnvar("LOG_CACHE_API").Short('a').String()

	port    = kingpin.Flag("port", "The port server should be running on.").Default("8080").OverrideDefaultFromEnvar("PORT").Short('p').Int64()
	verbose = kingpin.Flag("verbose", "Run the server in debugging mode.").Default("false").OverrideDefaultFromEnvar("DEBUG").Short('v').Bool()
)

var acceptFormats = []responder{
	responder{
		accept:      "text/plain",
		contentType: "text/plain; version=0.0.4; charset=utf-8",
	},
}

func main() {
	kingpin.Parse()

	log := logrus.New()

	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
		log.Debug("Verbose mode enabled")
	}

	srv, err := newServer(log, acceptFormats, *logCacheAPI)
	if err != nil {
		log.Panic(err)
	}

	srv.routes()
	srv.run(*port)
}
