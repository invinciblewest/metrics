package collectors

import (
	"github.com/invinciblewest/metrics/internal/storage"
	"log"
	"math/rand"
	"runtime"
)

type RuntimeCollector struct {
	st storage.Storage
}

func NewRuntimeCollector(st storage.Storage) *RuntimeCollector {
	return &RuntimeCollector{
		st: st,
	}
}

func (c *RuntimeCollector) Collect() error {
	log.Println("runtimeCollector: collecting metrics...")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	c.st.UpdateGauge("Alloc", float64(memStats.Alloc))
	c.st.UpdateGauge("BuckHashSys", float64(memStats.BuckHashSys))
	c.st.UpdateGauge("Frees", float64(memStats.Frees))
	c.st.UpdateGauge("GCCPUFraction", memStats.GCCPUFraction)
	c.st.UpdateGauge("GCSys", float64(memStats.GCSys))
	c.st.UpdateGauge("HeapAlloc", float64(memStats.HeapAlloc))
	c.st.UpdateGauge("HeapIdle", float64(memStats.HeapIdle))
	c.st.UpdateGauge("HeapInuse", float64(memStats.HeapInuse))
	c.st.UpdateGauge("HeapObjects", float64(memStats.HeapObjects))
	c.st.UpdateGauge("HeapReleased", float64(memStats.HeapReleased))
	c.st.UpdateGauge("HeapSys", float64(memStats.HeapSys))
	c.st.UpdateGauge("LastGC", float64(memStats.LastGC))
	c.st.UpdateGauge("Lookups", float64(memStats.Lookups))
	c.st.UpdateGauge("MCacheInuse", float64(memStats.MCacheInuse))
	c.st.UpdateGauge("MCacheSys", float64(memStats.MCacheSys))
	c.st.UpdateGauge("MSpanInuse", float64(memStats.MSpanInuse))
	c.st.UpdateGauge("MSpanSys", float64(memStats.MSpanSys))
	c.st.UpdateGauge("Mallocs", float64(memStats.Mallocs))
	c.st.UpdateGauge("NextGC", float64(memStats.NextGC))
	c.st.UpdateGauge("NumForcedGC", float64(memStats.NumForcedGC))
	c.st.UpdateGauge("NumGC", float64(memStats.NumGC))
	c.st.UpdateGauge("OtherSys", float64(memStats.OtherSys))
	c.st.UpdateGauge("PauseTotalNs", float64(memStats.PauseTotalNs))
	c.st.UpdateGauge("StackInuse", float64(memStats.StackInuse))
	c.st.UpdateGauge("StackSys", float64(memStats.StackSys))
	c.st.UpdateGauge("Sys", float64(memStats.Sys))
	c.st.UpdateGauge("TotalAlloc", float64(memStats.TotalAlloc))
	c.st.UpdateGauge("RandomValue", rand.Float64())
	c.st.UpdateCounter("PollCount", 1)

	pc, err := c.st.GetCounter("PollCount")
	if err != nil {
		return err
	}
	log.Printf("runtimeCollector: poll #%d is collected\r\n", pc)
	return nil
}
