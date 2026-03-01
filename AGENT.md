# intervals-cli — Agent Reference

Machine-readable reference for AI agents calling this CLI.

## Authentication

Credentials are required for all commands except `config` and `--version`.

**Preferred (env vars):**
```
INTERVALS_ATHLETE_ID=<athlete_id>
INTERVALS_API_KEY=<api_key>
```

**Alternatively:** run `intervals-cli config --athlete-id ID --api-key KEY`
once to persist credentials to `~/.intervals-cli/config.json`.

Env vars override the config file when both are present.

## Output Schema

All successful commands write to stdout:

```json
{
  "data": <object or array>,
  "metadata": {
    "count": <integer>,
    "tool": "intervals-cli",
    "version": "<semver>"
  }
}
```

`data` is an array for list commands and an object for single-resource commands.
`metadata.count` is the number of items in `data` (1 for object responses).

## Error Schema

Errors write to stderr and the process exits with a non-zero code:

```json
{"error": "<human-readable message>", "code": "<ERR_CODE>", "details": <optional>}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (API failure, I/O error, config error) |
| 2 | Invalid usage — missing required flag or bad argument |
| 3 | Data error — response could not be parsed |

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `json` | Output format: `json` or `text` |
| `--quiet` | false | Suppress non-essential output |
| `--version` | — | Print version JSON and exit |

## Commands

### `config`

Write credentials to `~/.intervals-cli/config.json`.

```
intervals-cli config --athlete-id ID --api-key KEY
```

| Flag | Required | Description |
|------|----------|-------------|
| `--athlete-id` | yes | Intervals.icu athlete ID |
| `--api-key` | yes | Intervals.icu API key |

---

### `activities`

List activities for a date range.

```
intervals-cli activities --oldest YYYY-MM-DD --newest YYYY-MM-DD
```

| Flag | Required | Description |
|------|----------|-------------|
| `--oldest` | yes | Start date (inclusive) |
| `--newest` | yes | End date (inclusive) |

`data` is an array of activity objects.

---

### `activity`

Get a single activity by ID.

```
intervals-cli activity --id ACTIVITY_ID
```

| Flag | Required | Description |
|------|----------|-------------|
| `--id` | yes | Activity ID |

`data` is an activity object.

---

### `activity-streams`

Get time-series streams for an activity (power, HR, cadence, etc.).

```
intervals-cli activity-streams --id ACTIVITY_ID
```

| Flag | Required | Description |
|------|----------|-------------|
| `--id` | yes | Activity ID |

`data` is the streams object as returned by the API.

---

### `create-activity`

Create a manual activity.

```
intervals-cli create-activity --name NAME --type TYPE --start ISO8601 [--distance M] [--moving-time S] [--description TEXT]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | yes | Activity name |
| `--type` | yes | Activity type (Run, Ride, Swim, etc.) |
| `--start` | yes | Start time in ISO 8601 format |
| `--distance` | no | Distance in meters |
| `--moving-time` | no | Moving time in seconds |
| `--description` | no | Activity description |

`data` is the created activity object.

---

### `update-activity`

Update fields on an existing activity. Only explicitly provided flags are sent.

```
intervals-cli update-activity --id ACTIVITY_ID [--name NAME] [--description TEXT] [--type TYPE]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--id` | yes | Activity ID |
| `--name` | no | New name |
| `--description` | no | New description |
| `--type` | no | New type |

`data` is the updated activity object.

---

### `events`

List calendar events for a date range.

```
intervals-cli events --oldest YYYY-MM-DD --newest YYYY-MM-DD [--category CATEGORY]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--oldest` | yes | Start date (inclusive) |
| `--newest` | yes | End date (inclusive) |
| `--category` | no | Filter by category (WORKOUT, NOTE, etc.) |

`data` is an array of event objects.

---

### `create-event`

Create a calendar event.

```
intervals-cli create-event --name NAME --category CATEGORY --start ISO8601 [--description TEXT] [--color COLOR]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | yes | Event name |
| `--category` | yes | Event category (WORKOUT, NOTE, etc.) |
| `--start` | yes | Start date/time in ISO 8601 format |
| `--description` | no | Event description |
| `--color` | no | Event color |

`data` is the created event object.

---

### `update-event`

Update fields on an existing event. Only explicitly provided flags are sent.

```
intervals-cli update-event --id EVENT_ID [--name NAME] [--description TEXT] [--category CATEGORY] [--color COLOR]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--id` | yes | Event ID |
| `--name` | no | New name |
| `--description` | no | New description |
| `--category` | no | New category |
| `--color` | no | New color |

`data` is the updated event object.

---

### `download-event`

Download a workout file to disk.

```
intervals-cli download-event --id EVENT_ID --dl-format FORMAT --output PATH
```

| Flag | Required | Description |
|------|----------|-------------|
| `--id` | yes | Event ID |
| `--dl-format` | yes | File format: `zwo`, `mrc`, `erg`, or `fit` |
| `--output` | yes | Output file path |

`data` on success:
```json
{"path": "/path/to/file.zwo", "bytes": 1234}
```

---

### `wellness`

List wellness records for a date range.

```
intervals-cli wellness --oldest YYYY-MM-DD --newest YYYY-MM-DD
```

| Flag | Required | Description |
|------|----------|-------------|
| `--oldest` | yes | Start date (inclusive) |
| `--newest` | yes | End date (inclusive) |

`data` is an array of wellness objects.

---

### `update-wellness`

Update a wellness record for a specific date. Only explicitly provided flags are sent.

```
intervals-cli update-wellness --date YYYY-MM-DD [--sleep-secs N] [--resting-hr N] [--weight N] [--fatigue N] [--soreness N] [--stress N] [--mood N]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--date` | yes | Date (YYYY-MM-DD) |
| `--sleep-secs` | no | Sleep duration in seconds |
| `--resting-hr` | no | Resting heart rate (bpm) |
| `--weight` | no | Body weight in kg |
| `--fatigue` | no | Fatigue level (1–5) |
| `--soreness` | no | Soreness level (1–5) |
| `--stress` | no | Stress level (1–5) |
| `--mood` | no | Mood level (1–5) |

`data` is the updated wellness object.

---

### `sport-settings`

Get settings for a sport.

```
intervals-cli sport-settings --sport SPORT
```

| Flag | Required | Description |
|------|----------|-------------|
| `--sport` | yes | Sport type (Ride, Run, Swim, etc.) |

`data` is the sport settings object.

---

### `update-sport-settings`

Update settings for a sport. Only explicitly provided flags are sent.

```
intervals-cli update-sport-settings --sport SPORT [--ftp N] [--lthr N] [--max-hr N] [--resting-hr N]
```

| Flag | Required | Description |
|------|----------|-------------|
| `--sport` | yes | Sport type (Ride, Run, Swim, etc.) |
| `--ftp` | no | Functional threshold power (watts) |
| `--lthr` | no | Lactate threshold heart rate (bpm) |
| `--max-hr` | no | Maximum heart rate (bpm) |
| `--resting-hr` | no | Resting heart rate (bpm) |

`data` is the updated sport settings object.
