# LRU Cache in Go

A fixed-size cache that throws out the least recently used item when it runs out of room. Built from scratch on my own doubly linked list and my own hash map. No standard library containers.

## Why an LRU

A cache holds data so you don't have to fetch it again. But memory is limited. When the cache is full, something has to go. An LRU picks the item nobody has touched in the longest time and evicts it. The bet is simple: if you haven't used it lately, you probably won't need it soon.

## How it works

It uses two data structures, because neither one can do the job alone.

- **Hash map:** fast lookup. Given a key, it finds the item in O(1). But it has no idea what order things were used in.
- **Doubly linked list:** tracks order. The front is the most recently used. The back is the least recently used. It can move or remove a node in O(1) if you already have a pointer to it. But finding a node by key means walking the whole list, which is O(n).

The LRU connects them. The hash map stores each key with a pointer to its node in the linked list. So lookup is fast, and reordering is fast.

**Get(key):** look up the node in the hash map. If it exists, move it to the front of the list and return its value.

**Put(key, value):** if the key exists, update the value and move the node to the front. If not, add a new node at the front and add it to the hash map. If the cache is now over capacity, evict the node at the back.

## The hard part

Building the LRU itself was easy. The hard part was making the linked list and hash map behave correctly on every edge case, like removing the only node, removing the head or tail, and resizing the map.

The trickiest bug was eviction. When I removed the last node from the list, I had to delete the same key from the hash map. But the node only held the value. It didn't know its own key. So the map entry was stuck there.

The fix: store the key inside each node. Now when a node gets evicted, it tells the hash map exactly which key to delete.

## Complexity

- Get: O(1)
- Put: O(1)
- Space: O(capacity)

## Run the tests

    go test ./...
    go test -cover ./...
