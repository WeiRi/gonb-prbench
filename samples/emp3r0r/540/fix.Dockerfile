FROM golang:1.22
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod handler_fixed.go handler.go verified_test.go ./
RUN rm -f handler.go && mv handler_fixed.go handler.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
