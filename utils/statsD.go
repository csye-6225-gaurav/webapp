package utils

import (
	"time"

	"github.com/smira/go-statsd"
)

var client *statsd.Client

func InitStatsD() {
	client = statsd.NewClient("localhost:8125",
		statsd.MaxPacketSize(1400),
		statsd.MetricPrefix("web."))
}

func CountIncrement(path string) {
	client.Incr("api.endpoint:"+path, 1)
}

func CountTimer(path string, startTime time.Time) {
	client.PrecisionTiming("api.endpoint.latency"+path, time.Since(startTime))
}
