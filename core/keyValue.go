package core

import "time"

// assumption: all values stored as string.
type Item struct {
	val    string
	expiry time.Time
}

type KV map[string]Item

// ToDo: callback channel for evicted values.
type DB struct {
	KV KV
}

func (db DB) set(key string, val string, ttl int32) {
	// convert ttl to sec
	exp := time.Now().Add(time.Second * time.Duration(ttl))
	db.KV[key] = Item{val: val, expiry: exp}
}
func (db DB) get(key string) string {
	return db.KV[key].val
}

func (db DB) del(key string) string {
	tmp := db.get(key)
	delete(db.KV, key)
	return tmp
}

// ToDo: solve race condition - mutex
// ToDo: minHeap for evict.
func (db DB) evict() {
	now := time.Now().UnixMicro()
	for k, v := range db.KV {
		if v.expiry.UnixMicro() < now {
			db.del(k)
		}
	}
}

func (db DB) ScheduledEvict() {
	ticker := time.NewTicker(TickerFrequency * time.Second)
	for {
		// ToDo: Use select to graceful shutdown of scheduledEvict
		<-ticker.C
		db.evict()
	}
}
