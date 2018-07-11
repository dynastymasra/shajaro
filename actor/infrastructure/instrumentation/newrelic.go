package instrumentation

import "github.com/newrelic/go-agent"

var newRelic *Newrelic

type Newrelic struct {
	application newrelic.Application
	config      newrelic.Config
}

func InitiateNewrelic(name, license string, enabled bool) error {
	config := newrelic.NewConfig(name, license)
	config.Enabled = enabled

	app, err := newrelic.NewApplication(config)
	if err != nil {
		return err
	}

	newRelic = &Newrelic{
		application: app,
		config:      config,
	}

	return nil
}

func NewRelicConfig() newrelic.Config {
	return newRelic.config
}

func NewRelicApp() newrelic.Application {
	return newRelic.application
}
