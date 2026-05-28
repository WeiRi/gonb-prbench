# kubernetes-98956

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/98956 |
| Bug commit | `d0a433fa4504` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kubelet/kubelet_pods.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0004e4124 by goroutine 428:
  k8s.io/kubernetes/pkg/kubelet.(*podKillerWithChannel).IsPodPendingTerminationByUID()
      /work/upstream/pkg/kubelet/kubelet_pods.go:1187 +0x21f
  k8s.io/kubernetes/pkg/kubelet.TestKillPodFollwedByIsPodPendingTermination_98956.func3()
      /work/upstream/pkg/kubelet/98956_handcrafted_race_test.go:65 +0x2f2

Previous write at 0x00c0004e4124 by goroutine 135:
  k8s.io/kubernetes/pkg/kubelet.(*podKillerWithChannel).KillPod()
      /work/upstream/pkg/kubelet/kubelet_pods.go:1226 +0x21f
==================
```

Before fix, `KillPod` enqueues a podKillingCh request; the marking-as-pending
happens later in `PerformPodKillingWork` (different goroutine). Between
`KillPod` return and the marking, `IsPodPendingTerminationByUID` returns false.
PR's fix: `KillPod` itself takes the lock and marks pending before enqueuing.

## Reproduce

```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-98956-bug .
docker run --rm --memory=4g --cpus=2 gonb-kubernetes-98956-bug \
    go test -race -count=3 -timeout=120s -run 'TestKillPodFollwedByIsPodPendingTermination_98956' .

DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-98956-fix .
docker run --rm --memory=4g --cpus=2 gonb-kubernetes-98956-fix \
    go test -race -count=3 -timeout=120s -run 'TestKillPodFollwedByIsPodPendingTermination_98956' .
```
