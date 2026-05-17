FROM golang:1.20
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod bucket_fixed.go bucket.go verified_test.go ./
RUN rm -f bucket.go && mv bucket_fixed.go bucket.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
