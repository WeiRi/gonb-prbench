# cockroach-88279

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/88279 |
| Bug commit | `4d6693190c7e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/util/admission/elastic_cpu_granter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000014240 by goroutine 12:
  ase/cockroach-88279.(*TokenBucket).Adjust()
      /work/elastic_cpu_granter.go:23 +0x111
  ase/cockroach-88279.(*elasticCPUGranter).tookWithoutPermission()
      /work/elastic_cpu_granter.go:41 +0xea
  ase/cockroach-88279.Test88279Race.func1()
      /work/verified_test.go:19 +0x135
  ase/cockroach-88279.Test88279Race.gowrap1()
      /work/verified_test.go:23 +0x41

Previous write at 0x00c000014240 by goroutine 8:
  ase/cockroach-88279.(*TokenBucket).TryToFulfill()
      /work/elastic_cpu_granter.go:17 +0x186
  ase/cockroach-88279.(*elasticCPUGranter).tryGet()
      /work/elastic_cpu_granter.go:36 +0x137
  ase/cockroach-88279.Test88279Race.func1()
      /work/verified_test.go:17 +0x142
  ase/cockroach-88279.Test88279Race.gowrap1()
      /work/verified_test.go:23 +0x41

Goroutine 12 (running) created at:
  ase/cockroach-88279.Test88279Race()
      /work/verified_test.go:13 +0x144
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/cockroach-88279.Test88279Race()
      /work/verified_test.go:13 +0x144
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-88279-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-88279-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-88279-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-88279-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-88279-bug .
# (then run as above, no --ssh flag)
```
