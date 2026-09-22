package LinkedList

import "fmt"

// Node Determines the Data, the Left Node, and the Right Node
type Node[T comparable] struct {
	Left, Right *Node[T]
	Data        T
}

// The LinkedList has the Head Node, Tail Node, and the size of the List.
// Each Head and Tail Node have right and left nodes. (Head left is nil; Tail Right is nil)
type LinkedList[T comparable] struct {
	Head *Node[T]
	Tail *Node[T]
	Size int
}

// Creates a NewNode
func NewNode[T comparable](value T) *Node[T] {
	return &Node[T]{Data: value}
}

// Inserts a Node Left of the given Node. Or before it
func (l *LinkedList[T]) insertLeft(data T, node *Node[T]) *Node[T] {
	newNode := &Node[T]{Data: data, Right: node}
	//If there is no Node. Pushes Front
	if node == nil {
		newNode = l.PushFront(data)
		return newNode
	}
	//If the node given does not have a left node. It makes the New Node head of the Linked List & assigns the left node to the new Node
	if node.Left == nil {
		node.Left = newNode
		l.Head = newNode
	} else {
		//If the Node Does have a left node. Assigns the Node.Left to NewNode.Left. The Right Node of the Left node to the newNode. The Left Node of the original node to the New Node.
		node.Left.Right = newNode
		newNode.Left = node.Left
		node.Left = newNode
	}
	l.Size++
	return newNode
}

// Inserts a Node to the Right of a given node.
func (l *LinkedList[T]) insertRight(data T, node *Node[T]) *Node[T] {
	newNode := &Node[T]{Data: data, Left: node}
	//If there is no Node. Pushes to the Back of the Linked List.
	if node == nil {
		newNode = l.PushBack(data)
		return newNode
	}
	//If the node.Right is empty. Assigns the node.Right = newNode. Makes new Node the tail
	if node.Right == nil {
		node.Right = newNode
		l.Tail = newNode
	} else {
		//If the Node is not Tail, assigns the node.Right's left node to the newNode(Insertion). Then the newNode.Right = node.Right. Then node.Right becomes the NewNode. Inserted into the LinkedList
		node.Right.Left = newNode
		newNode.Right = node.Right
		node.Right = newNode
	}
	l.Size++
	return newNode

}

// Pushes to the front of the Linked List
func (l *LinkedList[T]) PushFront(data T) *Node[T] {
	newNode := &Node[T]{Data: data}
	fmt.Println("Node Data:", data)
	//If there is no Head, The List is empty. Assigns New Node to Head & Tail. Increments Size
	if l.Head == nil {
		l.Head = newNode
		l.Tail = newNode
		l.Size++
	} else {
		//Otherwise adds it to the front.
		//Inserts the node left of the current Head
		newNode = l.insertLeft(data, l.Head)
		//Resets the Head to the New Node
		l.Head = newNode
	}
	return newNode
}
func (l *LinkedList[T]) PushBack(data T) *Node[T] {

	newNode := &Node[T]{Data: data}
	//Empty Case must increment size of List
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

// Inserts a function at a selected Index. The index must be less than the length of the List
func (l *LinkedList[T]) InsertIndex(data T, index int) (*Node[T], error) {
	if l.Size <= index {
		return nil, fmt.Errorf("Index Out of Range")
	} else {

		if index < l.Size/2 {
			currentNode := l.Head
			for i := 0; i < index; i++ {
				currentNode = currentNode.Right
			}
			newNode := l.insertLeft(data, currentNode)
			return newNode, nil
		} else {
			currentNode := l.Tail
			for i := l.Size - 1; i > index; i-- {
				currentNode = currentNode.Left
			}
			newNode := l.insertLeft(data, currentNode)
			return newNode, nil
		}

	}
}

// Searches the List Left to Right for the data
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

// Delete Data from the Data.
func (l *LinkedList[T]) DeleteData(data T) error {
	deleteNode, _, err := l.SearchData(data)
	if err != nil {
		return fmt.Errorf("Data Not Found, Did Not Delete")
	}
	l.DeleteNode(*deleteNode)
	return nil
}

// Deletes a Node. This is faster than data O(1)
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

// Takes data at an index and puts it in a new index. If Data not Found returns error
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

// Displays every node in Linked List
func (l *LinkedList[T]) Display() {
	currentNode := l.Head
	for i := 0; i <= l.Size-1; i++ {
		fmt.Println(currentNode.Data)
		currentNode = currentNode.Right
	}
}
