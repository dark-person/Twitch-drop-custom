# Twitch-drop-custom

A script to open specified link at specified time, customable by user. 

**Designed for Windows Only.**

## Build
```
go install
go build
twitch-drop-custom.exe
```

## Requirement before Run
1. **The prorgam must be running background to fullfill its purpose**, which means the computer should not be power off.
2. Your default browser should have the script/extentsion to auto claim the twitch drop.
3. **Mute the twitch stream** if you don't want to distrub others.

### Config Instructions
Example:

```
last_updated: 202212090115

task:
  - link : "https://youtu.be/dQw4w9WgXcQ"
    start_time : "2022-12-10 22:16:30"
    duration: 5
  - link : "https://youtu.be/dQw4w9WgXcQ"
    start_time : "2022-12-12 22:16:30"
    duration: 20
```

***Please ensure your YAML file format is correct before run/publish.***

`lastupdated` is for version checking, which is required.

`task` should include the detail of when to open the link.
The detail is include:
| Field Name   | Description                                         | Format | Require      |
| ------------ | --------------------------------------------------- | ------ | ------------ |
| `link`       | The link to open                                    | string | **Required** |
| `start_time` | The start time of the stream                        | string | **Required** |
| `duration`   | the time for that task to run, *Unit is in minutes* | number | **Required** |

## Reminder

1. Task can be at the same time, ***but not suggested***, as the browser need sometime to load the webpage, especially for twitch stream.
2. You can run at most five task synchronously.
3. Task with start time before current system time will be skipped
4. The webpage will not closed automatically closed. Please ensure you don't have too many browser tabs before running.

