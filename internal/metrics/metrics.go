package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TargetsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pentscan_targets_processed_total",
		Help: "The total number of processed targets (host:port)",
	})
	PortsOpen = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pentscan_ports_open_total",
		Help: "The total number of ports found to be open",
	})
	PortsClosed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pentscan_ports_closed_total",
		Help: "The total number of ports found to be closed",
	})
	PortsFiltered = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pentscan_ports_filtered_total",
		Help: "The total number of ports found to be filtered",
	})
)
