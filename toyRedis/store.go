package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var ErrKeyNotFound = errors.New("Key Not Found")

// The Vertex sets a value strings and a time string for the map
type Vertex struct {
	value string
	t     string
}

// The store struct stores the map and a sync.RWMutex
type Store struct {
	m   map[string]Vertex
	mu  sync.RWMutex
	now func() time.Time
}

// Creates a pointer to Store
func Create() *Store {
	return &Store{m: map[string]Vertex{}, mu: sync.RWMutex{}, now: time.Now}
}

// Sets a function in the map
func (s *Store) Set(command []string) (Vertex, error) {
	var expiryTime time.Time
	var v Vertex
	//Makes sure there are at least 2 strings in command to run code
	if len(command) >= 2 {
		//sets the key to zero and the value to 1
		key := command[0]
		value := command[1]

		v.value = value
		//If there are more than 3 strings in command checks to see if it is formatted for expiration
		if len(command) > 3 {
			num, err := strconv.Atoi(command[3])
			//Requires Expiration
			if command[2] != "EX" {
				return v, fmt.Errorf("Improper Format. Expected Blank or EX")
			} else if err != nil {
				//Parses time to text
				t, err := time.Parse(time.RFC3339, command[3])
				if err != nil {
					return v, fmt.Errorf("Expected integer or RFC3339 time value: %v", err)
				} else {
					expiryTime = t
				}
			} else {
				//User passes an integer in seconds declaring how much time to keep the key alive for
				duration := time.Duration(num) * time.Second
				expiryTime = s.now().Add(duration)
			}
		}
		//Checks to see if expiryTime is Zero
		if !expiryTime.IsZero() {
			//If not adds time to v.t
			v.t = expiryTime.Format(time.RFC3339)
		} else {
			v.t = ""
		}

		//Locks the RWMutex to update the key
		s.mu.Lock()
		s.m[key] = v
		s.mu.Unlock()

		return v, nil
	} else {
		return v, fmt.Errorf("Error Finding String")
	}
}

func parseExpiry(s string) (time.Time, error) {
	if s == "" || s == "0" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("Parse Expiry %q: %w", s, err)
	}
	return t, nil
}

func (s *Store) Get(key string) (string, error) {
	//Read Locks Mutex and defer Unlock to end
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.m[key]
	//Parses Vertex.t to see if it is in proper time format
	t, err := parseExpiry(value.t)
	if err != nil {
		return "", err
	}
	//Checks to see if the key is expired
	if t.Before(s.now()) && !t.IsZero() {
		delete(s.m, key)

		return "", fmt.Errorf("Key Expired: Deleted")
	}
	//If the key exists returns value, else it doesn't
	if exists {
		return value.value, nil
	} else {
		return "", ErrKeyNotFound
	}

}

func (s *Store) Delete(key string) {
	//Locks the map to delete the key
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)

}
func (s *Store) CleanUp() {
	//Locks the map to clean up the keys
	//Iterates through key to check for expiration
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.m {
		t, err := parseExpiry(v.t)
		//If no expiration continues to next index
		if err != nil {
			continue
		}
		if t.Before(s.now()) && !t.IsZero() {
			delete(s.m, k)
		}
	}

}
