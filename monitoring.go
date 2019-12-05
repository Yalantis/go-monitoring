// Package `monitoring` used for collecting and sending in-app Go stats to InfluxDB
package monitoring

import (
	"errors"
	"fmt"
	"time"

	ginflux "github.com/Yalantis/go-influx"
)

// Logger interface
type Logger interface {
	Warn(...interface{})
}

// Influx interface
type Influx interface {
	RegisterMeasurement(...ginflux.Measurement)
	Push(string, map[string]string, map[string]interface{}) error
}

// MeasurementOption type
type MeasurementOption func(*ginflux.Measurement)

// errors
var ErrInfluxInit = errors.New("influx is not initialized")

// const
const DefaultTimeout = time.Second

// Monitoring sends metrics to influx
type Monitoring struct {
	timeout    time.Duration
	appName    string
	collectors []Collector

	logger Logger
	influx Influx

	shutdown chan struct{}
}

// New creates Monitoring
func New(timeout time.Duration, appName string, influx Influx) (*Monitoring, error) {
	if influx == nil {
		return nil, ErrInfluxInit
	}

	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	return &Monitoring{
		timeout:  timeout,
		appName:  appName,
		influx:   influx,
		shutdown: make(chan struct{}),
	}, nil
}

// SetLogger sets logger
func (c *Monitoring) SetLogger(logger Logger) {
	c.logger = logger
}

// Start Monitoring
func (c *Monitoring) Start() {
	go func() {
		t := time.NewTicker(c.timeout)

	loop:
		for {
			select {
			case <-t.C:
				for _, cc := range c.collectors {
					name := cc.GetName()
					tags, fields := cc.Collect()
					if tags == nil {
						tags = make(map[string]string)
					}
					tags["app"] = c.appName

					err := c.influx.Push(name, tags, fields)
					if err != nil && c.logger != nil {
						c.logger.Warn(fmt.Sprintf("Influx[%s]", name), "error", err)
					}
				}
			case <-c.shutdown:
				break loop
			}
		}

		t.Stop()
	}()
}

// Shutdown Monitoring
func (c *Monitoring) Shutdown() {
	close(c.shutdown)
}

// AddCollector adds collector
func (c *Monitoring) AddCollector(collector Collector) {
	c.collectors = append(c.collectors, collector)
}

// InitDefault initializes Monitoring with default options, measurements
func InitDefault(timeout time.Duration, appName string, influx Influx, opts ...MeasurementOption) (*Monitoring, error) {
	monitoring, err := New(timeout, appName, influx)
	if err != nil {
		return nil, err
	}

	defaultMeasurement := ginflux.Measurement{
		Name:            AppStatsName,
		Database:        AppStatsDatabase,
		RetentionPolicy: ginflux.ShortTermRP,
		QueueSize:       ginflux.DefaultQueueSize,
		FlushInterval:   ginflux.DefaultFlushInterval,
	}

	for _, opt := range opts {
		opt(&defaultMeasurement)
	}

	influx.RegisterMeasurement(defaultMeasurement)

	monitoring.AddCollector(&AppStatistics{Name: AppStatsName})
	monitoring.Start()

	return monitoring, nil
}
