package collectors

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"log"
	"math/rand"
	"runtime"
)

type RuntimeCollector struct{}

func NewRuntimeCollector() *RuntimeCollector {
	return &RuntimeCollector{}
}

func (c *RuntimeCollector) Collect(st *storage.MemStorage) error {
	log.Println("runtimeCollector: collecting metrics...")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	st.UpdateGauge("Alloc", float64(memStats.Alloc))
	st.UpdateGauge("BuckHashSys", float64(memStats.BuckHashSys))
	st.UpdateGauge("Frees", float64(memStats.Frees))
	st.UpdateGauge("GCCPUFraction", memStats.GCCPUFraction)
	st.UpdateGauge("GCSys", float64(memStats.GCSys))
	st.UpdateGauge("HeapAlloc", float64(memStats.HeapAlloc))
	st.UpdateGauge("HeapIdle", float64(memStats.HeapIdle))
	st.UpdateGauge("HeapInuse", float64(memStats.HeapInuse))
	st.UpdateGauge("HeapObjects", float64(memStats.HeapObjects))
	st.UpdateGauge("HeapReleased", float64(memStats.HeapReleased))
	st.UpdateGauge("HeapSys", float64(memStats.HeapSys))
	st.UpdateGauge("LastGC", float64(memStats.LastGC))
	st.UpdateGauge("Lookups", float64(memStats.Lookups))
	st.UpdateGauge("MCacheInuse", float64(memStats.MCacheInuse))
	st.UpdateGauge("MCacheSys", float64(memStats.MCacheSys))
	st.UpdateGauge("MSpanInuse", float64(memStats.MSpanInuse))
	st.UpdateGauge("MSpanSys", float64(memStats.MSpanSys))
	st.UpdateGauge("Mallocs", float64(memStats.Mallocs))
	st.UpdateGauge("NextGC", float64(memStats.NextGC))
	st.UpdateGauge("NumForcedGC", float64(memStats.NumForcedGC))
	st.UpdateGauge("NumGC", float64(memStats.NumGC))
	st.UpdateGauge("OtherSys", float64(memStats.OtherSys))
	st.UpdateGauge("PauseTotalNs", float64(memStats.PauseTotalNs))
	st.UpdateGauge("StackInuse", float64(memStats.StackInuse))
	st.UpdateGauge("StackSys", float64(memStats.StackSys))
	st.UpdateGauge("Sys", float64(memStats.Sys))
	st.UpdateGauge("TotalAlloc", float64(memStats.TotalAlloc))
	st.UpdateGauge("RandomValue", rand.Float64())
	st.UpdateCounter("PollCount", 1)

	pc, err := st.GetCounter("PollCount")
	if err != nil {
		return err
	}
	log.Printf("runtimeCollector: poll #%d is collected\r\n", pc)
	return nil
}
