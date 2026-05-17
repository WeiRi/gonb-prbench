FROM gonb-nomad-14188-base:latest

WORKDIR /work/upstream/nomad/stream
RUN for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done

COPY verified_test.go ./nomad_14188_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
