package utils

import (
	"time"

	"github.com/smira/go-statsd"
)

var Client *statsd.Client

func InitStatsD() {
	Client = statsd.NewClient("localhost:8125",
		statsd.MaxPacketSize(1400),
		statsd.MetricPrefix("web."))
}

func CountIncrement(path string) {
	Client.Incr("api.endpoint.count"+path, 1)
}

func CountTimer(path string, startTime time.Time) {
	Client.PrecisionTiming("api.endpoint.latency"+path, time.Since(startTime))
}
