FROM gonb-etcd-1316-bug
CMD ["sh","-c","cd /work/upstream/pkg && go test -race -vet=off -count=10 -timeout=180s -run TestRace ."]
