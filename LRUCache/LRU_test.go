package lrucache

import (
	"testing"
)

func TestEviction(t *testing.T) {
	lru := CreateLRU[string](5)
	test := []struct {
		Key   string
		Value string
	}{
		{Key: "tjgrace", Value: "1"},
		{Key: "tyler", Value: "2"},
		{Key: "eve", Value: "3"},
		{Key: "john", Value: "4"},
		{Key: "elvis", Value: "5"},
		{Key: "willis", Value: "6"},
		{Key: "August", Value: "7"},
		{Key: "Hannah", Value: "8"},
		{Key: "george", Value: "9"},
		{Key: "will", Value: "10"},
	}
	//Make tjgrace maintain teh entire put method
	for i := 0; i < len(test)-1; i++ {
		lru.Put(test[i].Key, test[i].Value)
		_, err := lru.Get("tjgrace")
		if err != nil {
			t.Error("Error Retrieving tjgrace", err)
		}
		if test[i].Key == "August" {
			_, err = lru.Get("eve")
			if err == nil || err.Error() != "Key not Found" {

				t.Error("eve Found. Should be Evivted. Error", err)
			}
		}
	}
	_, e := lru.Get("tyler")
	if e == nil {
		t.Error("tyler found. Should be Evicted")
	}
}

func TestPutGet(t *testing.T) {
	lru := CreateLRU[string](5)
	test := []struct {
		Matching string
		Key      string
		Value    string
		WantErr  bool
	}{
		{Key: "tjgrace", Value: "1", WantErr: false},
		{Key: "tyler", Value: "2", WantErr: true},
		{Key: "eve", Value: "3", WantErr: false},
		{Key: "john", Value: "4", WantErr: false},
		{Key: "elvis", Value: "5", WantErr: false},
		{Key: "willis", Value: "6", WantErr: false},
		{Key: "August", Value: "7", WantErr: false},
		{Key: "Hannah", Value: "8", WantErr: false},
		{Key: "george", Value: "9", WantErr: false},
		{Key: "will", Value: "10", WantErr: false},
	}
	for _, u := range test {
		lru.Put(u.Key, u.Value)
		if u.WantErr {
			u.Matching = lru.link.Tail.Data.Key
		} else {
			u.Matching = lru.link.Head.Data.Key
		}
		t.Log("Object Put")
		node, e := lru.Get(u.Key)
		if e != nil {
			t.Errorf("LRU Put & Get Requests: %v, For user: %v, Want Error: %v", e, u.Key, u.WantErr)
		}
		if (u.Matching != u.Key) != u.WantErr {
			t.Errorf("LRU Does Not Match Wanted Key: %v, Current Key: %v, Want Error: %v", u.Matching, u.Key, u.WantErr)
		}
		t.Log("Node:", node)
	}

}
