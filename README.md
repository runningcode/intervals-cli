# intervals-cli

A command-line tool for the [Intervals.icu](https://intervals.icu) API. Manage
athlete training data — activities, calendar events, wellness records, and sport
settings — from your terminal or scripts.

## Installation

### Homebrew

```bash
brew tap runningcode/tap
brew install intervals-cli
```

### Build from source

Requires Go 1.22+.

```bash
git clone https://github.com/runningcode/intervals-cli
cd intervals-cli
go build -o intervals-cli .
```

## Configuration

You need an Intervals.icu athlete ID and API key. Find them under
**Settings → API** in the Intervals.icu web app.

```bash
intervals-cli config --athlete-id YOUR_ATHLETE_ID --api-key YOUR_API_KEY
```

This writes credentials to `~/.intervals-cli/config.json`. You can also use
environment variables, which take precedence over the config file:

```bash
export INTERVALS_ATHLETE_ID=your_id
export INTERVALS_API_KEY=your_key
```

## Usage

All commands output JSON by default. Pass `--format text` for a
human-readable table.

### Activities

```bash
# List activities for January 2024
intervals-cli activities --oldest 2024-01-01 --newest 2024-01-31

# Get a single activity
intervals-cli activity --id 12345678

# Get time-series streams for an activity
intervals-cli activity-streams --id 12345678

# Create a manual activity
intervals-cli create-activity --name "Morning Run" --type Run --start 2024-01-15T06:00:00

# Update an activity
intervals-cli update-activity --id 12345678 --name "Easy Run" --description "Recovery pace"
```

### Events

```bash
# List events for February 2024
intervals-cli events --oldest 2024-02-01 --newest 2024-02-29

# List only workouts
intervals-cli events --oldest 2024-02-01 --newest 2024-02-29 --category WORKOUT

# Create a workout event
intervals-cli create-event --name "VO2Max Intervals" --category WORKOUT --start 2024-02-15T18:00:00

# Update an event
intervals-cli update-event --id 98765 --name "Rescheduled VO2Max"

# Download a workout file
intervals-cli download-event --id 98765 --dl-format zwo --output workout.zwo
```

### Wellness

```bash
# List wellness records
intervals-cli wellness --oldest 2024-01-01 --newest 2024-01-31

# Log sleep and resting heart rate for a day
intervals-cli update-wellness --date 2024-01-15 --sleep-secs 28800 --resting-hr 52

# Log subjective metrics (1–5 scale)
intervals-cli update-wellness --date 2024-01-15 --fatigue 2 --mood 4
```

### Sport Settings

```bash
# Get cycling settings
intervals-cli sport-settings --sport Ride

# Update FTP
intervals-cli update-sport-settings --sport Ride --ftp 310
```

### Config

```bash
intervals-cli config --athlete-id YOUR_ID --api-key YOUR_KEY
```

## Output Format

Every command wraps its output in a consistent envelope:

```json
{
  "data": [ ... ],
  "metadata": {
    "count": 3,
    "tool": "intervals-cli",
    "version": "1.0.0"
  }
}
```

Errors go to stderr:

```json
{"error": "athlete ID and API key required", "code": "ERR_AUTH"}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid usage (missing required flag, bad arguments) |
| 3 | Data error (parse failure) |

## Publishing a Release

See [PUBLISHING.md](PUBLISHING.md) for one-time setup and release instructions.
