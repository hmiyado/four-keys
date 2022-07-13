package util

// isNearBy checks actual is in range of [expected*(1-epsilon), expected*(1+epsiolon)]
func IsNearBy(actual float64, expected float64, epsilon float64) bool {
	return actual >= expected*(1-epsilon) && actual <= expected*(1+epsilon)
}
