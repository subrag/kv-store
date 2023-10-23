package core_test

import (
	"testing"
	"time"

	core "github.com/subrag/kv-store/core"
)

func getKeyValueDB() core.DB {
	db := core.DB{KV: core.KV{}}
	return db
}
func TestSetKey(t *testing.T) {
	cases := []struct {
		key     string
		val     string
		ttl     int32
		sleep   int32
		wantErr bool
	}{
		{
			key:     "A",
			val:     "1",
			ttl:     5,
			sleep:   1,
			wantErr: false,
		},
		{
			key:     "B",
			val:     "2",
			ttl:     1,
			sleep:   2,
			wantErr: true,
		},
	}
	db := getKeyValueDB()

	for _, d := range cases {
		db.Set(d.key, d.val, d.ttl)
		time.Sleep(time.Second * time.Duration(d.sleep))
		db.Evict()
		val := db.Get(d.key)
		if d.wantErr {
			if val != "" {
				t.Fail()
			}
		} else {
			if val == "" {
				t.Fail()
			}
		}
	}
}
