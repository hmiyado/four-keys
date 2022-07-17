# four-keys

measure four keys metrics

$$
DeploymentFrequency = (NumOfSuccessfulReleases) / (NumOfDays)
\\
LeadTimeForChanges = mean( (ReleaseDateTime) - (DateTimeOfFirstCommitAfterPreviousRelease) )
\\
TimeToRestoreService = average( (RestoredReleaseDateTime) - (FailureReleaseDateTime) )
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
    "timeToRestoreServices": "00:12:34.567", # future works
    "changeFailureRate": 0.50 # future works
}
$ four-keys releases
{
    "option": {
        "since": "2022-01-01",
        "until": "2022-01-31"
    },
    "releases": [
        {
            "tag": "v1.0.0",
            "date": "2022-01-05 00:00",
            "leadTimeForChanges": {
                "value": "224.73468749999998",
                "unit": "day"
            },
            "result": { # future works
                "type": "success"
            }
        },
        {
            "tag": "v1.0.1",
            "date": "2022-01-10 00:00",
            "leadTimeForChanges": {
                "value": "66.9150462962963",
                "unit": "day"
            },
            "result": { # future works
                "type": "failure"
            }
        },
        {
            "tag": "v1.0.2",
            "date": "2022-01-15 00:00",
            "leadTimeForChanges": {
                "value": "130.77916666666667",
                "unit": "day"
            },
            "result": { # future works
                "type": "success",
                "timeToRestoreService": "120:00:00.000"
            }
        }
    ]
}
```
