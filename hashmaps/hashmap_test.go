package hashmaps

import (
	"fmt"
	"math/rand"
	"testing"
)

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
func TestDeleteKey(t *testing.T) {
	users := []Data[string]{
		{Key: "tjgrace", Value: "0"},
		{Key: "tyler", Value: "1"},
		{Key: "eve", Value: "2"},
		{Key: "john", Value: "3"},
		{Key: "elvis", Value: "4"},
		{Key: "willis", Value: "5"},
		{Key: "August", Value: "6"},
		{Key: "Hannah", Value: "7"},
		{Key: "george", Value: "8"},
		{Key: "will", Value: "9"},
	}
	hmap := Initiate[string](10)
	for _, u := range users {
		hmap.AddtoMap(u)
	}
	data, err := hmap.SearchMap("willis")
	if err != nil {
		t.Error("Error Finding User")
	}
	fmt.Println("User Found:", data.Key)
	//The order matters. This order checks every type of delete
	//Delete in the initial with no keys in secondary
	//Delete in the 1st chain.
	//Final Delete.
	//Deleting in the initial with 1 secondary and resetting the intial to having no next index
	numbers := []int{0, 5, 3, 4, 2, 6, 1, 9, 8, 7}
	for _, n := range numbers {
		u := users[n]
		hmap.DeleteKey(u.Key)
		_, err := hmap.SearchMap(u.Key)
		if err == nil {
			t.Errorf("Should not be able to access users: %v", err)
		}
	}

}

// Collision Chain Testing
func TestCollisionChain(t *testing.T) {
	//Change this variable to one of the variables in the slice to see where it lands in the hashmap. If you use the djb2 algorithm, eve should iterate through the 2nd slice in the hashmap twice before returning
	users := []Data[string]{
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
	hmap := Initiate[string](10)

	for _, u := range users {
		hmap.AddtoMap(u)
	}
	testusernams := []string{
		//A layer is the number of times a nextUserIndex has been called from the user struct
		//Eva will return a user who is two layers deep in the hash map
		//Henry will return a user who is not in the Hashmap, but will match a Hash code with eve, tyler, and tjgrace. This will iterate through the 2nd slice several times
		//tjgrace will be found directly in the initial map
		"eve", "Henry", "tjgrace", "mark",
	}
	for _, u := range testusernams {
		t.Run(u, func(t *testing.T) {
			username, err := hmap.SearchMap(u)

			if u == "Henry" || u == "mark" {
				if err == nil {
					t.Errorf("expected error for %q, git user %v", u, username)
				} else {
					fmt.Println("User not Found, Success")
				}
			} else if err != nil {
				t.Errorf("Error Searching for User %q", err)
			}
			fmt.Println("Username", username.Key, "Value", username.Value)
		})
	}

}
func TestResizeSearch(t *testing.T) {
	var hmap = Initiate[string](10)
	users := []Data[string]{
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
	//This will resize the hashmap to fit users in a different spot and search for a new user
	for i := 0; i < len(users); i++ {
		hmap.AddtoMap(users[i])
		//Change 0.1 to 0.9 to avoid resize
		hmap.Resize(2, 0.1)

	}
	us, err := hmap.SearchMap("eve")
	if err != nil {
		t.Error("Error finding user")
	}
	fmt.Println("User found", us)
}
func TestResize(t *testing.T) {
	var hmap = Initiate[string](10)
	users := []Data[string]{
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
	//This will resize the hashmap to fit users in a different spot and search for a new user
	for i := 0; i < len(users); i++ {
		hmap.AddtoMap(users[i])
		//Change 0.1 to 0.9 to avoid resize
		loadfactor := 0.1
		preResize := hmap.OccupiedCount
		initialLength := len(hmap.initial)
		resized := hmap.Resize(2, float32(loadfactor))

		//Did the occupancy actually trigger a resize or was it loop iterations?
		if resized && float32(preResize) <= float32(initialLength)*float32(loadfactor) {
			t.Error("Should not have resized")
		}
	}
}

func BenchmarkXxx(b *testing.B) {
	hmap := Initiate[string](1000000)
	var testuser Data[string]
	for b.Loop() {
		var u Data[string]
		i := rand.Intn(10000)
		u.Key = randomString(15)
		u.Value = "12345"
		hmap.AddtoMap(u)
		hmap.Resize(2, 0.75)
		//Takes user as test user for searching where the most recent iteration of i == 219 happens.
		//This would obviously not be used in a loading settings, but is useful for finding a random test user that is generated in this loop.
		if i == 219 {
			testuser = u
		}

	}
	if testuser.Key == "" {
		fmt.Println("Error: No user entered into test. Random integer didn't hit i")
		return
	}
	fmt.Println("Test User is: ", testuser.Key)
	fmt.Println("hashmap2 length: ", len(hmap.secondary))
	founduser, err := hmap.SearchMap(testuser.Key)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Found User: ", founduser.Key, " Password: ", founduser.Value)
}
