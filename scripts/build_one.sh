#!/bin/bash
set -e
SAMPLE_DIR="$1"
[ -d "$SAMPLE_DIR" ] || { echo "ERROR: $SAMPLE_DIR not a dir"; exit 1; }
[ -S "${SSH_AUTH_SOCK:-/tmp/ssh-agent-thu.sock}" ] || {
    echo "ERROR: SSH agent socket not found. Run: eval \$(ssh-agent -a /tmp/ssh-agent-thu.sock); ssh-add ~/.ssh/id_ed25519"
    exit 1
}
TAG="gonb-$(basename $(dirname $SAMPLE_DIR))-$(basename $SAMPLE_DIR)"
SSH_AUTH_SOCK="${SSH_AUTH_SOCK:-/tmp/ssh-agent-thu.sock}" DOCKER_BUILDKIT=1 docker build --ssh default -f "$SAMPLE_DIR/bug.Dockerfile" --network=host -t "$TAG" "$SAMPLE_DIR"
echo "Built: $TAG"
