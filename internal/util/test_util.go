package util

import "testing"

// isNearBy checks actual is in range of [expected*(1-epsilon), expected*(1+epsiolon)]
func IsNearBy(actual float64, expected float64, epsilon float64) bool {
	return actual >= expected*(1-epsilon) && actual <= expected*(1+epsilon)
}

// AssertIsNearBy checks actual is in range of [expected*(1-epsilon), expected*(1+epsiolon)]
//
func AssertIsNearBy(t *testing.T, actual float64, expected float64, epsilon float64) {
	if IsNearBy(actual, expected, epsilon) {
		return
	}
	t.Errorf("actual should be in [%v, %v] but %v", expected*(1-epsilon), expected*(1+epsilon), actual)
}
