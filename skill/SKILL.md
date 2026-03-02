---
name: intervals-cli
description: Query and update Intervals.icu training data — activities, events, wellness, and sport settings.
version: 0.0.1
metadata:
  openclaw:
    emoji: "🚴"
    homepage: https://github.com/runningcode/intervals-cli
    primaryEnv: INTERVALS_API_KEY
    requires:
      env:
        - INTERVALS_ATHLETE_ID
        - INTERVALS_API_KEY
      bins:
        - intervals-cli
    install:
      - kind: brew
        tap: runningcode/tap
        formula: intervals-cli
        bins: [intervals-cli]
---

# intervals-cli

CLI for the [Intervals.icu](https://intervals.icu) API. Use it to read and
write training data: activities, calendar events, wellness records, and sport
settings.

## Authentication

Set these environment variables before calling any command (except `config`):

```
INTERVALS_ATHLETE_ID=<athlete_id>
INTERVALS_API_KEY=<api_key>
```

Or run once to persist to `~/.intervals-cli/config.json`:

```bash
intervals-cli config --athlete-id ID --api-key KEY
```

## Output

All commands write JSON to stdout:

```json
{"data": <object|array>, "metadata": {"count": N, "tool": "intervals-cli", "version": "..."}}
```

Errors go to stderr as `{"error": "...", "code": "ERR_..."}` with a non-zero exit code.

## Commands

### Activities

```bash
# List activities in a date range
intervals-cli activities --oldest 2025-01-01 --newest 2025-01-31

# Get a single activity
intervals-cli activity --id ACTIVITY_ID

# Get time-series streams (power, HR, cadence, etc.)
intervals-cli activity-streams --id ACTIVITY_ID

# Create a manual activity
intervals-cli create-activity --name "Morning Run" --type Run --start 2025-01-15T07:00:00Z --distance 10000 --moving-time 3600

# Update an activity
intervals-cli update-activity --id ACTIVITY_ID --name "New Name" --description "Updated"
```

### Events (Calendar)

```bash
# List events in a date range
intervals-cli events --oldest 2025-01-01 --newest 2025-01-31

# Filter by category
intervals-cli events --oldest 2025-01-01 --newest 2025-01-31 --category WORKOUT

# Create an event
intervals-cli create-event --name "Long Ride" --category WORKOUT --start 2025-01-20T09:00:00Z

# Update an event
intervals-cli update-event --id EVENT_ID --name "Updated Name"

# Download a workout file
intervals-cli download-event --id EVENT_ID --dl-format zwo --output workout.zwo
```

Supported download formats: `zwo`, `mrc`, `erg`, `fit`.

### Wellness

```bash
# List wellness records
intervals-cli wellness --oldest 2025-01-01 --newest 2025-01-31

# Update a wellness record (only provided flags are sent)
intervals-cli update-wellness --date 2025-01-15 --sleep-secs 28800 --resting-hr 52 --weight 70.5
```

Wellness scale flags (`--fatigue`, `--soreness`, `--stress`, `--mood`) accept values 1–5.

### Sport Settings

```bash
# Get settings for a sport
intervals-cli sport-settings --sport Ride

# Update settings (only provided flags are sent)
intervals-cli update-sport-settings --sport Ride --ftp 280 --lthr 162 --max-hr 185
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (API failure, I/O, config) |
| 2 | Invalid usage — missing flag or bad argument |
| 3 | Data error — response could not be parsed |
