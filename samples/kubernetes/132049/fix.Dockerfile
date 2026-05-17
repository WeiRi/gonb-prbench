FROM golang:1.22
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod patch_fixed.go patch.go verified_test.go ./
RUN rm -f patch.go && mv patch_fixed.go patch.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
