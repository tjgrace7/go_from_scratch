# go_from_scratch
 
This repo is where I learn backend fundamentals in Go by building them myself. No frameworks where I can avoid them. No libraries that do the hard part for me.
 
I already shipped a full product in Python and FastAPI. It ran in production with real pilot users. But frameworks hid a lot from me. This repo is where I build the parts they hide, so I know what I am trusting when I use them.
 
## What is in here
 
| Project | Folder | What I learned |
|---|---|---|
| HTTP server on raw TCP | `cmd/rawtcp` | How bytes on a socket become an HTTP request |
| HTTP server on net/http | `cmd/nethttp` | Which layer the standard library replaces and which it leaves to me |
| Hash map | `hashmaps` | Hashing, collisions and resizing |
| Binary search tree | `binarytree` | Recursive insert, search and delete |
| Doubly linked list | `LinkedList` | Pointer updates that break quietly when you get them wrong |
| LRU cache | `LRUCache` | How a hash map and a list work together for O(1) get and put |
| toyRedis | `toyRedis` | Locks, race conditions, expiry and crash recovery |
 
Most folders have their own README with the full write-up and the bugs I hit along the way.
 
## The projects
 
### HTTP server, two ways
 
First I built an HTTP server on raw TCP. I wrote my own buffer loop to read the full request. The hard part was finding where the headers end and the body starts. A blank line splits them. After that line, `Content-Length` tells me how many body bytes are left.
 
Both versions check a Bearer token against an API key on every request.
 
Then I rebuilt it with `net/http`. It went from about 300 lines to about 100. The library took over the socket, the read loop and the parsing. It did not take over auth. That part is still mine.
 
Full write-up: [cmd/README.md](cmd/README.md)
 
### Hash map
 
A hash map built without Go's built-in `map`. It hashes keys with FNV-1a, which I wrote by hand. I also wrote djb2 and djb2a to compare.
 
Collisions go into chains. Each bucket points to the next entry by index, and -1 marks the end of a chain. I first used zero for that. But zero is a real index, so lookups looped forever. That bug taught me not to let a real value also mean "empty."
 
Resizing is called on purpose, not on every insert. The caller picks the load factor that triggers it and how much bigger the new map should be.
 
Full write-up: [hashmaps/README.md](hashmaps/README.md)
 
### Binary search tree
 
A binary search tree of integers. It supports insert, search, delete, in-order traversal, and finding a node's successor and predecessor. Delete removes a leaf directly. For any other node, it copies in the next value in order (the successor) or the one before it (the predecessor), then deletes that node lower in the tree.
 
Full write-up: [binarytree/README.md](binarytree/README.md)
 
### Linked list and LRU cache
 
The LRU cache holds a fixed number of items. When it is full, it drops the one used least recently. Get and put both run in O(1) time.
 
It uses two parts I built myself. The hash map finds a node fast. The doubly linked list keeps nodes in order of use, so moving or removing one is fast too.
 
The hardest bug was eviction. The node held its value but not its key. So when I dropped a node from the list, I could not remove it from the map. The fix was to store the key in each node.
 
Full write-ups: [LRUCache/README.md](LRUCache/README.md) and [LinkedList/README.md](LinkedList/README.md)
 
### toyRedis
 
A small key-value store in the style of Redis. Clients connect over TCP on port 6379 and send plain text commands:
 
```
SET name Tyler
SET name Tyler EX 10
GET name
DELETE name
```
 
- **Expiry.** `EX 10` makes a key expire in 10 seconds. GET checks expiry and deletes a stale key in one locked step, so a SET can't slip in between and get wiped out.
- **Persistence.** Every write goes to an append-only log. On restart, the server replays the log and then sweeps out keys that expired while it was down.
- **Concurrency.** A mutex guards the map, so many clients can connect at once without corrupting data. A race test runs 100 trials with the race detector on, and all pass.
The benchmark had a result I did not expect. Running 10 operations across goroutines was about 9.6 times slower than running them in order, 4,722 ns against 491 ns. Each operation is so cheap that the cost of goroutines and locks is bigger than the work. That is the same reason real Redis runs on one thread.
 
Full write-up: [toyRedis/README.md](toyRedis/README.md)
 
## Run it
 
You need Go 1.27 or later.
 
The HTTP servers read an API key from a `.env` file in the repo root:
 
```
API_KEY=your-key-here
```
 
Then:
 
```
go run ./cmd/rawtcp      # HTTP server on raw TCP, port 8080
go run ./cmd/nethttp     # HTTP server on net/http, port 8080
go run ./toyRedis        # key-value store, port 6379
```
 
Talk to toyRedis with netcat:
 
```
nc localhost 6379
```
 
## Tests
 
```
go test -cover ./...
go test -race ./toyRedis
```
 
| Package | Coverage |
|---|---|
| hashmaps | 90.1% |
| LinkedList | 88.9% |
| LRUCache | 82.1% |
| binarytree | 80.0% |
 
## What is next
 
- An outside code review from the Go community, with fixes applied in commits.
- CI with GitHub Actions to run the tests on every push.