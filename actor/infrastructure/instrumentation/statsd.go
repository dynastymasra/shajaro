package instrumentation

import (
	"fmt"

	"gopkg.in/alexcesaro/statsd.v2"
)

var (
	StatsD *statsd.Client
)

func InitiateStatsD(host, port, name string, enabled bool) error {
	if enabled {
		var err error

		address := fmt.Sprintf("%v:%v", host, port)

		StatsD, err = statsd.New(statsd.Address(address), statsd.Prefix(name))
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseStatsDClient() {
	if StatsD != nil {
		StatsD.Close()
	}
}

func NewTimingStatsD() *statsd.Timing {
	if StatsD != nil {
		t := StatsD.NewTiming()
		return &t
	}
	return nil
}

func TimingSend(bucket string, t *statsd.Timing) {
	if StatsD != nil {
		t.Send(bucket)
	}
}

func StatsDIncrement(bucket string) {
	if StatsD != nil {
		StatsD.Increment(bucket)
	}
}

func StatsDGauge(bucket string, value interface{}) {
	if StatsD != nil {
		StatsD.Gauge(bucket, value)
	}
}

func StatsDCount(bucket string, value interface{}) {
	if StatsD != nil {
		StatsD.Count(bucket, value)
	}
}

func StatsDFlush() {
	if StatsD != nil {
		StatsD.Flush()
	}
}

func StatsDHistogram(bucket string, value interface{}) {
	if StatsD != nil {
		StatsD.Histogram(bucket, value)
	}
}
