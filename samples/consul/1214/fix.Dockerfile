FROM golang:1.16
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod ./
COPY lock.go ./
COPY lock_fixed.go ./
COPY verified_test.go ./
RUN mv -f lock_fixed.go lock.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
