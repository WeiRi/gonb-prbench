FROM golang:1.22
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=1
WORKDIR /work
COPY go.mod event_broadcaster_fixed.go event_broadcaster.go verified_test.go ./
RUN rm -f event_broadcaster.go && mv event_broadcaster_fixed.go event_broadcaster.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
