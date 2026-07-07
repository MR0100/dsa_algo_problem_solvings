# Design Data Structures

> **Category:** Design / Composition of primitives
> **Typical complexity target:** O(1) — often *amortized* O(1) — per operation
> **Signature move:** compose two (or more) primitive structures so that each
> one covers the other's weakness.

---

## 1. What it is

"Design" problems don't ask you to compute an answer over a fixed input.
They ask you to **build a class** — a stateful object with a small API — where
**every operation must meet a per-operation complexity bound** and the
operations arrive in an arbitrary, adversarial order.

The problem gives you an interface, e.g.:

```
LRUCache(capacity int)
Get(key int) int
Put(key int, value int)
```

and the real question hiding underneath is always the same:

> *No single primitive data structure supports all of these operations at the
> required speed. Which **combination** of primitives does — and how do you
> keep the pieces in sync?*

That is the entire genre. A hash map gives O(1) lookup but no ordering. A
doubly linked list gives O(1) ordered insert/delete *at a known node* but O(n)
lookup. A heap gives O(log n) min/max but no random access. A dynamic array
gives O(1) index access and O(1) tail ops but O(n) middle deletion. Design
problems are solved by **welding structures together** so each query hits the
structure that answers it in O(1)/O(log n), and every mutation updates *all*
structures atomically.

### The standard "welds" (memorise these pairings)

| Requirement combination | Winning composition | Canonical problem |
|---|---|---|
| O(1) lookup **and** O(1) recency/order updates + eviction | hash map → nodes of a **doubly linked list** | LRU Cache (#146) |
| O(1) insert, delete, **and** getRandom | hash map (value → index) + **dynamic array** (swap-with-last delete) | Insert Delete GetRandom O(1) (#380) |
| Stack ops **and** O(1) getMin | main stack + **parallel min stack** (or pair-encoding) | Min Stack (#155) |
| Queue semantics using only stacks (or vice versa) | **two stacks**, lazy transfer | Queue via Stacks (#232), Stack via Queues (#225) |
| Prefix lookup / autocomplete / word dictionary | **Trie** (prefix tree), children map/array + `isEnd` flag | Implement Trie (#208), Word Dictionary (#211) |
| Running median of a stream | **two heaps** (max-heap low half, min-heap high half) | Find Median from Data Stream (#295) |
| Frequency-ordered eviction | hash map + **buckets of doubly linked lists keyed by frequency** | LFU Cache (#460) |
| Iterator with `peek` / flattening | wrap inner iterator, **cache one lookahead element** / explicit stack | Peeking Iterator (#284), Flatten Nested List (#341) |
| Timestamped versions of a key | hash map → **append-only sorted slice** + binary search | Time Based KV Store (#981) |
| O(1) counts of a sliding window | **ring buffer** or queue + running aggregates | Moving Average (#346), Hit Counter (#362) |

If you internalise this table, ~90 % of LeetCode design problems reduce to
"recognise the row, implement the weld carefully."

---

## 2. How to recognise a design problem (signals in the statement)

- **"Design a …" / "Implement the `X` class"** — the giveaway is literal: the
  problem statement is an API spec, not a function spec.
- **Input given as two parallel arrays** — `["LRUCache","put","get",...]` and
  `[[2],[1,1],[1],...]`. LeetCode serialises method-call sequences this way;
  seeing it means *design problem*.
- **Per-operation complexity constraints** — "each function must run in O(1)
  average time", "follow up: can you do both operations in O(1)?" A bound on
  *each call*, not on the whole algorithm, means you must pick/compose data
  structures, not design an algorithm.
- **"…in constant time" attached to an operation a single structure can't do
  in constant time** (getMin, getRandom, evict-least-recent). This is the
  explicit hint that composition is required.
- **Stream / online wording** — "numbers arrive one at a time", "support the
  following operations on a stream". You cannot re-scan history each call, so
  you must maintain incremental state.
- **Huge number of operations** (`up to 10^5 calls`) with small per-call work
  implied — an O(n)-per-call solution TLEs by construction.
- **"Amortized"** appearing anywhere — the setter is telling you a lazy /
  batched strategy (two-stack queue, rebuild-on-demand) is acceptable and
  probably intended.

**Interview meta-signal:** design questions test *engineering* more than
cleverness — invariant discipline, edge cases (empty, capacity 1, duplicate
keys, update-vs-insert), and clean APIs. Interviewers deliberately probe:
"what if `put` is called on an existing key?", "what happens at capacity 0?".

---

## 3. General templates in Go

### 3.1 The skeleton every design problem shares

```go
// The universal shape: a struct holding the composed structures,
// a constructor, and methods. In Go there is no class — a struct +
// pointer-receiver methods plays that role.

type MyStructure struct {
    // one field per primitive structure in the weld;
    // document the INVARIANT that ties them together.
}

func Constructor( /* params */ ) MyStructure {
    // allocate every internal structure — never leave a nil map!
    return MyStructure{ /* ... */ }
}

// Pointer receivers: methods mutate state, so *MyStructure, not MyStructure.
func (s *MyStructure) Op1( /* ... */ ) { /* keep ALL structures in sync */ }
func (s *MyStructure) Op2( /* ... */ ) int { /* ... */ return 0 }
```

Design-problem discipline in one sentence: **write down the invariant that
links your internal structures, and make every method restore it before
returning.**

### 3.2 Hash map + doubly linked list (the LRU weld)

```go
// Invariant: every key in `index` maps to exactly one node in the list,
// and the list is ordered most-recent (front) → least-recent (back).

type dllNode struct {
    key, val   int      // store key IN the node so eviction can delete from map
    prev, next *dllNode
}

type LRUCache struct {
    cap        int
    index      map[int]*dllNode // key → its node (the O(1) lookup half)
    head, tail *dllNode         // sentinel nodes (the O(1) reorder half)
}

func Constructor(capacity int) LRUCache {
    // Sentinels: dummy head/tail so insert/remove never branch on nil.
    h, t := &dllNode{}, &dllNode{}
    h.next, t.prev = t, h
    return LRUCache{cap: capacity, index: make(map[int]*dllNode), head: h, tail: t}
}

// remove unlinks a node from the list — O(1) because we HOLD the node.
func (c *LRUCache) remove(n *dllNode) {
    n.prev.next = n.next // bypass n going forward
    n.next.prev = n.prev // bypass n going backward
}

// pushFront inserts right after the head sentinel (most-recent position).
func (c *LRUCache) pushFront(n *dllNode) {
    n.next = c.head.next
    n.prev = c.head
    c.head.next.prev = n
    c.head.next = n
}

func (c *LRUCache) Get(key int) int {
    n, ok := c.index[key]
    if !ok {
        return -1 // miss
    }
    c.remove(n)     // a read IS a use →
    c.pushFront(n)  // move to most-recent position
    return n.val
}

func (c *LRUCache) Put(key, value int) {
    if n, ok := c.index[key]; ok {
        n.val = value  // UPDATE path: overwrite, refresh recency, done
        c.remove(n)
        c.pushFront(n)
        return
    }
    if len(c.index) == c.cap { // INSERT path at capacity: evict LRU first
        lru := c.tail.prev        // least-recent = node before tail sentinel
        c.remove(lru)             // O(1) list delete
        delete(c.index, lru.key)  // ← this is why nodes store their key
    }
    n := &dllNode{key: key, val: value}
    c.pushFront(n)
    c.index[key] = n
}
```

### 3.3 Hash map + array with swap-delete (the getRandom weld)

```go
// Invariant: pos[v] is the index of v inside vals; vals holds each value once.
type RandomizedSet struct {
    vals []int       // dense array → O(1) uniform random pick
    pos  map[int]int // value → index in vals → O(1) membership + locate
}

func (s *RandomizedSet) Remove(val int) bool {
    i, ok := s.pos[val]
    if !ok {
        return false
    }
    last := s.vals[len(s.vals)-1]
    s.vals[i] = last              // overwrite hole with LAST element…
    s.pos[last] = i               // …and fix that element's recorded index
    s.vals = s.vals[:len(s.vals)-1] // shrink; middle-delete became tail-delete
    delete(s.pos, val)
    return true
}
// Insert: append + record index. GetRandom: vals[rand.Intn(len(vals))].
```

### 3.4 Lazy / amortized pattern (queue from two stacks)

```go
// Invariant: queue order = out (top→bottom) followed by in (bottom→top).
type MyQueue struct{ in, out []int }

func (q *MyQueue) Push(x int) { q.in = append(q.in, x) } // always O(1)

func (q *MyQueue) Pop() int {
    if len(q.out) == 0 {           // refill ONLY when out is empty —
        for len(q.in) > 0 {        // each element is moved at most once
            n := len(q.in) - 1     // in its lifetime ⇒ amortized O(1)
            q.out = append(q.out, q.in[n])
            q.in = q.in[:n]
        }
    }
    n := len(q.out) - 1
    v := q.out[n]
    q.out = q.out[:n]
    return v
}
```

### 3.5 Trie node (prefix-lookup family)

```go
type TrieNode struct {
    children [26]*TrieNode // or map[byte]*TrieNode for sparse alphabets
    isEnd    bool          // marks "a word ends here" — NOT "no children"
}

func (t *TrieNode) Insert(word string) {
    cur := t
    for i := 0; i < len(word); i++ {
        c := word[i] - 'a'
        if cur.children[c] == nil {
            cur.children[c] = &TrieNode{} // create path lazily
        }
        cur = cur.children[c]
    }
    cur.isEnd = true // flag the terminal node
}
```

---

## 4. Worked example — LRU Cache, traced step by step

Sequence (LeetCode #146, Example 1), capacity 2:

```
put(1,1) put(2,2) get(1) put(3,3) get(2) put(4,4) get(1) get(3) get(4)
```

List shown front (most-recent) → back (least-recent); `H`/`T` are sentinels.

| # | Call | Path taken | List after (MRU→LRU) | Map keys | Returns |
|---|------|-----------|----------------------|----------|---------|
| 1 | `put(1,1)` | miss, size 0 < 2 → pushFront | `H (1,1) T` | {1} | — |
| 2 | `put(2,2)` | miss, size 1 < 2 → pushFront | `H (2,2) (1,1) T` | {1,2} | — |
| 3 | `get(1)` | hit → remove(1), pushFront(1) | `H (1,1) (2,2) T` | {1,2} | **1** |
| 4 | `put(3,3)` | miss, size == 2 → evict `tail.prev` = (2,2); delete key 2 from map; pushFront(3) | `H (3,3) (1,1) T` | {1,3} | — |
| 5 | `get(2)` | key 2 not in map | unchanged | {1,3} | **-1** |
| 6 | `put(4,4)` | miss, at capacity → evict (1,1) (it is LRU *because step 3 refreshed 1, then step 4 pushed 3 in front of it*); pushFront(4) | `H (4,4) (3,3) T` | {3,4} | **—** |
| 7 | `get(1)` | evicted at step 6 | unchanged | {3,4} | **-1** |
| 8 | `get(3)` | hit → move to front | `H (3,3) (4,4) T` | {3,4} | **3** |
| 9 | `get(4)` | hit → move to front | `H (4,4) (3,3) T` | {3,4} | **4** |

Expected output: `1, -1, -1, 3, 4` — matches.

**What the trace teaches:**
- Step 3 is the subtle one: a **read reorders**. If `Get` didn't move the node,
  step 4 would wrongly evict key 1 instead of key 2 and the whole tail of the
  trace would differ.
- Step 4 shows why nodes carry their key: eviction discovers the victim via
  `tail.prev` (a *node*), then must delete it from the *map* — impossible in
  O(1) unless the node knows its own key.
- Every row ends with the invariant intact: map keys == list nodes, list
  ordered by recency. That per-operation invariant check *is* the correctness
  proof.

---

## 5. Common pitfalls (and how to avoid them)

1. **Structures drifting out of sync.** Updating the list but forgetting
   `delete(map, key)` on eviction (or vice versa) is the #1 bug. *Fix:* state
   the invariant in a comment; in every method, pair each mutation of one
   structure with the matching mutation of the other — ideally inside small
   helpers (`remove`, `pushFront`) so it can't be forgotten.
2. **Missing the update-vs-insert branch in `Put`.** Putting an existing key
   must overwrite the value and refresh recency — **not** insert a duplicate
   node or trigger an eviction. Interviewers test this on purpose.
3. **Forgetting that `Get` mutates.** In LRU, a read changes recency order.
   Any "read-only" mental model of `Get` produces wrong evictions later.
4. **Value receivers in Go.** `func (c LRUCache) Put(...)` mutates a *copy*;
   state silently never changes. Always use pointer receivers (`*LRUCache`)
   for stateful methods, and note maps/slices inside a value receiver *partly*
   work (shared backing) — which makes this bug intermittent and nastier.
5. **Nil maps / unlinked sentinels.** `var m map[int]int` panics on write;
   sentinels must be linked to each other (`h.next = t; t.prev = h`) in the
   constructor. Initialise **everything** in `Constructor`.
6. **Skipping sentinel nodes.** Hand-rolling head/tail nil-checks quadruples
   the branch count in `remove`/`pushFront` and breeds off-by-one bugs at size
   0 and 1. Two dummy nodes erase every special case.
7. **Middle-deletion from a slice with `copy`/`append`.** That's O(n) and
   blows the O(1) contract. Use **swap-with-last** (§3.3) when order doesn't
   matter, or a linked list when it does. And remember swap-delete must fix
   the moved element's index in the map — forgetting `pos[last] = i` is the
   classic RandomizedSet bug (it even breaks when removing the last element
   itself only if you reorder the lines carelessly — do map-fix *before*
   shrink, delete the target key *last*).
8. **Eager instead of lazy transfer in two-stack queues.** Transferring on
   every `Push` *and* `Pop` degrades to O(n) per op. Transfer only when `out`
   is empty; each element moves once → amortized O(1).
9. **Trie: confusing "node exists" with "word exists".** `search("app")` must
   check `isEnd` at the final node; `startsWith("app")` must not. Merging the
   two checks fails half the test cases.
10. **Capacity edge cases.** capacity 1 (every insert at capacity evicts;
    get-then-put ordering matters), capacity 0 if the constraints allow it,
    and `Get` on an empty structure. Walk these by hand before submitting.
11. **Ignoring "amortized" vs "worst-case".** If the follow-up demands
    worst-case O(1) (e.g., Min Stack), a rebuild-on-demand trick is invalid —
    you need the parallel-structure design instead.
12. **Testing methods in isolation.** Design bugs are *sequence* bugs (the
    step-3→step-4 interaction above). Always run the full example call
    sequence in `main()`, printing after each call.

---

## 6. Problems in this repo

No problem in the currently written range (0001–0130) is a design problem —
the classic ones (LRU Cache #146, Min Stack #155, Implement Trie #208, Queue
via Stacks #232, Find Median from Data Stream #295, RandomizedSet #380, LFU
Cache #460) fall later in the numbering. Nearest touchpoint so far:

- [0001 — Two Sum](../0001_two_sum/README.md) — its Related Problems section
  points at Two Sum III (#170), the data-structure-design variant of Two Sum
  (design a class supporting `add` and `find`).

> Problems 0131–0400 are being written concurrently; a later pass will link
> the design problems from that range (starting with #146 LRU Cache and
> #155 Min Stack) here.
