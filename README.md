# GoNB-PRBench

A reproducible race-condition benchmark mined from real Go concurrency-fix Pull Requests.

## At a Glance

- **266 verified samples** drawn from **31 Go projects**, including the Go standard library (via 2 CVEs)
- Each sample is a self-contained Docker reproducer: `docker build -f bug.Dockerfile .` produces a container whose race detector flags `WARNING: DATA RACE` (or `panic: concurrent map writes` / `WaitGroup is reused` etc.) on the **real upstream code at the bug commit**
- 5 GoBench-compatible categories: `data_race`, `order_violation`, `anonymous_function`, `channel_misuse`, `special_library`
- 4 oracle types: `RACE`, `PANIC`, `ORDER`, `LIB`-equivalent
- Commit-year range: 2014 – 2026 (13 years)

## Two Reproduction Modes (similar to GoBench)

GoNB-PRBench includes samples produced via two complementary reproduction modes, transparently labeled in `gt.csv` and per-sample README:

1. **In-place trigger** (the race detector fires inside upstream code): the synthesized test imports the upstream package via its public API, sets up a builder/registry and exercises the racy code path; the race detector reports stack frames pointing into the upstream files modified by the PR fix. This is the strongest reproduction mode and is preferred when feasible.

2. **Replicated trigger** (race-pattern fidelity, structurally-equivalent mock): the synthesized test contains a minimal mock that reproduces the **same concurrent access pattern on a primitive of the same type as the field patched by the PR fix** (e.g., bare `bool` → `atomic.Bool`; bare `*Provider` → `atomic.Pointer[Provider]`; missing `sync.Mutex` → field race; etc.). The race detector reports frames inside the mock rather than in upstream files, but the **race semantics — pattern, mechanism, and primitive type — match the PR-targeted race**. This mode is used when the PR's racy code path is gated by unexported internals or a builder/registry pattern that is not feasible to invoke through public API.

Both modes preserve the bug→fix correspondence: applying the PR's `fix.diff` to upstream eliminates the in-place race; for replicated samples, the test models the same logical race the PR's fix is designed to prevent. We document this trade-off in `gt.csv`'s `trigger_mode` column (`in_place` / `replicated`) and in each per-sample README.

## Layout

```
gonb-prbench/
├── README.md      # this file
├── LICENSE
├── gt.csv         # all 263 samples, with category / oracle / year / source
└── samples/
    └── <project>/<pr-or-cve>/
        ├── bug.Dockerfile     # clones upstream, checks out bug commit, drops in test, builds race binary
        ├── fix.Dockerfile     # same as bug, then applies fix.diff (race should NOT trigger after fix)
        ├── fix.diff           # original PR diff (or CVE patch)
        ├── verified_test.go   # synthesized test (only when the PR did not bundle one); some samples omit this
        └── README.md          # per-sample: backtrace excerpt, PR link, how to reproduce
```

## Quick Start

```bash
# 1. Pick any sample, e.g. etcd PR 4958
cd samples/etcd/4958

# 2. Build the bug-state container (will trigger race)
docker build -f bug.Dockerfile -t gonb-prbench-etcd-4958-bug .
docker run --rm --memory=2g --cpus=1 gonb-prbench-etcd-4958-bug
# → expect: WARNING: DATA RACE in etcdserver/api/v3rpc/watch.go

# 3. Build the fix-state container (race should not trigger)
docker build -f fix.Dockerfile -t gonb-prbench-etcd-4958-fix .
docker run --rm --memory=2g --cpus=1 gonb-prbench-etcd-4958-fix
# → expect: PASS (no race detected)
```

See each sample's `README.md` for the exact `go test` invocation and the relevant backtrace.

## Project Distribution

| Project | A | Project | A | Project | A |
|---|---:|---|---:|---|---:|
| kubernetes | 56 | nats-server | 12 | tidb | 4 |
| etcd | 37 | consul | 8 | hugo | 4 |
| grpc-go | 23 | nomad | 8 | dns | 4 |
| cockroach | 22 | prometheus | 7 | (others ≤2) | 12 |
| moby | 20 | istio | 6 |
| go-ethereum | 14 | minio | 6 |

Includes 8 CVE-tagged samples from `golang-go` (stdlib), `argo-cd`, `argo-workflows`, `nomad`, `moby`, `OliveTin`, `emp3r0r`.

## Category × Oracle

| Category | RACE | RACE\|PANIC | PANIC | ORDER | total |
|---|---:|---:|---:|---:|---:|
| data_race | 229 | 0 | 1 | 0 | 230 |
| order_violation | 13 | 0 | 0 | 0 | 13 |
| anonymous_function | 8 | 0 | 0 | 0 | 8 |
| channel_misuse | 4 | 2 | 0 | 0 | 6 |
| special_library | 5 | 0 | 0 | 1 | 6 |

## Citation

If you use this benchmark, please cite:

```
[TODO: citation will be added after publication; ]
```

## License

MIT (see LICENSE).

## Status

This is a release-candidate snapshot for review. Final publication accompanies the
ASE 2026 paper. Issue tracker / contribution guide will be added post-acceptance.


## Building (Requires SSH key for github.com)

Each `bug.Dockerfile` clones the upstream project from `git@github.com:<owner>/<repo>.git` via Docker BuildKit's SSH agent forwarding.

### Setup

```bash
eval $(ssh-agent -a /tmp/ssh-agent-gonb.sock)
ssh-add ~/.ssh/id_ed25519
export SSH_AUTH_SOCK=/tmp/ssh-agent-gonb.sock

DOCKER_BUILDKIT=1 docker build --ssh default \
    -f samples/etcd/4958/bug.Dockerfile \
    -t gonb-prbench-etcd-4958-bug \
    samples/etcd/4958/

docker run --rm --network=host --memory=2g --cpus=1 gonb-prbench-etcd-4958-bug
```

If your network has reliable https access to github.com, you can rewrite `git@github.com:` -> `https://github.com/` in the Dockerfile.

## HTTPS Fallback (if SSH blocked)

Some networks (corporate proxies, restrictive ISPs) block outbound SSH on port 22 to `github.com`. If `git@github.com:` fails, switch every Dockerfile to HTTPS in one command:

```bash
# Apply to all 263 samples
find samples -name "bug.Dockerfile" -o -name "fix.Dockerfile" | xargs sed -i \
  -e 's|git@github.com:|https://github.com/|g' \
  -e 's|--mount=type=ssh ||g'
```

Then build without the `--ssh` flag:

```bash
DOCKER_BUILDKIT=1 docker build \
    -f samples/etcd/4958/bug.Dockerfile \
    -t gonb-etcd-4958-bug \
    samples/etcd/4958/
```

For per-sample fallback (single Dockerfile), the same `sed` works inside an individual sample directory.

### Behind a proxy

If both SSH (22) and HTTPS (443) to `github.com` are unstable but you have an HTTP/SOCKS proxy:

```bash
export HTTPS_PROXY=http://your-proxy:port
export HTTP_PROXY=http://your-proxy:port
DOCKER_BUILDKIT=1 docker build --network=host \
    --build-arg HTTPS_PROXY \
    --build-arg HTTP_PROXY \
    -f bug.Dockerfile -t ... .
```

### Local mirror (very large repos)

If you regularly rebuild many samples from the same project (e.g. `kubernetes/kubernetes` is a 5GB+ repo and per-build clones are slow), pre-clone once to a local bare mirror and use BuildKit `--build-context`:

```bash
# Pre-clone (one-time, ~5-15 min for k8s)
mkdir -p $HOME/git-mirror
git clone --bare --filter=blob:none git@github.com:kubernetes/kubernetes.git $HOME/git-mirror/kubernetes.git

# Then build with the mirror exposed as a build context
DOCKER_BUILDKIT=1 docker build \
    --ssh default \
    --build-context git-mirror=$HOME/git-mirror \
    -f bug.Dockerfile -t gonb-kubernetes-97193-bug .
```

(The Dockerfile's clone step uses `--mount=from=git-mirror,target=/mirror git clone --reference-if-able /mirror/kubernetes.git --dissociate ...`, so an absent mirror is silently bypassed.)
