package hashmaps

// FNV1a Hash Alogorithm
func FNV1aHash(key string) uint32 {
	var hash uint32 = 2166136261 //offset basis, fixed constant
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= 16777619
	}
	return hash
}

// Djb2 Hash Algorithm
func Djb2(key string) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(key); i++ {
		hash = hash*33 + uint32(key[i])
	}
	return hash
}

// Djb2a Hash Algorithm
func Djb2a(key string) uint32 {
	var hash uint32 = 5381
	for i := 0; i < len(key); i++ {
		hash = hash*33 ^ uint32(key[i])
	}
	return hash
}
