package binarytree

import "fmt"

type Node struct {
	Left, Right *Node
	Value       int
}

// Inserts Node
func insert(node *Node, value int) *Node {
	//If there is no node. Returns node. (Inserted)
	if node == nil {
		return &Node{Value: value}
	} else if value < node.Value {
		//If the value is less than the given nodes value, recursively runs this function until it reaches a node where the value is not less than this value
		node.Left = insert(node.Left, value)
	} else if value > node.Value {
		//If the value is greater than the given nodes value, recursively runs this function until it reaches a node where the value is not greater than this value
		node.Right = insert(node.Right, value)
	}
	//returns the given node
	return node
}

// Displays nodes in value order
func display(root *Node) {
	if root != nil {
		display(root.Left)
		fmt.Println(root.Value)
		display(root.Right)
	}
}

// Searches from the root Node for the value requested
func search(root *Node, data int) bool {
	if root == nil {
		return false
	} else if root.Value == data {
		return true
	} else if root.Value > data {
		return search(root.Left, data)
	} else {
		return search(root.Right, data)
	}
}

// Deletes a nodes data
func delete(root *Node, data int) *Node {
	//If not root node. Returns root
	if root == nil {
		return root
	} else if root.Value > data {
		//If the value is greater than data. Runs Delete Recursively on the Left Node of Root
		root.Left = delete(root.Left, data)
	} else if root.Value < data {
		//If the value is less than data. Runs delete recursively on the right. This pattern cause
		root.Right = delete(root.Right, data)
	} else {
		//Once the data matches.
		if root.Left == nil && root.Right == nil {
			//If there is no Right or left, just clear the node
			root = nil
		} else if root.Right != nil { //Find Successor to replace this node
			//If the right node != nil run successor
			root.Value = successor(root)
			root.Right = delete(root.Right, root.Value)
		} else {
			//If the left value is not nil run predecessor
			root.Value = predecessor(root)
			root.Left = delete(root.Left, root.Value)
		}
	}
	return root
}

// Find the node that is Right than the node that is the furthest left of right. This replaces the initial
func successor(root *Node) int {
	root = root.Right
	for root.Left != nil {
		root = root.Left
	}
	return root.Value
}

// Finds the node that is Left than the node that is the furthest right of left. This replaces the root
func predecessor(root *Node) int {
	root = root.Left
	for root.Right != nil {
		root = root.Right
	}
	return root.Value
}
