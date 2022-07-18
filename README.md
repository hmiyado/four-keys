# four-keys

measure four keys metrics

$$
DeploymentFrequency = (NumOfSuccessfulReleases) / (NumOfDays)
\\
LeadTimeForChanges = mean( (ReleaseDateTime) - (DateTimeOfFirstCommitAfterPreviousRelease) )
\\
TimeToRestore = average( (RestoredReleaseDateTime) - (FailureReleaseDateTime) )
\\
ChangeFailureRate = (NumOfFailureRelease) / (NumOfReleases)
$$

## Example

```sh
$ cd some-repo
$ four-keys
{
    "option": {
        "since": "2022-01-01",
        "until": "2022-01-31"
    },
    "deploymentFrequency": 0.5,
    "leadTimeForChanges": {
        "value": "98.84710648148149",
        "unit": "day"
    },
    "timeToRestore": {
        "value": 2.7969791666666666,
        "unit": "day"
    },
    "changeFailureRate": 0.50
}
$ four-keys releases
{
    "option": {
        "since": "2022-01-01",
        "until": "2022-01-31"
    },
    "releases": [
        {
            "tag": "v1.0.2",
            "date": "2022-01-15 00:00",
            "leadTimeForChanges": {
                "value": "130.77916666666667",
                "unit": "day"
            },
            "result": {
                "isSuccess": true,
                "timeToRestore": "120:00:00.000" # future works
            }
        },
        {
            "tag": "v1.0.1",
            "date": "2022-01-10 00:00",
            "leadTimeForChanges": {
                "value": "66.9150462962963",
                "unit": "day"
            },
            "result": {
                "isSucecss": false
            }
        },

        {
            "tag": "v1.0.0",
            "date": "2022-01-05 00:00",
            "leadTimeForChanges": {
                "value": "224.73468749999998",
                "unit": "day"
            },
            "result": {
                "isSuccess": true
            }
        }
    ]
}
```
