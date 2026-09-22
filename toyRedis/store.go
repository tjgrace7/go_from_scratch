package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var ErrKeyNotFound = errors.New("Key Not Found")
var Path = "/Users/roxystudio/Public/go_from_scratch/toyRedis/append.txt"

type Vertex struct {
	value string
	t     string
}

type Store struct {
	m  map[string]Vertex
	mu sync.RWMutex
}

func Create() *Store {
	return &Store{m: map[string]Vertex{}, mu: sync.RWMutex{}}
}

func (s *Store) Set(command []string) (Vertex, error) {
	var expiryTime time.Time
	var v Vertex
	if len(command) >= 2 {
		key := command[0]
		value := command[1]

		v.value = value

		if len(command) > 3 {
			num, err := strconv.Atoi(command[3])
			if command[2] != "EX" {
				return v, fmt.Errorf("Improper Format. Expected Blank or EX")
			} else if err != nil {
				t, err := time.Parse(time.RFC3339, command[3])
				if err != nil {
					return v, fmt.Errorf("Expected integer or RFC3339 time value: %v", err)
				} else {
					expiryTime = t
				}
			} else {
				duration := time.Duration(num) * time.Second
				expiryTime = time.Now().Add(duration)
			}
		}
		if !expiryTime.IsZero() {
			v.t = expiryTime.Format(time.RFC3339)
		} else {
			v.t = ""
		}
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.m[key]

	t, err := parseExpiry(value.t)
	if err != nil {
		return "", err
	}
	if t.Before(time.Now()) && !t.IsZero() {
		delete(s.m, key)

		return "", fmt.Errorf("Key Expired: Deleted")
	}
	if exists {
		return value.value, nil
	} else {
		return "", ErrKeyNotFound
	}

}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)

}
func (s *Store) CleanUp() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.m {
		t, err := parseExpiry(v.t)
		if err != nil {
			continue
		}
		if t.Before(time.Now()) && !t.IsZero() {
			s.Delete(k)
		}
	}

}
