# serving-10729

| Field | Value |
|---|---|
| Project | serving |
| Reference | https://github.com/knative/serving/pull/10729 |
| Bug commit | `c9fc401a8a48` |
| Category | data_race |
| Oracle | PANIC |
| Primary diff file | `pkg/activator/handler/concurrency_reporter.go` |


## Race report excerpt

```
=== RUN   TestConcurrencyReporterRace_10729
==================
WARNING: DATA RACE
Read at 0x00c0004e2038 by goroutine 75:
  knative.dev/serving/pkg/activator/handler.(*ConcurrencyReporter).computeReport()
      /work/upstream/pkg/activator/handler/concurrency_reporter.go:153 +0x339
  knative.dev/serving/pkg/activator/handler.(*ConcurrencyReporter).report()
      /work/upstream/pkg/activator/handler/concurrency_reporter.go:133 +0x93
  knative.dev/serving/pkg/activator/handler.TestConcurrencyReporterRace_10729.func2()
      /work/upstream/pkg/activator/handler/10729_handcrafted_race_test.go:51 +0xed

Previous write at 0x00c0004e2038 by goroutine 74:
  knative.dev/serving/pkg/activator/handler.(*ConcurrencyReporter).computeReport()
      /work/upstream/pkg/activator/handler/concurrency_reporter.go:154 +0x35d
==================
```

`ConcurrencyReporter.report()` deletes a stat with zero AverageConcurrency.
A concurrent `ReqIn` that grabbed the stat just before the deletion is silently
dropped; subsequent `ReqOut` re-creates the stat with NEGATIVE concurrency.
PR avoids the deletion when a request raced the reporter.

## Reproduce

```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-serving-10729-bug .
docker run --rm --memory=4g --cpus=2 gonb-serving-10729-bug \
    go test -race -count=5 -timeout=120s -run 'TestConcurrencyReporterRace_10729' .

DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-serving-10729-fix .
docker run --rm --memory=4g --cpus=2 gonb-serving-10729-fix \
    go test -race -count=5 -timeout=120s -run 'TestConcurrencyReporterRace_10729' .
```
