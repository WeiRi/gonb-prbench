FROM golang:1.22
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod ./
COPY agent.go ./
COPY agent_fixed.go ./
COPY verified_test.go ./
RUN mv -f agent_fixed.go agent.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
