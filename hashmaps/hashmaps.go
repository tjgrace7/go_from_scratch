package hashmaps

import "fmt"

type Data[T comparable] struct {
	Value     T
	Key       string
	Occupied  bool
	nextIndex int
}
type Hashmap[T comparable] struct {
	initial                []Data[T]
	secondary              []Data[T]
	OccupiedCount          int
	SecondaryOccupiedCount int
}

// Hash returns the Hash code using FNV1a
func (h *Hashmap[T]) hash(key string) uint32 {
	hash := FNV1aHash(key)
	return hash % uint32(len(h.initial))
}

// Creates a HashMap
func Initiate[T comparable](size int) Hashmap[T] {
	initial := make([]Data[T], size)
	secondary := make([]Data[T], 0, size)
	return Hashmap[T]{initial: initial, secondary: secondary}

}

// Adds given data to map
func (h *Hashmap[T]) AddtoMap(data Data[T]) {
	index := h.hash(data.Key)
	//Uses index = -1 instead of 0
	data.nextIndex = -1
	data.Occupied = true
	//If the initial map is occupied begin collision chain
	if h.initial[index].Occupied {
		h.secondary = append(h.secondary, data)
		h.SecondaryOccupiedCount++
		currentData := h.initial[index]
		currentIndex := int(index)
		initial := true
		for currentData.nextIndex != -1 {
			currentIndex = currentData.nextIndex
			currentData = h.secondary[currentData.nextIndex]
			initial = false
		}
		if initial {
			h.initial[index].nextIndex = len(h.secondary) - 1
		} else {
			h.secondary[currentIndex].nextIndex = len(h.secondary) - 1
		}
	} else {
		h.initial[index] = data
		h.OccupiedCount++
	}
}

// Searches map for data
func (h *Hashmap[T]) SearchMap(key string) (Data[T], error) {
	index := h.hash(key)
	//If the key is in initial return
	if h.initial[index].Key == key {
		return h.initial[index], nil

	} else if !h.initial[index].Occupied {
		//No key found
		return h.initial[index], fmt.Errorf("Key not Found")
	} else {
		//Begin Collision Chain
		layers := 1
		currentData := h.initial[index]
		for currentData.Key != key {
			if currentData.nextIndex == -1 {
				return currentData, fmt.Errorf("Key not Found")
			}
			currentData = h.secondary[currentData.nextIndex]
			layers++
		}
		return currentData, nil
	}
}

// Resize Hashmap. Loadfactor determins how much of the hashmap is occupied before resizing. Size Factor determines how large to make the new size
func (h *Hashmap[T]) Resize(sizefactor int, loadfactor float32) bool {
	fmt.Println("Occupied Count:", h.OccupiedCount)
	if float32(h.OccupiedCount) <= float32(len(h.initial))*loadfactor {
		return false
	}
	initial := h.initial
	//Creates new slices
	h.initial = make([]Data[T], len(h.initial)*sizefactor)
	h.OccupiedCount = 0
	for i := 0; i < len(initial); i++ {
		if initial[i].Occupied {
			h.AddtoMap(initial[i])
		}
	}
	secondary := h.secondary
	//Creates new slice. Eliminating Dead Delete Keys
	h.secondary = make([]Data[T], 0, len(h.initial))
	h.SecondaryOccupiedCount = 0
	for i := 0; i < len(secondary); i++ {
		if secondary[i].Occupied {
			h.AddtoMap(secondary[i])
		}
	}
	fmt.Println("Occupied Count:", h.OccupiedCount)
	return true
}

// Deletes given key.
func (h *Hashmap[T]) DeleteKey(key string) error {
	index := h.hash(key)
	currentData := h.initial[index]
	initial := true
	var zero T
	currentIndex := int(index)
	//If Key is not in intial list. Begin Chaining
	for currentData.Key != key {

		if currentData.nextIndex == -1 {
			return fmt.Errorf("Key Not Found: Delete Key Function. Key: ")
		}
		currentIndex = currentData.nextIndex
		currentData = h.secondary[currentData.nextIndex]
		initial = false
	}
	//Current Data no shows the given key
	nextIndex := currentData.nextIndex
	//If the key was found in initial and there is no chain. Delete initial and lower occupied count then return
	if initial && nextIndex == -1 {
		h.initial[index] = Data[T]{Key: "", Value: zero, Occupied: false, nextIndex: -1}
		h.OccupiedCount--
		return nil
	} else if initial && nextIndex != -1 {
		//If found in intial copy next index into initial
		h.initial[index] = h.secondary[nextIndex]
		h.initial[index].nextIndex = nextIndex
		currentIndex = nextIndex
		nextIndex = h.secondary[currentIndex].nextIndex
	}
	//If the Delete key is on the first iteration of the delete process in secondary. Primary key needs update
	first := true
	//For as long as the next index != -1. Copy the key into the key before it
	for nextIndex != -1 {
		h.secondary[currentIndex] = h.secondary[nextIndex]
		if h.secondary[nextIndex].nextIndex != -1 {
			h.secondary[currentIndex].nextIndex = nextIndex

		}
		currentIndex = nextIndex
		nextIndex = h.secondary[nextIndex].nextIndex
		first = false
	}
	//Resets initial index to having no follow up if no additional keys in secondary
	if first {
		fmt.Println("Reset initial next index:", h.initial[index].Key)
		h.initial[index].nextIndex = -1
	}
	//deletes final key. Leaving Dead Key in Secondary. This will no be in a chain, but if you look up the index, it will be empty
	h.secondary[currentIndex] = Data[T]{Key: "", Value: zero, Occupied: false, nextIndex: -1}
	h.SecondaryOccupiedCount--
	return nil

}
