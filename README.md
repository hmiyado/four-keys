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
    "repository": {
        "name": "some-repo"
    },
    "option": {
        "startDate": "2022-01-01",
        "endDate": "2022-01-31"
    },
    "deploymentFrequency": 0.5,
    "leadTimeForChanges": "12:34:56.789",
    "timeToRestoreServices": "00:12:34.567",
    "changeFailureRate": 0.50
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
            "leadTimeForChanges": "11:22:33.444",
            "result": {
                "type": "success"
            }
        },
        {
            "tag": "v1.0.1",
            "date": "2022-01-10 00:00",
            "leadTimeForChanges": "12:34:56.000",
            "result": {
                "type": "failure"
            }
        },
        {
            "tag": "v1.0.2",
            "date": "2022-01-15 00:00",
            "leadTimeForChanges": "12:34:56.000",
            "result": {
                "type": "success",
                "timeToRestoreService": "120:00:00.000"
            }
        }
    ]
}
```
