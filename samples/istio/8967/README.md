# istio-8967

| Field | Value |
|---|---|
| Project | istio |
| Reference | https://github.com/istio/istio/pull/8967 |
| Bug commit | `bbee2cec0972` |
| Category | order_violation |
| Oracle | RACE |
| Primary diff file | `pkg/file/fssource.go` |

## Two modes

| Mode | Dockerfile | Test | Description |
|---|---|---|---|
| **in-place** (default) | `bug.Dockerfile` | `verified_test_inplace.go` | Test runs inside the upstream package via SSH-cloned source at the bug commit. Race detector frames hit `fssource.go` — the PR's actual diff target. |
| **verify** (replicated) | `bug.Dockerfile.verify`* | `verified_test.go` | Mock-based stress test in isolated `/work/pr2t-test`, package `main`. Same race semantics. |

\* `bug.Dockerfile.verify` is the legacy SSH-agent build. Only present if the verify mode was preserved during migration.

## Build (in-place)

```bash
docker build --secret id=ssh_key,src=$HOME/.ssh/id_ed25519 \
  -f bug.Dockerfile -t gonb-istio-8967-bug .
docker run --rm --memory=2g --cpus=2 gonb-istio-8967-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=180s -run TestRace"
# Expected: WARNING: DATA RACE + FAIL
```

## Build (fix verification)

```bash
docker build --secret id=ssh_key,src=$HOME/.ssh/id_ed25519 \
  -f fix.Dockerfile -t gonb-istio-8967-fix .
docker run --rm --memory=2g --cpus=2 gonb-istio-8967-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=180s -run TestRace"
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug_inplace.txt` for the full in-place trace and `race_report_bug.txt` for the replicated trace.

## HTTPS fallback

If SSH is unavailable, replace `git@github.com:` with `https://github.com/` in `bug.Dockerfile` (and remove the `--secret` build flag).
