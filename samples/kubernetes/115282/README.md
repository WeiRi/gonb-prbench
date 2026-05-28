# kubernetes-115282

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/115282 |
| Bug commit | `ac6d67d27c6` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apiserver/pkg/server/config.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000190000 by goroutine 26:
  ase/kubernetes-115282.(*warningRecorder).AddWarning()
      /work/config.go:12 +0x111
  ase/kubernetes-115282.(*Chain).ServeRequest.func1()
      /work/config.go:31 +0x51

Previous read at 0x00c000190000 by goroutine 10:
  ase/kubernetes-115282.(*warningRecorder).Snapshot()
      /work/config.go:16 +0x139
  ase/kubernetes-115282.(*Chain).ServeRequest()
      /work/config.go:36 +0x11a
  ase/kubernetes-115282.TestWarningWithRequestTimeout_115282.func1()
      /work/verified_test.go:19 +0x104
==================
```

`DefaultBuildHandlerChain` wraps `WithWarningRecorder` OUTSIDE the timeout
handler. After timeout fires, the inner handler keeps calling
`warningRecorder.AddWarning()` which writes to the recorder's shared slice,
while the timeout filter's response writer reads concurrently → race.
Fix wraps `WithWarningRecorder` INSIDE the timeout wrapper.

## Reproduce

```bash
# BUG state (expect FAIL: DATA RACE)
docker build -f bug.Dockerfile -t gonb-kubernetes-115282-bug .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-115282-bug

# FIX state (expect ok)
docker build -f fix.Dockerfile -t gonb-kubernetes-115282-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-115282-fix
```
