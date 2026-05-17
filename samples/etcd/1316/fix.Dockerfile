FROM golang:1.20
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod leader_fixed.go leader.go verified_test_inplace.go  verified_test.go ./
RUN  rm -f leader.go && mv leader_fixed.go leader.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
