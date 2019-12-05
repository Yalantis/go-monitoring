package monitoring

import (
	"reflect"
	"sync"
	"testing"
	"time"

	ginflux "github.com/Yalantis/go-influx"
)

type fakeLogger struct {
	WarnFn func()
}

func (f fakeLogger) Warn(...interface{}) {
	if f.WarnFn != nil {
		f.WarnFn()
	}
}

type fakeInflux struct {
	RegisterMeasurementFn func([]ginflux.Measurement)
	PushFn                func(string, map[string]string, map[string]interface{}) error
}

func (f fakeInflux) RegisterMeasurement(m ...ginflux.Measurement) {
	if f.RegisterMeasurementFn != nil {
		f.RegisterMeasurementFn(m)
	}
}

func (f fakeInflux) Push(name string, tags map[string]string, fields map[string]interface{}) error {
	if f.PushFn != nil {
		return f.PushFn(name, tags, fields)
	}
	return nil
}

func TestInitDefault(t *testing.T) {
	// invalid
	monitoring, err := InitDefault(time.Second, "", nil)
	if err != ErrInfluxInit {
		t.Fatalf("expected error. want: %v, got: %v", ErrInfluxInit, err)
	}
	if monitoring != nil {
		t.Fatalf("expected to be nil. got: %v", monitoring)
	}

	// ok
	monitoring, err = InitDefault(time.Second, "", fakeInflux{})
	if err != nil {
		t.Fatalf("expected error to be nil. got: %v", err)
	}
	if monitoring == nil {
		t.Fatalf("expected not to be nil. got: %v", monitoring)
	}
	monitoring.Shutdown()

	// applies MeasurementOption
	var measurement ginflux.Measurement
	influx := fakeInflux{
		RegisterMeasurementFn: func(m []ginflux.Measurement) {
			measurement = m[0]
		},
	}

	option := func(measurement *ginflux.Measurement) {
		measurement.Name = "measurement.Name"
		measurement.Database = "measurement.Database"
		measurement.RetentionPolicy = "measurement.RetentionPolicy"
		measurement.QueueSize = 42
		measurement.FlushInterval = 42
	}

	monitoring, err = InitDefault(time.Second, "", influx, option)
	if err != nil {
		t.Fatalf("expected error to be nil. got: %v", err)
	}
	if monitoring == nil {
		t.Fatalf("expected not to be nil. got: %v", monitoring)
	}
	monitoring.Shutdown()

	if measurement.Name != "measurement.Name" {
		t.Fatalf("expected to match. want: %v, got: %v", measurement.Name, "measurement.Name")
	}

	if measurement.Database != "measurement.Database" {
		t.Fatalf("expected to match. want: %v, got: %v", measurement.Database, "measurement.Database")
	}

	if measurement.RetentionPolicy != "measurement.RetentionPolicy" {
		t.Fatalf("expected to match. want: %v, got: %v", measurement.RetentionPolicy, "measurement.RetentionPolicy")
	}

	if measurement.QueueSize != 42 {
		t.Fatalf("expected to match. want: %v, got: %v", measurement.QueueSize, 42)
	}

	if measurement.FlushInterval != 42 {
		t.Fatalf("expected to match. want: %v, got: %v", measurement.FlushInterval, 42)
	}
}

func TestNew(t *testing.T) {
	monitoring, _ := New(0, "", fakeInflux{})
	if monitoring.timeout == 0 {
		t.Fatal("expected not to be zero")
	}
}

func TestInfluxDB_SetLogger(t *testing.T) {
	monitoring, _ := New(0, "", fakeInflux{})
	defer monitoring.Shutdown()

	var logged bool

	monitoring.SetLogger(fakeLogger{
		WarnFn: func() {
			logged = true
		},
	})

	monitoring.logger.Warn()

	if logged != true {
		t.Fatal("expected to be true")
	}
}

func TestMonitoring_Start(t *testing.T) {
	var once sync.Once
	done := make(chan struct{})

	expectedTags := map[string]string{"app": "app"}
	expectedFields := make(map[string]interface{})

	var actualName string
	var actualTags map[string]string
	var actualFields map[string]interface{}

	pushFn := func(name string, tags map[string]string, fields map[string]interface{}) error {
		actualName = name
		actualTags = tags
		actualFields = fields
		once.Do(func() { close(done) })
		return nil
	}

	monitoring, _ := InitDefault(time.Millisecond, "app", fakeInflux{PushFn: pushFn})
	monitoring.Start()

	select {
	case <-time.After(time.Second):
		t.Fatal("out of time")
	case <-done:
		monitoring.Shutdown()
	}

	if actualName != AppStatsName {
		t.Fatalf("expected to match. want: %v, got: %v", AppStatsName, actualName)
	}

	if !reflect.DeepEqual(expectedTags, actualTags) {
		t.Fatalf("expected to match. want: %v, got: %v", expectedTags, actualTags)
	}

	if reflect.DeepEqual(expectedFields, actualFields) {
		t.Fatalf("expected to not match. want: %v, got: %v", expectedFields, actualFields)
	}
}
