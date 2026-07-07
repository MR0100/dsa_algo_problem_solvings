# 0284 — Peeking Iterator

> LeetCode #284 · Difficulty: Medium
> **Categories:** Design, Iterator, Array

---

## Problem Statement

Design an iterator that supports the `peek` operation on an existing iterator in addition to the `hasNext` and the `next` operations.

Implement the `PeekingIterator` class:

- `PeekingIterator(Iterator<int> nums)` Initializes the object with the given integer iterator `iterator`.
- `int next()` Returns the next element in the array and moves the pointer to the next element.
- `boolean hasNext()` Returns `true` if there are still elements in the array.
- `int peek()` Returns the next element in the array **without** moving the pointer.

**Note:** Each language may have a different implementation of the constructor and `Iterator`, but they all support the `int next()` and `boolean hasNext()` functions.

**Example 1:**
```
Input:
["PeekingIterator", "next", "peek", "next", "next", "hasNext"]
[[[1, 2, 3]], [], [], [], [], []]
Output:
[null, 1, 2, 2, 3, false]

Explanation:
PeekingIterator peekingIterator = new PeekingIterator([1, 2, 3]); // [1,2,3]
peekingIterator.next();    // return 1, the pointer moves to the next element [1,2,3].
peekingIterator.peek();    // return 2, the pointer does not move [1,2,3].
peekingIterator.next();    // return 2, the pointer moves to the next element [1,2,3].
peekingIterator.next();    // return 3, the pointer moves to the next element [1,2,3].
peekingIterator.hasNext(); // return False.
```

**Constraints:**
- `1 <= nums.length <= 1000`
- `1 <= nums[i] <= 1000`
- All the calls to `next` and `peek` are valid.
- At most `1000` calls will be made to `next`, `hasNext`, and `peek`.

**Follow-up:** How would you extend your design to be generic and work with all types, not just integer?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Iterator / adapter design** — wrap a forward-only iterator to add a new capability (`peek`) → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Lookahead buffering** — cache one element so `peek` can report the future without consuming it → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview
| # | Approach | Time (each op) | Space | When to use |
|---|----------|----------------|-------|-------------|
| 1 | Cache One Element (eager) | O(1) | O(1) | Standard; fetch one ahead at construction |
| 2 | Lazy Peek (fetch on demand) | O(1) | O(1) | Avoids constructor-time fetch; same guarantees |

---

## Approach 1 — Cache One Element (Optimal)

### Intuition
The base iterator only moves forward and can't be rewound, so to `peek` we pre-fetch one element into a buffer. `peek` reads the buffer; `next` returns the buffer then refills it from the base (if any); `hasNext` is true whenever a buffered element is held. One element of extra state, all operations O(1).

### Algorithm
1. Constructor: if `base.hasNext()`, `next = base.Next()`, `hasPeeked = true`.
2. `peek`: return `next`.
3. `next`: `v = next`; if `base.hasNext()` refill `next = base.Next()`, else `hasPeeked = false`; return `v`.
4. `hasNext`: return `hasPeeked`.

### Complexity
- **Time:** `peek`, `next`, `hasNext` all O(1).
- **Space:** O(1) — one cached value and a flag.

### Code
```go
type CachePeeking struct {
	base      *Iterator
	next      int
	hasPeeked bool
}

func NewCachePeeking(it *Iterator) *CachePeeking {
	p := &CachePeeking{base: it}
	if it.HasNext() {
		p.next = it.Next()
		p.hasPeeked = true
	}
	return p
}

func (p *CachePeeking) Peek() int { return p.next }

func (p *CachePeeking) Next() int {
	v := p.next
	if p.base.HasNext() {
		p.next = p.base.Next()
	} else {
		p.hasPeeked = false
	}
	return v
}

func (p *CachePeeking) HasNext() bool { return p.hasPeeked }
```

### Dry Run
Base iterator over `[1,2,3]`. Constructor pre-fetches `next=1`, `hasPeeked=true`, base position now at index 1.

| call | returns | buffered `next` (after) | base position (after) | hasPeeked |
|------|---------|-------------------------|------------------------|-----------|
| Next() | 1 | 2 | idx 2 | true |
| Peek() | 2 | 2 | idx 2 | true |
| Next() | 2 | 3 | idx 3 (end) | true |
| Next() | 3 | — | end | false |
| HasNext() | false | — | end | false |

Output: `[1, 2, 2, 3, false]`.

---

## Approach 2 — Lazy Peek

### Intuition
Rather than always buffering one ahead, buffer lazily. Keep a flag `peeked`. `peek`: if not yet peeked, consume one from the base, stash it, set `peeked=true`; return the stash. `next`: if `peeked`, return the stash and clear the flag; otherwise delegate to `base.Next()`. `hasNext`: true if we hold a stashed value **or** the base still has one. Avoids the constructor-time fetch while keeping O(1) ops.

### Algorithm
1. `peek`: if `!peeked` then `cache = base.Next()`, `peeked = true`; return `cache`.
2. `next`: if `peeked` then `peeked = false`, return `cache`; else return `base.Next()`.
3. `hasNext`: return `peeked || base.HasNext()`.

### Complexity
- **Time:** `peek`, `next`, `hasNext` all O(1).
- **Space:** O(1) — one cached value and a flag.

### Code
```go
type LazyPeeking struct {
	base   *Iterator
	cache  int
	peeked bool
}

func NewLazyPeeking(it *Iterator) *LazyPeeking {
	return &LazyPeeking{base: it}
}

func (p *LazyPeeking) Peek() int {
	if !p.peeked {
		p.cache = p.base.Next()
		p.peeked = true
	}
	return p.cache
}

func (p *LazyPeeking) Next() int {
	if p.peeked {
		p.peeked = false
		return p.cache
	}
	return p.base.Next()
}

func (p *LazyPeeking) HasNext() bool {
	return p.peeked || p.base.HasNext()
}
```

### Dry Run
Base iterator over `[1,2,3]`, nothing pre-fetched (`peeked=false`, base at idx 0).

| call | branch | returns | peeked (after) | base position (after) |
|------|--------|---------|-----------------|------------------------|
| Next() | !peeked → base.Next() | 1 | false | idx 1 |
| Peek() | !peeked → cache=base.Next()=2 | 2 | true | idx 2 |
| Next() | peeked → return cache | 2 | false | idx 2 |
| Next() | !peeked → base.Next() | 3 | false | idx 3 (end) |
| HasNext() | peeked(false) \|\| base.HasNext()(false) | false | false | end |

Output: `[1, 2, 2, 3, false]`.

---

## Key Takeaways
- **Peek = one-element lookahead buffer.** The whole trick is caching the upcoming value plus a "is it valid?" flag.
- **Eager vs lazy** buffering are behaviourally identical here; eager is simpler to state, lazy avoids touching the base until needed.
- For the generic follow-up, make the wrapper hold a value of the element type `T` and a boolean — no logic changes, just parameterise the type.

---

## Related Problems
- LeetCode #341 — Flatten Nested List Iterator (iterator design)
- LeetCode #251 — Flatten 2D Vector (iterator over 2D data)
- LeetCode #281 — Zigzag Iterator (wrapping/merging iterators)
- LeetCode #173 — Binary Search Tree Iterator (controlled advancement)
