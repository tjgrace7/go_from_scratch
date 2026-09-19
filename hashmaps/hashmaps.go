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

func (h *Hashmap[T]) hash(key string) uint32 {
	hash := FNV1aHash(key)
	return hash % uint32(len(h.initial))
}
func Initiate[T comparable](size int) Hashmap[T] {
	initial := make([]Data[T], size)
	secondary := make([]Data[T], 0, size)
	return Hashmap[T]{initial: initial, secondary: secondary}

}
func (h *Hashmap[T]) AddtoMap(data Data[T]) {
	index := h.hash(data.Key)
	data.nextIndex = -1
	data.Occupied = true
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
func (h *Hashmap[T]) SearchMap(key string) (Data[T], error) {
	index := h.hash(key)
	if h.initial[index].Key == key {
		return h.initial[index], nil

	} else if !h.initial[index].Occupied {
		return h.initial[index], fmt.Errorf("Key not Found")
	} else {

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
func (h *Hashmap[T]) Resize(sizefactor int, loadfactor float32) bool {
	fmt.Println("Occupied Count:", h.OccupiedCount)
	if float32(h.OccupiedCount) <= float32(len(h.initial))*loadfactor {
		return false
	}
	initial := h.initial
	h.initial = make([]Data[T], len(h.initial)*sizefactor)
	h.OccupiedCount = 0
	for i := 0; i < len(initial); i++ {
		if initial[i].Occupied {
			h.AddtoMap(initial[i])
		}
	}
	secondary := h.secondary
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

func (h *Hashmap[T]) DeleteKey(key string) error {
	index := h.hash(key)
	currentData := h.initial[index]
	initial := true
	var zero T
	currentIndex := int(index)
	for currentData.Key != key {

		if currentData.nextIndex == -1 {
			return fmt.Errorf("Key Not Found: Delete Key Function. Key: ")
		}
		currentIndex = currentData.nextIndex
		currentData = h.secondary[currentData.nextIndex]
		initial = false
	}

	nextIndex := currentData.nextIndex
	if initial && nextIndex == -1 {
		h.initial[index] = Data[T]{Key: "", Value: zero, Occupied: false, nextIndex: -1}
		h.OccupiedCount--
		return nil
	} else if initial && nextIndex != -1 {
		h.initial[index] = h.secondary[nextIndex]
		h.initial[index].nextIndex = nextIndex
		currentIndex = nextIndex
		nextIndex = h.secondary[currentIndex].nextIndex
	}
	//If the Delete key is on the first iteration of the delete process in secondary. Primary key needs update
	first := true

	for nextIndex != -1 {

		h.secondary[currentIndex] = h.secondary[nextIndex]

		if h.secondary[nextIndex].nextIndex != -1 {
			h.secondary[currentIndex].nextIndex = nextIndex

		}
		currentIndex = nextIndex
		nextIndex = h.secondary[nextIndex].nextIndex
		first = false
	}
	//Resets initial index to having no follow up if the secondary only has one chain
	if first {
		fmt.Println("Reset initial next index:", h.initial[index].Key)
		h.initial[index].nextIndex = -1
	}
	h.secondary[currentIndex] = Data[T]{Key: "", Value: zero, Occupied: false, nextIndex: -1}
	h.SecondaryOccupiedCount--
	return nil

}
