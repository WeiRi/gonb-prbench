# bug.Dockerfile for kubernetes-90348 (in-place BUG state)
FROM gonb-kubernetes-90348-base:latest

WORKDIR /work/upstream/staging/src/k8s.io/legacy-cloud-providers/vsphere

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'vsphere.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY vsphere.go ./vsphere.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'TestRace_90348_LoopVarCapture' .
