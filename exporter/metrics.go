package exporter

import (
	"eos_exporter/config"

	"github.com/prometheus/client_golang/prometheus"
)

func AddMetrics(tokens []config.TokenContract) map[string]*prometheus.Desc {
	metrics := make(map[string]*prometheus.Desc)
	metrics["CpuUsed"] = prometheus.NewDesc(
		prometheus.BuildFQName("eos", "account", "cpu_used"),
		"Current value of used CPU for given account",
		[]string{"account"}, nil,
	)
	metrics["CpuMax"] = prometheus.NewDesc(
		prometheus.BuildFQName("eos", "account", "cpu_max"),
		"Maximum amount of CPU that can be used by given account",
		[]string{"account"}, nil,
	)
	metrics["NetUsed"] = prometheus.NewDesc(
		prometheus.BuildFQName("eos", "account", "net_used"),
		"Current value of used NET for given account",
		[]string{"account"}, nil,
	)
	metrics["NetMax"] = prometheus.NewDesc(
		prometheus.BuildFQName("eos", "account", "net_max"),
		"Maximum amount of NET than can be used by given account",
		[]string{"account"}, nil,
	)
	metrics["RamUsed"] = prometheus.NewDesc(
		prometheus.BuildFQName("eos", "account", "ram_used"),
		"Total amount of used ram for given account",
		[]string{"account"}, nil,
	)
	metrics["RamQuota"] = prometheus.NewDesc(
		prometheus.BuildFQName("eos", "account", "ram_quota"),
		"Amount of available ram for given account",
		[]string{"account"}, nil,
	)
	metrics["CurrencyBalance"] = prometheus.NewDesc(
		prometheus.BuildFQName("eos", "account", "balance"),
		"Currency balance for given account",
		[]string{"account", "token"}, nil,
	)
	return metrics
}

func (e *Exporter) processMetrics(data []*AccountInfo, ch chan<- prometheus.Metric) error {
	for _, x := range data {
		ch <- prometheus.MustNewConstMetric(e.Metrics["CpuUsed"], prometheus.GaugeValue, float64(x.CPULimit.Used), x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["CpuMax"], prometheus.GaugeValue, float64(x.CPULimit.Max), x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["NetUsed"], prometheus.GaugeValue, float64(x.NetLimit.Used), x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["NetMax"], prometheus.GaugeValue, float64(x.NetLimit.Max), x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["RamUsed"], prometheus.GaugeValue, float64(x.RAMUsage), x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["RamQuota"], prometheus.GaugeValue, float64(x.RAMQuota), x.AccountName)
		for tokenSymbol, balance := range x.CurrencyBalances {
			ch <- prometheus.MustNewConstMetric(e.Metrics["CurrencyBalance"], prometheus.GaugeValue, float64(balance), x.AccountName, tokenSymbol)
		}
	}
	return nil
}
