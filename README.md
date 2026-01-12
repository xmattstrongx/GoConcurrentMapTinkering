# GoConcurrentMapTinkering

A tiny playground app that demonstrates usage of
`github.com/xmattstrongx/go_concurrent_map` with expiration, purge loops, and
concurrent access.

## Run the demo

```sh
go run ./cmd/demo
```

## What the demo covers

- Build a concurrent map with default expiration and purge interval.
- Set/get basic entries.
- Per-entry expiration and non-expiring entries.
- Concurrent readers and writers.

## Notes

If your Go build cache is not writable, set `GOCACHE` to a writable path:

```sh
GOCACHE=/path/to/writable/cache go run ./cmd/demo
```

## Benchmarks and tests

Run the expiration tests and benchmarks:

```sh
TMPDIR=/path/to/writable/tmp GOCACHE=/path/to/writable/cache /path/to/go/bin/go test -run TestExpiration -bench .
```
