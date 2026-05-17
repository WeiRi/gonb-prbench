FROM golang:1.22
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod ./
COPY agent.go ./
COPY verified_test.go ./
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
