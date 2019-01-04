package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.Metrics {
		ch <- m
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	var data, err = e.gatherAccountData()

	if err != nil {
		log.Errorf("Error gathering data from remote API %v", err)
		return
	}

	err = e.processMetrics(data, ch)

	if err != nil {
		log.Error("Error processing metrics", err)
		return
	}

	log.Infoln("All metrics were successfully gathered")
}
