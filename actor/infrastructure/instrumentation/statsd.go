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

func NewTimingStatsD(enabled bool) *statsd.Timing {
	if enabled && StatsD != nil {
		t := StatsD.NewTiming()
		return &t
	}
	return nil
}

func TimingSend(key string, t *statsd.Timing, enabled bool) {
	if enabled && StatsD != nil {
		t.Send(key)
	}
}

func StatsDIncrement(key string, enabled bool) {
	if enabled && StatsD != nil {
		StatsD.Increment(key)
	}
}

func StatsDGauge(key string, value interface{}, enabled bool) {
	if enabled && StatsD != nil {
		StatsD.Gauge(key, value)
	}
}
