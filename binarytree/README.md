# Binary Search Tree
 
A binary search tree of integers, built from scratch in Go.
 
## What it does
 
- **insert:** adds a value in the right spot. Smaller values go left, larger values go right. A value that is already in the tree is ignored.
- **search:** returns true if a value is in the tree.
- **delete:** removes a value and keeps the tree in order.
- **display:** prints every value in order, smallest to largest. This is an in-order traversal: left side, then the node, then the right side.
- **successor and predecessor:** find the replacement value when deleting a node.
## How I learned it
 
The insert, search and traversal code was easy to follow. To make sure I understood it, I drew out a large tree by hand and walked each function through it. Once I could see the tree, the recursion made sense.
 
## The hard part: delete
 
Deleting a leaf is easy. You just remove it. Deleting a node with children is harder, because something has to take its place without breaking the order of the tree.
 
It took me a minute to see the answer. Go to the node's right side, then go as far left as you can. That node is the successor. It is the smallest value on the right side. So it is larger than the node you are deleting and everything on its left, and smaller than everything else on its right. That makes it the perfect replacement.
 
The code copies the successor's value into the node being deleted. Then it deletes the successor from the right side, which is always an easier case.
 
If the node has no right side, the code does the same thing in reverse. It goes left, then as far right as it can, to find the predecessor. That is the largest value on the left side.
 
## Complexity
 
| Operation | Balanced tree | Worst case |
|---|---|---|
| insert | O(log n) | O(n) |
| search | O(log n) | O(n) |
| delete | O(log n) | O(n) |
| display | O(n) | O(n) |
 
The worst case happens when values go in already sorted. The tree turns into a straight line, which is just a slow linked list. This tree does not balance itself. Self-balancing trees like AVL and red-black trees fix that by rotating nodes after inserts and deletes.
 
## Run the tests
 
```
go test -cover ./binarytree
```
 
Coverage: 80.0%