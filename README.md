# logslice

Fast log segmentation tool that splits large log files by time window or pattern.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

## Usage

Split a log file into 1-hour time windows:

```bash
logslice -input app.log -window 1h -output ./slices/
```

Split by pattern match:

```bash
logslice -input app.log -pattern "ERROR|FATAL" -output errors.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-input` | Path to the source log file | required |
| `-output` | Output file or directory | `./out` |
| `-window` | Time window size (e.g. `15m`, `1h`, `24h`) | `1h` |
| `-pattern` | Regex pattern to filter log lines | none |
| `-format` | Timestamp format in log lines | auto-detect |

### Example

```bash
# Extract the last 30 minutes of logs matching WARNING or above
logslice -input /var/log/app.log -window 30m -pattern "WARNING|ERROR|FATAL" -output ./recent_errors.log
```

## Requirements

- Go 1.21+

## License

MIT © 2024 yourusername