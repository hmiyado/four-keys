package cli

import (
	"encoding/json"
	"time"
)

type TimeUnit string

const (
	TimeUnitDay TimeUnit = "day"
)

type LeadTimeForChangesOutput struct {
	*time.Duration
	timeUnit TimeUnit
}

func getLeadTimeForChangesOutput(duration time.Duration) LeadTimeForChangesOutput {
	p := LeadTimeForChangesOutput{
		Duration: &duration,
		timeUnit: TimeUnitDay,
	}
	return p
}

func (p *LeadTimeForChangesOutput) Present() float64 {
	presentByDay := func(duration *time.Duration) float64 {
		return duration.Hours() / float64(24)
	}
	return presentByDay(p.Duration)
}

func (p *LeadTimeForChangesOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(&MetricsOutput{
		Value: p.Present(),
		Unit:  string(p.timeUnit),
	})
}

func (p *LeadTimeForChangesOutput) UnmarshalJSON(data []byte) error {
	var out MetricsOutput
	err := json.Unmarshal(data, &out)
	if err != nil {
		return err
	}
	p.timeUnit = TimeUnit(out.Unit)
	duration := time.Duration(float64(time.Hour) * float64(24) * out.Value)
	p.Duration = &duration
	return nil
}

type MetricsOutput struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
