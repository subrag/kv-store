package core

type DB struct {
	hmap KV
}

type KV map[string]string

func (kv KV) set(a, b string) {
	kv[a] = b
}

func (kv KV) get(a string) string {
	return kv[a]
}
