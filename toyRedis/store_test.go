package main

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestGetSetDelete(t *testing.T) {
	s := Create()

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
	time.Sleep(time.Duration(num) * time.Second)
	_, err = s.Get("name")
	if err == nil {
		t.Error("Key Should Not Exist", err)
	}
}

func TestCleanUp(t *testing.T) {
	s := Create()
	s.Set([]string{"name", "Tyler", "EX", "15"})

	s.Set([]string{"age", "27", "EX", "60"})
	s.Set([]string{"home", "208 Maple Ave", "EX", "3"})
	time.Sleep(time.Duration(5) * time.Second)
	_, err := s.Get("home")
	if err == nil {
		t.Error("Key Should Be Expired")
	}
	time.Sleep(time.Duration(10) * time.Second)
	s.CleanUp()
	_, err = s.Get("name")
	//Tests Clean Up

	if err != ErrKeyNotFound {
		t.Error("Cleanup Should have made key not exist", err)
	}
	_, err = s.Get("age")
	if err != nil {
		t.Error("Key Should Exist")
	}
}

func TestConcurrency(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := Create()
		var wg sync.WaitGroup

		s.Set([]string{"name", "tyler", "EX", "1"})
		time.Sleep(time.Duration(1) * time.Second)
		wg.Add(2)
		go func() { defer wg.Done(); s.Get("name") }()
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

func BenchmarkConcurrency(b *testing.B) {

	s := Create()
	var wg sync.WaitGroup

	for b.Loop() {
		wg.Add(10)

		go func() { defer wg.Done(); s.Set([]string{"name", "Tyler"}) }()
		go func() { defer wg.Done(); s.Get("name") }()
		go func() { defer wg.Done(); s.Set([]string{"age", "26"}) }()
		go func() { defer wg.Done(); s.Set([]string{"age", "27"}) }()
		go func() { defer wg.Done(); s.Get("age") }()
		go func() { defer wg.Done(); s.Delete("age") }()
		go func() { defer wg.Done(); s.Set([]string{"boobs", "perky"}) }()
		go func() { defer wg.Done(); s.Set([]string{"boobs", "saggy"}) }()
		go func() { defer wg.Done(); s.Get("boobs") }()
		go func() { defer wg.Done(); s.Delete("boobs") }()
		wg.Wait()
	}
}
func BenchmarkSequential(b *testing.B) {

	s := Create()

	for b.Loop() {

		s.Set([]string{"name", "Tyler"})
		s.Get("name")
		s.Set([]string{"age", "26"})
		s.Set([]string{"age", "27"})
		s.Get("age")
		s.Delete("age")
		s.Set([]string{"boobs", "perky"})
		s.Set([]string{"boobs", "saggy"})
		s.Get("boobs")
		s.Delete("boobs")
	}

}
