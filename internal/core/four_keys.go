package core

import "time"

func GetDeploymentFrequency(releases []*Release, option Option) float64 {
	return GetDeploymentFrequencyByTimeunit(releases, option.Since, option.Until, "day")
}

// GetDeploymentFrequencyByTimeunit returns deployment frequency by timeunit.
// timeunitKind is one of "day", "week", "month". Default is "day".
// If timeunitKind is not one of them, "day" is used.
func GetDeploymentFrequencyByTimeunit(releases []*Release, since time.Time, until time.Time, timeunitKind string) float64 {
	duration := until.Sub(since)
	frequency := float64(duration / (time.Hour * 24))
	switch timeunitKind {
	case "day":
		frequency = float64(duration / (time.Hour * 24))
		if frequency == 0 {
			frequency = 1
		}
	case "week":
		frequency = float64(duration / (time.Hour * 24 * 7))
		if frequency == 0 {
			frequency = 1
		}
	case "month":
		frequency = float64((until.Year()-since.Year())*12 + int(until.Month()-since.Month()))
		if frequency == 0 {
			frequency = 1
		}
	}

	releasesCount := len(releases)
	return float64(releasesCount) / frequency
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
