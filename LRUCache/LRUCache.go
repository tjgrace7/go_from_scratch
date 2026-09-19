package lrucache

import (
	"fmt"

	"github.com/tjgrace7/go_from_scratch/LinkedList"
	"github.com/tjgrace7/go_from_scratch/hashmaps"
)

type LRUNode[T comparable] struct {
	Key   string
	Value T
}

type LRU[T comparable] struct {
	hashmap hashmaps.Hashmap[*LinkedList.Node[LRUNode[T]]]
	link    LinkedList.LinkedList[LRUNode[T]]
	size    int
}

func CreateLRU[T comparable](size int) LRU[T] {
	return LRU[T]{size: size, link: LinkedList.LinkedList[LRUNode[T]]{}, hashmap: hashmaps.Initiate[*LinkedList.Node[LRUNode[T]]](size)}
}

func (l *LRU[T]) Put(key T, value T) {
	currentNode, err := l.Get(key)
	skey := fmt.Sprintf("%v", key)
	data := hashmaps.Data[*LinkedList.Node[LRUNode[T]]]{Key: skey, Value: &LinkedList.Node[LRUNode[T]]{Data: LRUNode[T]{Key: skey, Value: value}}}
	//Overwrite Data
	if err != nil {
		l.link.PushFront(LRUNode[T]{Key: skey, Value: value})

	} else {
		currentNode.Data.Value = value
		l.link.PushFront(currentNode.Data)
		l.hashmap.DeleteKey(data.Key)
	}
	l.hashmap.AddtoMap(data)
	fmt.Println("LRU Size:", l.size, "Hashmap Occupancy:", l.hashmap.OccupiedCount+l.hashmap.SecondaryOccupiedCount)
	//Evict if Occupied is greater than size
	if l.size < l.hashmap.OccupiedCount+l.hashmap.SecondaryOccupiedCount {
		err := l.hashmap.DeleteKey(l.link.Tail.Data.Key)
		if err != nil {
			fmt.Println("Error Deleting User:", err)
		}
		l.link.DeleteNode(*l.link.Tail)

	}
}
func (l *LRU[T]) Get(key T) (*LinkedList.Node[LRUNode[T]], error) {
	skey := fmt.Sprintf("%v", key)
	data, err := l.hashmap.SearchMap(skey)
	if err != nil {
		return nil, err
	}
	fmt.Println("Data: ", data.Key)
	node, _, er := l.link.SearchData(data.Value.Data)
	if er != nil {
		return nil, er
	}
	l.LeastRecentUsed(*node)
	return node, nil
}
func (l *LRU[T]) LeastRecentUsed(node LinkedList.Node[LRUNode[T]]) {
	l.link.DeleteNode(node)
	l.link.PushFront(node.Data)
}
