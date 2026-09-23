toyRedis
An in-memory key-value store with a Redis-style text protocol, TTL expiry, crash-safe persistence, and benchmarked/tested concurrent access. Built from scratch in Go as part of a self-taught backend engineering curriculum.
What it does
toyRedis runs a TCP server that speaks a simple text protocol:
SET name Tyler
GET name
DELETE name
SET name Tyler EX 10
SET stores the second word as the key and the third as the value. Adding EX followed by an integer sets an expiration in seconds, SET name Tyler EX 10 expires the key name in 10 seconds.
GET returns the value for a key. It also checks expiration on every read: if the key has expired, it deletes it from the map as part of the same lookup instead of returning stale data.
DELETE removes a key from the map directly.
Persistence
Every write is appended to append.txt before it's acknowledged, so the store can rebuild its state after a crash or restart.
Writing the log required a custom string parser, not just splitting on spaces. Timestamps are written in RFC3339 format, which contains spaces of its own (2026-09-22T08:10:03-05:00 has no spaces, but the log entries needed a way to group a field that might). I added a parenthesis-delimited field: if a token starts with (, the parser holds off splitting until it finds the matching ). A missing closing parenthesis is treated as a malformed log entry and returns an error rather than silently misparsing. The parentheses are inserted by the program when writing, never typed by hand.
On startup, the server replays append.txt to restore state, then runs an expiry sweep to clean up any keys that expired while the process was down. The log stores the actual expiration timestamp, not the raw TTL the user typed, so restoring a key doesn't require re-running a countdown. A key logged with an expiration in the past is already correctly expired the moment it's loaded back in.
Concurrency
Correctness
Multiple clients can hit the store concurrently without corrupting data. The map is guarded by a mutex; every read and write to the map happens inside a lock.
SET does its prep work, timestamping and building the record to insert, before acquiring the lock, so that work happens while a goroutine is waiting its turn rather than while it's holding the lock. The lock itself only covers the actual map write.
GET's expiry check and delete needed to happen as a single atomic step, not two separate lock/unlock pairs. Checking expiry and removing the key as two separate operations opens a window where a concurrent SET on the same key can land in between, get silently deleted by the stale expiry check, and lose data. Both steps happen under one lock.
I also found and fixed a lock leak in GET: an early return inside a conditional branch was skipping the unlock on some code paths. A test that should have taken 10 seconds hung indefinitely instead. Running with go test -timeout 5s forced a goroutine stack dump, which showed a SET goroutine permanently blocked waiting on the lock, with no GET goroutine visible in the dump at all. A goroutine that's actually deadlocked shows up in the dump; one that's finished doesn't. That absence was the tell: GET had already returned without releasing the lock. The fix was defer mu.Unlock() immediately after acquiring the lock, so every return path releases it regardless of which branch it takes.
Verified with a test that runs the exact race, a SET racing a GET on a key at its expiry boundary, 10 times per run, across 10 separate go test invocations with the race detector on: 100 trials total, all passing.
Performance
A benchmark comparing the same 10 operations (4 SET, 4 GET, 2 DELETE) run sequentially versus concurrently across goroutines:
Mode	ns/op	vs. Sequential
Sequential (no locking)	~491.5	1x (baseline)
Concurrent, Mutex	~4,722	~9.6x slower
Concurrent, RWMutex	~4,872	~9.9x slower (statistically the same as Mutex)
Each number is an average across 3 runs (-count=3 -benchtime=3s).
Running the operations concurrently was slower than running them sequentially. Each operation is a single map read or write, on the order of 50ns of actual work. Spawning 10 goroutines and coordinating them through a sync.WaitGroup costs roughly 4,200ns of overhead, over 8x the work being parallelized. The management cost of waiting in line for the lock, plus goroutine scheduling, outweighed the work itself. This lines up with why real Redis is single-threaded: cheap per-operation work doesn't justify per-call concurrency overhead. A larger per-operation cost, or a larger dataset, could change this result; this benchmark reflects this workload at this scale.
RWMutex did not improve on plain Mutex, it was, if anything, marginally slower, though the difference is within run-to-run noise. The benchmark mix is half writes (SET/DELETE need the exclusive lock), which leaves limited room for concurrent readers to overlap under RWMutex, and RWMutex carries slightly more overhead per lock than a plain Mutex.
An earlier pass at this benchmark included fmt.Println calls inside the timed section, which added stdout's own internal lock on top of the store's mutex and produced an inflated, misleading number. Removing logging from the hot path was necessary to get a clean measurement.
Requirements checklist
 TCP server with a text protocol (SET, GET, DELETE)
 TTL/expiry on keys
 Append-only log for persistence across restarts
 Benchmark: ops/sec with 10 concurrent clients
 Multiple concurrent clients without data corruption
What I'd do differently
Given how cheap each operation is here, I'd reach for concurrency only if the per-operation cost were high enough to justify the coordination overhead, or if I were optimizing for something other than raw throughput on a single machine. Real key-value stores like Redis make this same call: single-threaded, no per-command locking, because the operations are cheap enough that concurrency costs more than it saves.
