package core

import "time"

// assumption: all values stored as string.
type Item struct {
	val    string
	expiry time.Time
}

type KV map[string]Item

type DB struct {
	KV KV
}

func (db DB) set(a string, b string, ttl int32) {
	exp := time.Now().Add(time.Second * time.Duration(ttl))
	db.KV[a] = Item{val: b, expiry: exp}
}
func (db DB) get(a string) string {
	return db.KV[a].val
}

func (db DB) evict() {
}

// func (db DB)
// ToDo: solve race condition - mutex
// ToDo: minHeap for evict
