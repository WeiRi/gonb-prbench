# bug.Dockerfile for tidb-38374 — worktree at BUG state, use BUG-API test
FROM gonb-tidb-38374-base:latest

WORKDIR /work/upstream/sessionctx/variable
# Delete all _test.go to avoid old conflicts; plant our BUG-API test
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null

COPY verified_test_bug.go ./tidb_38374_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
