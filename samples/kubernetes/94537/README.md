# kubernetes-94537

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/94537 |
| Bug commit | `a84419f02766` |
| Category | data_race |
| Oracle | PANIC |
| Primary diff file | `staging/src/k8s.io/legacy-cloud-providers/azure/cache/azure_cache.go` |


## Race report excerpt

The following PANIC oracle output is captured when running the bug build (`go test -race`):

```
=== RUN   TestCacheNoConcurrentGet_94537
goroutine dump (PANIC oracle 94537 — getter call count > 1):
goroutine 12 [running]:
k8s.io/legacy-cloud-providers/azure/cache.TestCacheNoConcurrentGet_94537.func1(...)
    /work/upstream/staging/src/k8s.io/legacy-cloud-providers/azure/cache/94537_handcrafted_race_test.go
k8s.io/legacy-cloud-providers/azure/cache.(*TimedCache).getInternal(...)
    /work/upstream/staging/src/k8s.io/legacy-cloud-providers/azure/cache/azure_cache.go
...
--- FAIL: TestCacheNoConcurrentGet_94537
```

`TimedCache.getInternal` does not re-check `Store.GetByKey` inside the lock, so
concurrent goroutines each create their own `AzureCacheEntry{Data: nil}` and
each invokes the user-supplied Getter for the same key. PR adds the
under-lock re-check.

## Reproduce

```bash
# BUG state (expect FAIL: PANIC oracle fires)
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-94537-bug .
docker run --rm --memory=4g --cpus=2 gonb-kubernetes-94537-bug \
    go test -race -count=10 -timeout=120s -run 'TestCacheNoConcurrentGet_94537' .

# FIX state (expect ok)
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-94537-fix .
docker run --rm --memory=4g --cpus=2 gonb-kubernetes-94537-fix \
    go test -race -count=10 -timeout=120s -run 'TestCacheNoConcurrentGet_94537' .
```
