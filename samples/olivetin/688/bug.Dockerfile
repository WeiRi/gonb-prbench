FROM gonb-olivetin-688-base:latest
WORKDIR /work/upstream/service/internal/entities
RUN find . -maxdepth 1 -name '*_test.go' -delete 2>/dev/null || true
COPY verified_test.go ./olivetin_688_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=180s -run TestRace_688 ."]
