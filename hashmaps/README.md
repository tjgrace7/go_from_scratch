Hash Map (Built From Scratch, No Standard Library Shortcuts)
A hash map implemented in Go with a custom hash function, separate chaining for collisions, and dynamic resizing. Built as part of Phase 1 of a backend engineering curriculum to understand how maps and dictionaries actually work under the hood.
Why build this
Python dictionaries and Go's built-in maps hide the mechanics. Building one by hand forces you to deal with hashing, collisions, memory layout, and resizing directly instead of trusting a library to handle it. This project is that exercise.
Collision handling: separate chaining
Collisions are handled with separate chaining (also called closed addressing). When two keys hash to the same bucket, the second key is appended to a chain at that bucket instead of being placed in a different slot.
The chain-termination bug
The first attempt at chaining checked for a zero value to know when a chain ended, or checked an Occupied bool on the next slot. This broke on the second chain. The last real index in a chain would read as zero, zero would also read as occupied, and the lookup would loop through the same slots forever.
The fix was to use -1 as the standard "end of chain" marker instead of zero or a bool. Zero is a valid index, so it can never safely mean "nothing here." -1 cannot be confused with a real index, so a single check against -1 reliably marks the end of a chain. This removed the infinite loop.
Lesson: don't overload a valid data value (zero) to also mean "empty." Use a value that is impossible to collide with real data.
The delete bug
Deleting a key from the middle of a chain means shifting every later entry up one index and clearing the final slot.
The first version of this had a bug. When shifting an entry up, the code copied the entry's key and value but read the nextIndex from the wrong source, the secondary list's next index, which pointed one link further down the chain than it should have. This caused two problems: some keys were skipped entirely during shifts, and the program sometimes crashed outright because the nextIndex it followed didn't always point to a real entry.
Lesson: when shifting linked structures, track pointers (or indices) explicitly for each node you touch. Don't assume the "next" value on a neighboring struct is the one you need. Copy structure and pointers separately and verify each one.
Resizing
Resizing is not automatic on every insert. It's triggered by a Resize() call, controlled by two caller-supplied values:
Load factor: the fill percentage that triggers a resize. A load factor of .75 means the map resizes once it's 75% full.
Size factor: the multiplier applied to the current size when resizing. A size factor of 10 makes the new map 10 times the size of the old one.
This keeps resize behavior explicit and tunable by the caller instead of hardcoded.

