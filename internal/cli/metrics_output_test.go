package cli

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPresentShouldReturnDaysWhenTimeUnitIsDay(t *testing.T) {
	day := time.Duration(time.Hour * 24)
	actual := &LeadTimeForChangesOutput{
		Duration: &day,
		timeUnit: TimeUnitDay,
	}
	if actual.Present() != float64(1) {
		t.Errorf("expected: 1, actual: %v", actual.Present())
	}
}

func TestPresentShouldEqualMarshalAndUnmarshalJSON(t *testing.T) {
	day := time.Duration(time.Hour * 24)
	expected := &LeadTimeForChangesOutput{
		Duration: &day,
		timeUnit: TimeUnitDay,
	}
	marshaled, err := json.Marshal(expected)
	if err != nil {
		t.Errorf(err.Error())
	}

	var actual LeadTimeForChangesOutput
	unmarshalError := json.Unmarshal(marshaled, &actual)
	if unmarshalError != nil {
		t.Errorf(err.Error())
	}

	if expected.Duration.Seconds() == actual.Duration.Seconds() && expected.timeUnit == actual.timeUnit {
		return
	}
	t.Errorf("expected: %v, actual: %v", expected, actual)
}
