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
$ four-keys releases --repository https://github.com/go-git/go-git --since 2015-12-20 --until 2016-01-12 | jq
{
  "option": {
    "since": "2015-12-20T00:00:00Z",
    "until": "2016-01-12T23:59:59Z",
    "ignorePattern": null
  },
  "releases": [
    {
      "tag": "v2.1.2",
      "date": "2016-01-11T12:09:15+01:00",
      "leadTimeForChanges": {
        "value": 0.017638888888888888,
        "unit": "day"
      },
      "result": {
        "isSuccess": true,
        "timeToRestore": {
          "value": 2.7969791666666666,
          "unit": "day"
        }
      }
    },
    {
      "tag": "v2.1.1",
      "date": "2016-01-08T17:01:36+01:00",
      "leadTimeForChanges": {
        "value": 0.00863425925925926,
        "unit": "day"
      },
      "result": {
        "isSuccess": false,
        "timeToRestore": null
      }
    },
    {
      "tag": "v2.1.0",
      "date": "2015-12-23T09:48:11+01:00",
      "leadTimeForChanges": {
        "value": 6.587986111111111,
        "unit": "day"
      },
      "result": {
        "isSuccess": true,
        "timeToRestore": null
      }
    }
  ]
}
```
