Linked List
A custom doubly linked list, written from scratch. This was not a standalone requirement of the Data Structures From Scratch project. It was built as a supplement to the LRU cache, since an LRU cache needs O(1) insertion and removal at both ends. Of the three structures in that project, this was the second hardest to write.
What it does
PushFront: insert a node at the head of the list.
PushBack: insert a node at the tail of the list.
InsertIndex: insert a node at a specific index. InsertIndex cannot add a node to the end of the list, that case is handled by PushBack instead. InsertIndex works by inserting a node before an existing node, and the tail has no next node to insert before, so it does not have a valid target to anchor to.
Index traversal: walk the list by index to find the correct node.
Delete by Data: remove a node by matching its stored value.
Delete by Node: remove a node directly if you already hold a reference to it.
Reindex: recalculate index positions across the list.
Bug I hit and needed help finding
In PushFront, the if statement that checked whether head and tail were both nil (the empty list case) did not increase the list's size. That meant the very first node inserted into the list was never counted. Any later search by index was off by one from the start, because the size the code thought it had did not match the number of nodes actually in the list.
The fix was making sure the empty list branch increments size the same as every other insertion path does. Easy to miss because the branch looks like it is only setting head and tail, not touching the count, when in fact every insertion path needs to touch the count.
Why this matters for the LRU cache
The LRU cache needs to move a node to the front on access and evict from the back when it's full. Both of those are O(1) operations only if the linked list gives you O(1) push and O(1) delete when you already hold the node. Building this by hand instead of using a library forces you to actually reason about pointer updates, which is the part that breaks quietly if you get it wrong.
