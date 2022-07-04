# four-keys

measure four keys metrics

$$
DeploymentFrequency = (NumOfSuccessfulReleases) / (NumOfDays)
\\
LeadTimeForChanges = average( max( (ReleaseDateTime) - (CommitDateTime) ) )
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
        "startDate": "2022-01-01",
        "endDate": "2022-01-31"
    },
    "deploymentFrequency": 0.5,
    "leadTimeForChanges": "12:34:56.789", # future works
    "timeToRestoreServices": "00:12:34.567", # future works
    "changeFailureRate": 0.50 # future works
}
$ four-keys releases
{
    "option": {
        "startDate": "2022-01-01",
        "endDate": "2022-01-31"
    },
    "releases": [
        {
            "tag": "v1.0.0",
            "date": "2022-01-05 00:00",
            "leadTimeForChanges": "11:22:33.444", # future works
            "result": { # future works
                "type": "success"
            }
        },
        {
            "tag": "v1.0.1",
            "date": "2022-01-10 00:00",
            "leadTimeForChanges": "12:34:56.000", # future works
            "result": { # future works
                "type": "failure"
            }
        },
        {
            "tag": "v1.0.2",
            "date": "2022-01-15 00:00",
            "leadTimeForChanges": "12:34:56.000", # future works
            "result": { # future works
                "type": "success",
                "timeToRestoreService": "120:00:00.000"
            }
        }
    ]
}
```
