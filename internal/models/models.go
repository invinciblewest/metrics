package models

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func CheckType(mType string) bool {
	return mType == TypeGauge || mType == TypeCounter
}

func (m *Metrics) CheckType() bool {
	return CheckType(m.MType)
}
