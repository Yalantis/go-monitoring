package monitoring

import "runtime"

const (
	AppStatsName     = "go_memstats"
	AppStatsDatabase = "go_statistics"
)

// AppStatistics AppStatistics
type AppStatistics struct {
	Name     string
	memStats runtime.MemStats
}

// GetName returns name of AppStatistics
func (a AppStatistics) GetName() string { return a.Name }

// Collect collects app statistics
func (a *AppStatistics) Collect() (map[string]string, map[string]interface{}) {
	// collect measurements
	runtime.ReadMemStats(&a.memStats)

	fields := map[string]interface{}{
		"alloc":         int(a.memStats.Alloc),
		"total_alloc":   int(a.memStats.TotalAlloc),
		"sys":           int(a.memStats.Sys),
		"mallocs":       int(a.memStats.Mallocs),
		"frees":         int(a.memStats.Frees),
		"heap_alloc":    int(a.memStats.HeapAlloc),
		"heap_sys":      int(a.memStats.HeapSys),
		"heap_objects":  int(a.memStats.HeapObjects),
		"next_gc":       int(a.memStats.NextGC),
		"num_goroutine": runtime.NumGoroutine(),
	}

	return nil, fields
}
