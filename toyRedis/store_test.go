package main

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

// Tests the GET/SET/DELETE Methods
func TestGetSetDelete(t *testing.T) {
	s := Create()
	current := time.Now()
	s.now = func() time.Time { return current }
	//Set/GET Test
	s.Set([]string{"name", "tyler"})

	_, err := s.Get("name")
	if err != nil {
		t.Error("Error Getting Name", err)
	}
	//Delete Test
	s.Delete("name")
	_, err = s.Get("name")
	if err == nil {
		t.Error("Name Not Deleted", err)
	}
	//Get Blank Test
	_, err = s.Get("hulu")
	if err == nil {
		t.Error("Key Does Not Exist. Should not Be Found", err)
	}
	//Expiration test. Duration in Seconds
	duration := "10"
	num, err := strconv.Atoi(duration)
	if err != nil {
		t.Error("Error Converting String:", err)
	}
	//EX required to set Expiration
	_, err = s.Set([]string{"pilot", "George", "EX", duration})
	if err != nil {
		t.Error("Error Creating Key/Value Pair:", err)
	}
	_, err = s.Get("pilot")
	if err != nil {
		t.Error("Key Should Be Found", err)
	}
	current = current.Add(time.Duration(num) * time.Second)
	_, err = s.Get("pilot")
	if err == nil {
		t.Error("Key Should Not Exist", err)
	}
}

// Tests Clean up function
func TestCleanUp(t *testing.T) {
	s := Create()
	current := time.Now()
	s.now = func() time.Time { return current }
	//Sets function with expirations
	s.Set([]string{"name", "Tyler", "EX", "15"})

	s.Set([]string{"age", "27", "EX", "60"})
	s.Set([]string{"home", "208 Maple Ave", "EX", "3"})
	//Sleeps
	current = current.Add(5 * time.Second)
	//Attempts to GET expired Key
	_, err := s.Get("home")
	if err == nil {
		t.Error("Key Should Be Expired")
	}
	current = current.Add(20 * time.Second)
	//Cleans up Expired Keys
	s.CleanUp()
	_, err = s.Get("name")
	//Tests Clean Up. Making sure the Error Matches Key Not Found, instead of Expired
	if err != ErrKeyNotFound {
		t.Error("Cleanup Should have made key not exist", err)
	}
	_, err = s.Get("age")
	if err != nil {
		t.Error("Key Should Exist")
	}
}

// Tests Concurrent Functions for Race Conditions
func TestConcurrency(t *testing.T) {
	//Loops 10 times
	for i := 0; i < 10; i++ {
		s := Create()
		current := time.Now()
		s.now = func() time.Time { return current }
		var wg sync.WaitGroup

		s.Set([]string{"name", "tyler", "EX", "1"})
		current = current.Add(1 * time.Second)
		wg.Add(2)
		//Attempts to GET Expired Key
		go func() { defer wg.Done(); s.Get("name") }()
		//Sets New Key
		go func() { defer wg.Done(); s.Set([]string{"name", "George"}) }()
		wg.Wait()

		value, err := s.Get("name")
		if err != nil {
			t.Error("Expected Value:", err)
		}
		if value != "George" {
			t.Errorf("iteration %d: expected George, got %q", i, value)
		}
	}
}

// Benchmark speed of Concurrency
func BenchmarkConcurrency(b *testing.B) {

	s := Create()
	var wg sync.WaitGroup

	for b.Loop() {
		wg.Add(10)
		//10 Concurrent Functions have to wait in line
		go func() { defer wg.Done(); s.Set([]string{"name", "Tyler"}) }()
		go func() { defer wg.Done(); s.Get("name") }()
		go func() { defer wg.Done(); s.Set([]string{"age", "26"}) }()
		go func() { defer wg.Done(); s.Set([]string{"age", "27"}) }()
		go func() { defer wg.Done(); s.Get("age") }()
		go func() { defer wg.Done(); s.Delete("age") }()
		go func() { defer wg.Done(); s.Set([]string{"terrain", "woods"}) }()
		go func() { defer wg.Done(); s.Set([]string{"terrain", "river"}) }()
		go func() { defer wg.Done(); s.Get("terrain") }()
		go func() { defer wg.Done(); s.Delete("terrain") }()
		wg.Wait()
	}
}

// Benchmarks same 10 functions as Concurrency ran Sequentially
func BenchmarkSequential(b *testing.B) {

	s := Create()

	for b.Loop() {

		s.Set([]string{"name", "Tyler"})
		s.Get("name")
		s.Set([]string{"age", "26"})
		s.Set([]string{"age", "27"})
		s.Get("age")
		s.Delete("age")
		s.Set([]string{"terrain", "woods"})
		s.Set([]string{"terrain", "river"})
		s.Get("terrain")
		s.Delete("terrain")
	}

}
