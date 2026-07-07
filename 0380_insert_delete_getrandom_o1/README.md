# 0380 — Insert Delete GetRandom O(1)

> LeetCode #380 · Difficulty: Medium
> **Categories:** Design, Hash Table, Array, Math, Randomization

---

## Problem Statement

Implement the `RandomizedSet` class:

- `RandomizedSet()` Initializes the `RandomizedSet` object.
- `bool insert(int val)` Inserts an item `val` into the set if not present. Returns `true` if the item was not present, `false` otherwise.
- `bool remove(int val)` Removes an item `val` from the set if present. Returns `true` if the item was present, `false` otherwise.
- `int getRandom()` Returns a random element from the current set of elements (it's guaranteed that at least one element exists when this method is called). Each element must have the **same probability** of being returned.

You must implement the functions of the class such that each function works in **average** `O(1)` time complexity.

**Example 1:**

```
Input
["RandomizedSet", "insert", "remove", "insert", "getRandom", "remove", "insert", "getRandom"]
[[], [1], [2], [2], [], [1], [2], []]
Output
[null, true, false, true, 2, true, false, 2]

Explanation
RandomizedSet randomizedSet = new RandomizedSet();
randomizedSet.insert(1); // Inserts 1 to the set. Returns true as 1 was inserted successfully.
randomizedSet.remove(2); // Returns false as 2 does not exist in the set.
randomizedSet.insert(2); // Inserts 2 to the set, returns true. Set now contains [1,2].
randomizedSet.getRandom(); // getRandom() should return either 1 or 2 randomly.
randomizedSet.remove(1); // Removes 1 from the set, returns true. Set now contains [2].
randomizedSet.insert(2); // 2 was already in the set, so return false.
randomizedSet.getRandom(); // Since 2 is the only number in the set, getRandom() will always return 2.
```

**Constraints:**

- `-2^31 <= val <= 2^31 - 1`
- At most `2 * 10^5` calls will be made to `insert`, `remove`, and `getRandom`.
- There will be **at least one** element in the data structure when `getRandom` is called.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design Data Structures** — combine two structures so every operation hits O(1) → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Hash Map** — value→index map gives O(1) membership and locate → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Array** — a contiguous slice gives O(1) uniform random indexing; swap-remove keeps it dense → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)

---

## Approaches Overview

| # | Approach | insert | remove | getRandom | Space | When to use |
|---|----------|--------|--------|-----------|-------|-------------|
| 1 | Slice Only (Brute Force) | O(n) | O(n) | O(1) | O(n) | Baseline; scans for membership |
| 2 | Hash Map + Slice / Swap-Remove (Optimal) | O(1) avg | O(1) avg | O(1) | O(n) | The intended answer; all ops O(1) |

---

## Approach 1 — Slice Only (Brute Force)

### Intuition

Keep everything in a slice. GetRandom is already O(1) (random index), but membership must be checked by a linear scan, making Insert and Remove O(n). Remove still uses the key swap-with-last trick so deletion itself is O(1) once the element is located.

### Algorithm

1. `insert`: scan for `val`; if found return false; else append and return true.
2. `remove`: scan for `val`; if not found return false; else swap it with the last element and shrink the slice; return true.
3. `getRandom`: return `vals[rand(len)]`.

### Complexity

- **Time:** Insert O(n), Remove O(n), GetRandom O(1).
- **Space:** O(n).

### Code

```go
type SliceSet struct {
	vals []int
	rng  *rand.Rand
}

func NewSliceSet() *SliceSet {
	return &SliceSet{rng: rand.New(rand.NewSource(1))}
}

func (s *SliceSet) Insert(val int) bool {
	for _, v := range s.vals {
		if v == val {
			return false
		}
	}
	s.vals = append(s.vals, val)
	return true
}

func (s *SliceSet) Remove(val int) bool {
	for i, v := range s.vals {
		if v == val {
			last := len(s.vals) - 1
			s.vals[i] = s.vals[last]
			s.vals = s.vals[:last]
			return true
		}
	}
	return false
}

func (s *SliceSet) GetRandom() int {
	return s.vals[s.rng.Intn(len(s.vals))]
}
```

### Dry Run

Example 1 operation sequence (values in `vals`):

| Op | scan result | return | vals after |
|----|-------------|--------|-----------|
| insert(1) | absent | true | [1] |
| remove(2) | absent | false | [1] |
| insert(2) | absent | true | [1,2] |
| getRandom() | rand index | 1 or 2 | [1,2] |
| remove(1) | at i=0, swap last | true | [2] |
| insert(2) | present | false | [2] |
| getRandom() | only element | 2 | [2] |

Returns `[null, true, false, true, (1|2), true, false, 2]` ✔

---

## Approach 2 — Hash Map + Slice / Swap-Remove (Optimal)

### Intuition

Two structures, each solving one requirement:

- a **slice** `vals` gives O(1) uniform random access for `getRandom`;
- a **map** `idx: value → position in vals` gives O(1) membership for `insert`/`remove`.

The subtlety is O(1) `remove` without leaving a hole in the slice: overwrite the target slot with the **last** element, update that element's index in the map, then pop the tail. Because the set is unordered, reordering is harmless.

### Algorithm

1. `insert(val)`: if `val` is in `idx` return false; record `idx[val] = len(vals)`, append `val`, return true.
2. `remove(val)`: if `val` not in `idx` return false; let `i = idx[val]`; move the last element into slot `i` and set its `idx` to `i`; pop the tail; `delete(idx, val)`; return true.
3. `getRandom`: return `vals[rand(len)]`.

### Complexity

- **Time:** Insert O(1) avg, Remove O(1) avg, GetRandom O(1).
- **Space:** O(n) — slice plus map.

### Code

```go
type RandomizedSet struct {
	vals []int
	idx  map[int]int
	rng  *rand.Rand
}

func NewRandomizedSet() *RandomizedSet {
	return &RandomizedSet{
		idx: make(map[int]int),
		rng: rand.New(rand.NewSource(1)),
	}
}

func (s *RandomizedSet) Insert(val int) bool {
	if _, ok := s.idx[val]; ok {
		return false
	}
	s.idx[val] = len(s.vals)
	s.vals = append(s.vals, val)
	return true
}

func (s *RandomizedSet) Remove(val int) bool {
	i, ok := s.idx[val]
	if !ok {
		return false
	}
	last := len(s.vals) - 1
	lastVal := s.vals[last]
	s.vals[i] = lastVal
	s.idx[lastVal] = i
	s.vals = s.vals[:last]
	delete(s.idx, val)
	return true
}

func (s *RandomizedSet) GetRandom() int {
	return s.vals[s.rng.Intn(len(s.vals))]
}
```

### Dry Run

Example 1: track `vals` and `idx`.

| Op | Action | return | vals | idx |
|----|--------|--------|------|-----|
| insert(1) | append | true | [1] | {1:0} |
| remove(2) | 2 ∉ idx | false | [1] | {1:0} |
| insert(2) | append | true | [1,2] | {1:0, 2:1} |
| getRandom() | rand | 1 or 2 | [1,2] | {1:0, 2:1} |
| remove(1) | i=0; move vals[1]=2 into slot 0; idx[2]=0; pop | true | [2] | {2:0} |
| insert(2) | 2 ∈ idx | false | [2] | {2:0} |
| getRandom() | only element | 2 | [2] | {2:0} |

Returns `[null, true, false, true, (1|2), true, false, 2]` ✔

*(In `main.go` the `getRandom` outputs are checked to be currently-present elements — and the final one to equal 2 — rather than fixed literals, since the value is random.)*

---

## Key Takeaways

- **Two structures, one for each need:** array for O(1) random indexing, hash map for O(1) membership/locate. Neither alone gives O(1) on all three ops.
- **Swap-with-last delete** is the core trick for O(1) removal from an array when order doesn't matter — remember to update the moved element's index in the map.
- Uniform `getRandom` requires the elements to sit in **contiguous** storage of exactly the current size — a plain hash set can't index randomly in O(1).
- This pattern generalises: #381 extends it to allow duplicates by mapping each value to a *set of indices*.

---

## Related Problems

- LeetCode #381 — Insert Delete GetRandom O(1) - Duplicates allowed (value → set of indices)
- LeetCode #146 — LRU Cache (combined map + doubly linked list design)
- LeetCode #379 — Design Phone Directory (pool + membership design)
- LeetCode #710 — Random Pick with Blacklist (remap + random index)
