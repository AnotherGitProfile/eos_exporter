package main

import (
	"flag"
	"fmt"
	"net/http"

	"eos_exporter/config"
	"eos_exporter/exporter"

	log "github.com/prometheus/common/log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	mets map[string]*prometheus.Desc
	argv struct {
		configPath string
		port       uint
		help       bool
	}
)

func init() {
	flag.UintVar(&argv.port, "port", 9386, "port to listen")
	flag.BoolVar(&argv.help, "h", false, "show this help")
	flag.StringVar(&argv.configPath, "config.file", "config.yml", "path to configuration file")
	flag.Parse()
}

func main() {
	if argv.help {
		flag.Usage()
		return
	}
	cfg, err := config.LoadConfig(argv.configPath)
	if err != nil {
		log.Errorf("Error reading config file. Error %s", err)
		return
	}
	mets = exporter.AddMetrics(cfg.Tokens)
	exporter := exporter.Exporter{
		Metrics: mets,
		Config:  cfg,
	}
	prometheus.MustRegister(&exporter)
	http.Handle("/metrics", prometheus.Handler())
	log.Infof("Beginning to server on port %d", argv.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", argv.port), nil))
}
