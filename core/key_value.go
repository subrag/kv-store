package core

import (
	"log"
	"os"
	"os/signal"
	"time"
)

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

func (db DB) Set(key string, val string, ttl int32) {
	// convert ttl to sec
	exp := time.Now().Add(time.Second * time.Duration(ttl))
	db.KV[key] = Item{val: val, expiry: exp}
}
func (db DB) Get(key string) string {
	return db.KV[key].val
}

func (db DB) Del(key string) string {
	tmp := db.Get(key)
	delete(db.KV, key)
	return tmp
}

// ToDo: solve race condition - mutex
// ToDo: minHeap for evict.
func (db DB) Evict() {
	now := time.Now().UnixMicro()
	for k, v := range db.KV {
		if v.expiry.UnixMicro() < now {
			db.Del(k)
		}
	}
}

func (db DB) ScheduledEvict() {
	ticker := time.NewTicker(TickerFrequency * time.Second)
	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt)
	for {
		select {
		case <-ticker.C:
			db.Evict()
		case <-shutDown:
			log.Println("closing the scheduled evict on os interrupt")
			os.Exit(1)
		}
	}
}
