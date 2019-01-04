package exporter

import (
	"encoding/json"
	"eos_exporter/config"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	Float64 float64

	Exporter struct {
		Metrics map[string]*prometheus.Desc
		Config  *config.Config
	}

	AccountInfo struct {
		AccountName string `json:"account_name"`
		CPULimit    struct {
			Used Float64 `json:"used"`
			Max  Float64 `json:"max"`
		} `json:"cpu_limit"`
		NetLimit struct {
			Used Float64 `json:"used"`
			Max  Float64 `json:"max"`
		} `json:"net_limit"`
		RAMUsage         Float64 `json:"ram_usage"`
		RAMQuota         Float64 `json:"ram_quota"`
		CurrencyBalances map[string]Float64
	}
)

func (f *Float64) UnmarshalJSON(bs []byte) error {
	var i float64
	if err := json.Unmarshal(bs, &i); err == nil {
		*f = Float64(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return fmt.Errorf("expected a string or a float %v", err)
	}
	if err := json.Unmarshal([]byte(s), &i); err != nil {
		return err
	}
	*f = Float64(i)
	return nil
}
