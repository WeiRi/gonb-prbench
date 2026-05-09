# kubernetes-100638

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/100638 |
| Bug commit | `6572fe4d9017` |
| Category | data_race |
| Oracle | PANIC |
| Primary diff file | `queueset.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000383d0 by goroutine 96:
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.(*queueSet).setConfiguration()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/queueset.go:207 +0x24b
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.(*queueSetCompleter).Complete()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/queueset.go:160 +0x4e8
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.TestRace_100638.func2()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/race_100638_test.go:71 +0x276
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.TestRace_100638.func5()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/race_100638_test.go:73 +0x41

Previous read at 0x00c0000383d0 by goroutine 118:
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.(*queueSet).StartRequest.func1()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/queueset.go:300 +0x189

Goroutine 96 (running) created at:
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.TestRace_100638()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/race_100638_test.go:57 +0x4f3
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 118 (running) created at:
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.(*queueSet).StartRequest()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/queueset.go:290 +0xd9c
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.TestRace_100638.func1()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/race_100638_test.go:43 +0x188
  k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset.TestRace_100638.func4()
      /work/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset/race_100638_test.go:51 +0x41
==================
E0507 14:53:14.949700    3876 runtime.go:76] Observed a panic: sync: negative WaitGroup counter
goroutine 31 [running]:
k8s.io/apimachinery/pkg/util/runtime.logPanic({0xe55e60?, 0x1079d00?})
	/work/staging/src/k8s.io/apimachinery/pkg/util/runtime/runtime.go:74 +0xdd
```

(Full trace in `race_report_bug.txt`.)

## How to reproduce

### 1. SSH agent setup (one-time)
```bash
eval $(ssh-agent -a /tmp/ssh-agent-gonb.sock)
ssh-add ~/.ssh/id_ed25519
export SSH_AUTH_SOCK=/tmp/ssh-agent-gonb.sock
```

### 2. Build bug image
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-100638-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-100638-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-100638-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-100638-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-100638-bug .
# (then run as above, no --ssh flag)
```
