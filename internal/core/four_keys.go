package core

import "time"

func GetDeploymentFrequency(releases []*Release, option Option) float64 {
	duration := option.Until.Sub(option.Since)
	daysCount := int(duration.Hours() / 24)
	releasesCount := len(releases)
	return float64(releasesCount) / float64(daysCount)
}

func GetMeanLeadTimeForChanges(releases []*Release) time.Duration {
	if len(releases) == 0 {
		return time.Duration(0)
	}
	sum := time.Duration(0)
	for _, release := range releases {
		sum = release.LeadTimeForChanges + sum
	}
	return time.Duration(int64(sum) / int64(len(releases)))
}

func GetTimeToRestore(releases []*Release) time.Duration {
	sum := time.Duration(0)
	countOfRestore := 0
	for _, release := range releases {
		if release.Result.TimeToRestore != nil {
			sum += *release.Result.TimeToRestore
		}
	}
	if countOfRestore == 0 {
		return sum
	}
	return sum / time.Duration(countOfRestore)
}

func GetChangeFailureRate(releases []*Release) float64 {
	if len(releases) == 0 {
		return 0
	}

	sumOfFailure := 0
	for _, release := range releases {
		if !release.Result.IsSuccess {
			sumOfFailure += 1
		}
	}
	return float64(sumOfFailure) / float64(len(releases))
}
