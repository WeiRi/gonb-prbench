package clientv3

import "context"

// Local metadata stub mimicking google.golang.org/grpc/metadata.MD
type MDLocal map[string][]string

func PairsLocal(kvs ...string) MDLocal {
	md := MDLocal{}
	for i := 0; i < len(kvs); i += 2 {
		md[kvs[i]] = append(md[kvs[i]], kvs[i+1])
	}
	return md
}

func (m MDLocal) Set(k, v string) {
	m[k] = []string{v}
}

func (m MDLocal) Get(k string) []string {
	return m[k]
}

type ctxKeyMD struct{}

func NewOutgoingContextLocal(ctx context.Context, md MDLocal) context.Context {
	return context.WithValue(ctx, ctxKeyMD{}, md)
}

func FromOutgoingContextLocal(ctx context.Context) (MDLocal, bool) {
	v, ok := ctx.Value(ctxKeyMD{}).(MDLocal)
	return v, ok
}
