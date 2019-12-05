package monitoring

// Collector interface
type Collector interface {
	GetName() string
	Collect() (tags map[string]string, fields map[string]interface{})
}

// Collect func shortcut
type Collect func() (tags map[string]string, fields map[string]interface{})

// StatsCollector responsible for stats collecting
type StatsCollector struct {
	name string
	fn   Collect
}

// NewCollector creates Collector
func NewCollector(name string, fn Collect) Collector {
	return StatsCollector{
		name: name,
		fn:   fn,
	}
}

// GetName returns name of Collector
func (s StatsCollector) GetName() string { return s.name }

// Collect executes underling fn
func (s StatsCollector) Collect() (map[string]string, map[string]interface{}) {
	return s.fn()
}
