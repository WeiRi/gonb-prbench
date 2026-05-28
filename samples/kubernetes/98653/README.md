# kubernetes-98653

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/98653 |
| Bug commit | `2783f2f76ec5` |
| Category | data_race |
| Oracle | PANIC |
| Primary diff file | `staging/src/k8s.io/apimachinery/pkg/watch/streamwatcher.go` |


## Race report excerpt

The bug is a goroutine leak — `StreamWatcher.receive()` blocks forever on send
when ResultChan isn't drained and Stop() ran after the `stopping()` check.
Race detector doesn't catch this directly; oracle is `runtime.NumGoroutine()` count:

```
=== RUN   TestStreamWatcherRace_98653
goroutine dump (PANIC oracle 98653, leak):
goroutine 5 [running]:
k8s.io/apimachinery/pkg/watch.TestStreamWatcherRace_98653(...)
    /work/upstream/staging/src/k8s.io/apimachinery/pkg/watch/98653_handcrafted_race_test.go
...
--- FAIL: TestStreamWatcherRace_98653 (90 leaked goroutines per iter)
```

PR adds a `done` channel signaled by `Stop()` so the pending send selects out.

## Reproduce

```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-98653-bug .
docker run --rm --memory=4g --cpus=2 gonb-kubernetes-98653-bug \
    go test -race -count=3 -timeout=120s -run 'TestStreamWatcherRace_98653' .

DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-98653-fix .
docker run --rm --memory=4g --cpus=2 gonb-kubernetes-98653-fix \
    go test -race -count=3 -timeout=120s -run 'TestStreamWatcherRace_98653' .
```
