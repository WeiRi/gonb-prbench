# syntax=docker/dockerfile:1.4
# bug.Dockerfile (in-place) for istio-8214
# Build:
#   docker build --secret id=ssh_key,src=$HOME/.ssh/id_ed25519 -f bug.Dockerfile -t img .

FROM golang:1.10 AS clone
RUN --mount=type=secret,id=ssh_key \
    mkdir -p /root/.ssh && \
    cp /run/secrets/ssh_key /root/.ssh/id_ed25519 && \
    chmod 600 /root/.ssh/id_ed25519 && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts && \
    mkdir -p /go/src/istio.io && \
    git clone --depth=200 git@github.com:istio/istio.git /go/src/istio.io/istio && \
    cd /go/src/istio.io/istio && \
    git fetch --depth=200 origin be6aaeb4ede98b198c2c818b38bd3a5a0018304d && \
    git checkout --detach be6aaeb4ede98b198c2c818b38bd3a5a0018304d~1 && \
    rm -rf /root/.ssh

FROM golang:1.10
ENV GOPATH=/go GO111MODULE=off CGO_ENABLED=1
COPY --from=clone /go/src/istio.io/istio /go/src/istio.io/istio
WORKDIR /go/src/istio.io/istio
RUN find /go/src/istio.io/istio/pkg/cache -name "*_test.go" ! -name "istio-8214_race_test.go" -delete 2>/dev/null || true
COPY verified_test_inplace.go /go/src/istio.io/istio/pkg/cache/istio-8214_race_test.go
WORKDIR /go/src/istio.io/istio/pkg/cache
