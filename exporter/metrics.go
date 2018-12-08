package exporter

import "github.com/prometheus/client_golang/prometheus"

func AddMetrics() map[string]*prometheus.Desc {
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
	return metrics
}

func (e *Exporter) processMetrics(data []*AccountInfo, ch chan<- prometheus.Metric) error {
	for _, x := range data {
		ch <- prometheus.MustNewConstMetric(e.Metrics["CpuUsed"], prometheus.GaugeValue, x.CPULimit.Used, x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["CpuMax"], prometheus.GaugeValue, x.CPULimit.Max, x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["NetUsed"], prometheus.GaugeValue, x.NetLimit.Used, x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["NetMax"], prometheus.GaugeValue, x.NetLimit.Max, x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["RamUsed"], prometheus.GaugeValue, x.RAMUsage, x.AccountName)
		ch <- prometheus.MustNewConstMetric(e.Metrics["RamQuota"], prometheus.GaugeValue, x.RAMQuota, x.AccountName)
	}
	return nil
}
