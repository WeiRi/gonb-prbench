# bug.Dockerfile for kubernetes-119729 (in-place BUG state)
FROM gonb-kubernetes-119729-base:latest

WORKDIR /work/upstream/pkg/scheduler

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'schedule_one.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY schedule_one.go ./schedule_one.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'TestScheduleOneBindingFailureRace' .
