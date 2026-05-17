FROM gonb-kubernetes-134390-bug:latest
ENV GOFLAGS= GOWORK=off
WORKDIR /work/pr2t-test
COPY metric_fixed.go ./metric.go
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./kubernetes_134390_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=180s -run TestRace_134390 ."]
