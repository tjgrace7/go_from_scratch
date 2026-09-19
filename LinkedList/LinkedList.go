package LinkedList

import "fmt"

type Node[T comparable] struct {
	Left, Right *Node[T]
	Data        T
}
type LinkedList[T comparable] struct {
	Head *Node[T]
	Tail *Node[T]
	Size int
}

func NewNode[T comparable](value T) *Node[T] {
	return &Node[T]{Data: value}
}
func (l *LinkedList[T]) insertLeft(data T, node *Node[T]) *Node[T] {
	newNode := &Node[T]{Data: data, Right: node}
	if node == nil {
		newNode = l.PushFront(data)
		return newNode
	}

	if node.Left == nil {
		node.Left = newNode
		l.Head = newNode
	} else {
		node.Left.Right = newNode
		newNode.Left = node.Left
		node.Left = newNode
	}
	l.Size++
	return newNode
}
func (l *LinkedList[T]) insertRight(data T, node *Node[T]) *Node[T] {
	newNode := &Node[T]{Data: data, Left: node}
	if node == nil {
		newNode = l.PushBack(data)
		return newNode
	}

	if node.Right == nil {
		node.Right = newNode
		l.Tail = newNode
	} else {
		node.Right.Left = newNode
		newNode.Right = node.Right
		node.Right = newNode
	}
	l.Size++
	return newNode

}
func (l *LinkedList[T]) PushFront(data T) *Node[T] {
	newNode := &Node[T]{Data: data}
	fmt.Println("Node Data:", data)
	if l.Head == nil {
		l.Head = newNode
		l.Tail = newNode
		l.Size++
	} else {
		newNode = l.insertLeft(data, l.Head)
		l.Head = newNode
	}
	return newNode
}
func (l *LinkedList[T]) PushBack(data T) *Node[T] {

	newNode := &Node[T]{Data: data}
	if l.Tail == nil {
		l.Head = newNode
		l.Tail = newNode
		l.Size++
	} else {
		newNode = l.insertRight(data, l.Tail)
		l.Tail = newNode
	}
	return newNode
}
func (l *LinkedList[T]) InsertIndex(data T, index int) (*Node[T], error) {
	if l.Size <= index {
		return nil, fmt.Errorf("Index Out of Range")
	} else {
		currentNode := l.Head
		for i := 0; i < index; i++ {
			currentNode = currentNode.Right
		}
		newNode := l.insertLeft(data, currentNode)
		return newNode, nil
	}
}
func (l *LinkedList[T]) SearchData(data T) (*Node[T], int, error) {
	currentNode := l.Head
	fmt.Println("Search Data:", data, "Current Node:", currentNode.Data, "Linked List Size:", l.Size)
	for i := 0; i < l.Size; i++ {
		if data == currentNode.Data {
			return currentNode, i, nil
		} else if currentNode.Right == nil {
			break
		}
		currentNode = currentNode.Right
	}
	return nil, -1, fmt.Errorf("Data not in Linked List")
}
func (l *LinkedList[T]) DeleteData(data T) error {
	deleteNode, _, err := l.SearchData(data)
	if err != nil {
		return fmt.Errorf("Data Not Found, Did Not Delete")
	}
	l.DeleteNode(*deleteNode)
	return nil
}
func (l *LinkedList[T]) DeleteNode(node Node[T]) error {
	if node.Left == nil && node.Right == nil && l.Size == 1 {
		l.Head = nil
		l.Tail = nil
		l.Size = 0
	} else if node.Left == nil && node.Right == nil {
		//Node is not in linked list

		return fmt.Errorf("Node not in Linked List")
	} else if node.Left != nil && node.Right != nil {
		node.Left.Right = node.Right
		node.Right.Left = node.Left
		l.Size--
	} else if node.Left != nil && node.Right == nil {
		node.Left.Right = nil
		l.Tail = node.Left
		l.Size--
	} else {
		node.Right.Left = nil
		l.Head = node.Right
		l.Size--
	}
	return nil
}

// If target Index is out of Range, it will push it to the back of the LinkedList instead of erroring
func (l *LinkedList[T]) ReIndex(data T, targetIndex int) error {
	_, _, err := l.SearchData(data)
	if err != nil {
		return err
	}
	l.DeleteData(data)
	if targetIndex >= l.Size {
		l.PushBack(data)
	} else {
		l.InsertIndex(data, targetIndex)
	}
	return nil
}
func (l *LinkedList[T]) Display() {
	currentNode := l.Head
	for i := 0; i <= l.Size-1; i++ {
		fmt.Println(currentNode.Data)
		currentNode = currentNode.Right
	}
}
