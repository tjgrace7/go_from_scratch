package LinkedList

import (
	"fmt"
	"testing"
)

// Int LinkedList Insert

func TestIntegerLinkedList(t *testing.T) {
	values := []int{4, 7, 1, 30, 9, 20, 0, 89, 50}
	link := &LinkedList[int]{}
	for _, v := range values {
		link.PushBack(v)
	}

	node, index, err := link.SearchData(9)
	if err != nil {
		t.Error(err)

	}
	fmt.Println("Node Found ", node, "at index", index)

	node, index, err = link.SearchData(90)
	if err == nil {
		t.Error("Expecting Error Here. Search Data not in list")
	}
	fmt.Println("Node Found ", node, "at index", index)
	err = link.DeleteData(20)
	if err != nil {
		t.Error(err)
	}
	link.PushFront(500)
	link.ReIndex(89, 7)

	link.Display()
}

func TestSearchReindex(t *testing.T) {
	values := []string{"hello", "world", "tyler", "james", "hannah", "chloe", "programming"}
	link := &LinkedList[string]{}
	for _, v := range values {
		link.PushBack(v)
	}

	node, index, err := link.SearchData("hannah")
	if err != nil {
		t.Error(err)

	}
	fmt.Println("Node Found ", node, "at index", index)

	node, index, err = link.SearchData("dylan")
	if err == nil {
		t.Error("Expecting Error Here. Search Data not in list")
	}
	fmt.Println("Node Found ", node, "at index", index)

	link.PushFront("Roxy Music")
	link.ReIndex("chloe", 2)
	//link.Head.Right.Right.Data == index 2
	//link.Tail.Left.Data = Chloe's initial position
	if "chloe" != link.Head.Right.Right.Data {
		t.Error("Reindex Failed")
	} else if "chloe" == link.Tail.Left.Data {
		t.Error("Reindex Delete Failed")
	}
	node, index, err = link.SearchData("Roxy Music")
	if err != nil {
		t.Error(err)
	}
}

func TestPushFront(t *testing.T) {
	link := &LinkedList[string]{}
	firststring := "Test"
	link.PushFront(firststring)

	if link.Head.Data != firststring && link.Tail.Data != firststring {
		t.Error("Error Pushing Front")
	}
	secondstring := "Test2"
	link.PushFront(secondstring)
	if link.Head.Data != secondstring && link.Tail.Data != firststring {
		t.Error("Error Pushing Front 2")
	}
	thirdstring := "Test3"
	link.PushFront(thirdstring)
	if link.Head.Data != thirdstring && link.Head.Right.Data != secondstring && link.Tail.Left.Data != secondstring && link.Tail.Data != firststring {
		t.Error("Error Pushing Front 3")
	}
}
func TestPushBack(t *testing.T) {
	link := &LinkedList[string]{}
	firststring := "Test"
	link.PushBack(firststring)

	if link.Head.Data != firststring && link.Tail.Data != firststring {
		t.Error("Error Pushing Front")
	}
	secondstring := "Test2"
	link.PushBack(secondstring)
	if link.Head.Data != firststring && link.Tail.Data != secondstring {
		t.Error("Error Pushing Front 2")
	}
	thirdstring := "Test3"
	link.PushBack(thirdstring)
	if link.Head.Data != firststring && link.Head.Right.Data != secondstring && link.Tail.Left.Data != secondstring && link.Tail.Data != thirdstring {
		t.Error("Error Pushing Front 3")
	}
}

func TestInsertIndex(t *testing.T) {
	link := &LinkedList[int]{}
	_, err := link.InsertIndex(123, 5)
	//Error Expected
	if err == nil {
		t.Error("Index Out of Range Error Expected")
	}
	//Add Initial Index. Insert Index requires non zero list to work

	link.PushFront(1)

	node, err := link.InsertIndex(123, 0)
	if err != nil {
		t.Error("Index should be inserted")
	}
	if node.Data != link.Head.Data {
		t.Error("Index should be tail")
	}
	//Cannot add integers to end of list, expected. Use PushBack for that
	node, err = link.InsertIndex(51, 2)
	//Error Expected
	if err == nil {
		t.Error("Index Out of Range Error Expected")
	}
	//Inserting into middle
	node, err = link.InsertIndex(51, 1)
	if err != nil {
		t.Error("Index should be working")
	}
	//Compare Order
	//Expected Order
	numbers := []int{123, 51, 1}
	node = link.Head
	for _, n := range numbers {
		if node.Data != n {
			t.Error("Numbers should match")
		}
		node = node.Right
	}
	if link.Tail.Data != numbers[2] {
		t.Error("End of Linked List Expected")
	}
}
func TestDeleteDataNode(t *testing.T) {
	values := []int{4, 7, 1, 30, 9, 20, 0, 89, 50}
	link := &LinkedList[int]{}
	for _, v := range values {
		link.PushBack(v)
	}
	err := link.DeleteData(1)
	if err != nil {
		t.Error("Delete Data Error:", err)
	}
	if link.Head.Right.Right.Data != 30 {
		t.Error("Delete Data Not Properly entered")
	}
	err = link.DeleteData(100)
	//Expected Error
	if err == nil {
		t.Error("Data Not Found Error Expected")
	}
	//Remaining List {4, 7, 30, 9, 20, 0, 89, 50}
	//Delete error
	node := Node[int]{
		Data:  7,
		Right: nil,
		Left:  nil,
	}
	err = link.DeleteNode(node)
	if err == nil {
		t.Error("Expected Error")
	}
	//Delete Middle
	node = *link.Head.Right
	err = link.DeleteNode(node)
	if err != nil {
		t.Error(err)
	}
	t.Log("Delete Middle Test Complete")
	//Remaining List {4, 30, 9, 20, 0, 89, 50}
	//Delete Tail
	node = *link.Tail
	err = link.DeleteNode(node)
	t.Log("Delete Tail Test Complete")
	//Remaining List {4, 30, 9, 20, 0, 89}
	//Delete Head
	node = *link.Head
	err = link.DeleteNode(node)
	t.Log("Delete Head Test Complete")
	numbers := []int{30, 9, 20, 0, 89}
	node = *link.Head

	for _, n := range numbers {
		if node.Data != n {
			t.Error("Matching Failed")
		}
		if node.Right != nil {
			node = *node.Right
		}
	}
	t.Log("Compare Remaining Test Complete")
	link.DeleteData(30)
	link.DeleteData(9)
	link.DeleteData(20)
	link.DeleteData(0)
	node = *link.Head
	//Final Delete Test. Clear Linked List
	err = link.DeleteNode(node)
	if err != nil {
		t.Error("Error Deleting Node")
	} else if link.Size != 0 {
		t.Error("Linked List Sizing Error")
	}
}
