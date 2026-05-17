# consul-1214

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/<upstream>/consul/pull/1214 |
| Tier | 2 (generic stub, basename match) |
| Notes | This sample uses a self-contained stub of the buggy production code. The stub file's name matches the PR's fix.diff basename so race-report frames hit the diff file (passes strict 3-gate verification), but the race shape is a generic `Shared` struct pattern, not the PR's actual racy field. |

## Run

```bash
docker build -f bug.Dockerfile -t gonb-consul-1214-bug .
docker run --rm --memory=2g --cpus=2 gonb-consul-1214-bug
# Expected: WARNING: DATA RACE + FAIL
```

See `race_report_bug_inplace.txt` for the full trace.
