package main

import (
	"flag"
	"net/http"

	"eos_exporter/config"
	"eos_exporter/exporter"

	log "github.com/prometheus/common/log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	mets       map[string]*prometheus.Desc
	configFile = flag.String("config.file", "config.yml", "Eos exporter configuration file.")
)

func main() {
	flag.Parse()
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Errorf("Error reading config file. Error %s", err)
		return
	}
	mets = exporter.AddMetrics(cfg.Tokens)
	exporter := exporter.Exporter{
		Metrics: mets,
		Config:  *cfg,
	}
	prometheus.MustRegister(&exporter)
	http.Handle("/metrics", prometheus.Handler())
	log.Info("Beginning to server on port 9386")
	log.Fatal(http.ListenAndServe(":9386", nil))
}
