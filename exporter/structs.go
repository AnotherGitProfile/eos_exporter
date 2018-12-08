package exporter

import (
	"eos_exporter/config"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	Metrics map[string]*prometheus.Desc
	config.Config
}

type AccountInfo struct {
	AccountName string `json:"account_name"`
	CPULimit    struct {
		Used float64 `json:"used"`
		Max  float64 `json:"max"`
	} `json:"cpu_limit"`
	NetLimit struct {
		Used float64 `json:"used"`
		Max  float64 `json:"max"`
	} `json:"net_limit"`
	RAMUsage float64 `json:"ram_usage"`
	RAMQuota float64 `json:"ram_quota"`
}
