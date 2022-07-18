package cli

import (
	"encoding/json"
	"time"
)

type TimeUnit string

const (
	TimeUnitDay TimeUnit = "day"
)

type DurationWithTimeUnit struct {
	*time.Duration
	timeUnit TimeUnit
}

func getDurationWithTimeUnit(duration time.Duration) DurationWithTimeUnit {
	p := DurationWithTimeUnit{
		Duration: &duration,
		timeUnit: TimeUnitDay,
	}
	return p
}

func (p *DurationWithTimeUnit) Present() float64 {
	presentByDay := func(duration *time.Duration) float64 {
		return duration.Hours() / float64(24)
	}
	return presentByDay(p.Duration)
}

func (p *DurationWithTimeUnit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&DurationOutput{
		Value: p.Present(),
		Unit:  string(p.timeUnit),
	})
}

func (p *DurationWithTimeUnit) UnmarshalJSON(data []byte) error {
	var out DurationOutput
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}
	p.timeUnit = TimeUnit(out.Unit)
	duration := time.Duration(float64(time.Hour) * float64(24) * out.Value)
	p.Duration = &duration
	return nil
}

type DurationOutput struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
