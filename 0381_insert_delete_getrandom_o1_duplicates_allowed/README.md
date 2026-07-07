# 0381 — Insert Delete GetRandom O(1) - Duplicates allowed

> LeetCode #381 · Difficulty: Hard
> **Categories:** Array, Hash Table, Math, Design, Randomized

---

## Problem Statement

`RandomizedCollection` is a data structure that contains a collection of numbers, possibly duplicates (i.e., a multiset). It should support inserting and removing specific elements and also reporting a random element.

Implement the `RandomizedCollection` class:

- `RandomizedCollection()` Initializes the empty `RandomizedCollection` object.
- `bool insert(int val)` Inserts an item `val` into the multiset, even if the item is already present. Returns `true` if the item is not present, `false` otherwise.
- `bool remove(int val)` Removes an item `val` from the multiset if present. Returns `true` if the item is present, `false` otherwise. Note that if `val` has multiple occurrences in the multiset, we only remove one of them.
- `int getRandom()` Returns a random element from the current multiset of elements. The probability of each element being returned is **linearly related to the number of the same values the multiset contains**.

You must implement the functions of the class such that each function works on **average** `O(1)` time complexity.

**Note:** The test cases are generated such that `getRandom` will only be called if there is **at least one** item in the `RandomizedCollection`.

**Example 1:**
```
Input
["RandomizedCollection", "insert", "insert", "insert", "getRandom", "remove", "getRandom"]
[[], [1], [1], [2], [], [1], []]
Output
[null, true, false, true, 2, true, 1]

Explanation
RandomizedCollection randomizedCollection = new RandomizedCollection();
randomizedCollection.insert(1);   // return true since the collection does not contain 1.
                                  // Inserts 1 into the collection.
randomizedCollection.insert(1);   // return false since the collection contains 1.
                                  // Inserts another 1 into the collection. Collection now contains [1,1].
randomizedCollection.insert(2);   // return true since the collection does not contain 2.
                                  // Inserts 2 into the collection. Collection now contains [1,1,2].
randomizedCollection.getRandom(); // getRandom should:
                                  // - return 1 with probability 2/3, or
                                  // - return 2 with probability 1/3.
randomizedCollection.remove(1);   // return true since the collection contains 1.
                                  // Removes 1 from the collection. Collection now contains [1,2].
randomizedCollection.getRandom(); // getRandom should return 1 or 2, both equally likely.
```

**Constraints:**
- `-2³¹ <= val <= 2³¹ - 1`
- At most `2 * 10⁵` calls **in total** will be made to `insert`, `remove`, and `getRandom`.
- There will be at least one element in the data structure when `getRandom` is called.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map (value → index set)** — map each value to the *set* of positions it occupies so duplicates are tracked and any occurrence can be located in O(1) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Swap-with-last O(1) deletion** — array trick: overwrite the removed slot with the last element and shrink, avoiding O(n) shifting → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Randomized / Uniform Sampling** — a contiguous backing slice lets `rand.Intn(n)` pick multiplicity-weighted uniformly → see [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md)

---

## Approaches Overview

Let n = current number of elements.

| # | Approach | Insert | Remove | GetRandom | Space | When to use |
|---|----------|--------|--------|-----------|-------|-------------|
| 1 | Brute Force (plain slice) | O(n) | O(n) | O(1) | O(n) | Baseline; too slow for the constraints |
| 2 | Slice + Index Map (Optimal) | O(1) avg | O(1) avg | O(1) | O(n) | The required solution |

---

## Approach 1 — Brute Force (Plain Slice)

### Intuition
The simplest multiset is a flat list. Insert appends. GetRandom indexes a random position (automatically weighted by multiplicity). Remove is the slow part: linear-scan for any occurrence and delete it by shifting. Correct, but Remove (and the presence check in Insert) is O(n).

### Algorithm
1. Insert: scan to learn if `val` already appears, then append; return `!existed`.
2. Remove: linear-scan for `val`; if found, delete that index (shift tail) and return `true`; else `false`.
3. GetRandom: return `elems[rand.Intn(len)]`.

### Complexity
- **Time:** Insert O(n) (presence scan), Remove O(n) (scan + shift), GetRandom O(1).
- **Space:** O(n) — the element list.

### Code
```go
func (c *SliceCollection) Insert(val int) bool {
	existed := false
	for _, v := range c.elems {
		if v == val {
			existed = true // val already present → return false later
			break
		}
	}
	c.elems = append(c.elems, val) // always store the new copy
	return !existed
}

func (c *SliceCollection) Remove(val int) bool {
	for i, v := range c.elems {
		if v == val {
			c.elems = append(c.elems[:i], c.elems[i+1:]...) // delete index i
			return true
		}
	}
	return false // val not found
}

func (c *SliceCollection) GetRandom() int {
	return c.elems[c.rng.Intn(len(c.elems))]
}
```

### Dry Run
Operations `insert(1), insert(1), insert(2), getRandom, remove(1), getRandom`.

| Op | scan result | elems after | return |
|----|-------------|-------------|--------|
| insert(1) | not found | `[1]` | true |
| insert(1) | found | `[1,1]` | false |
| insert(2) | not found | `[1,1,2]` | true |
| getRandom | idx∈{0,1,2} | `[1,1,2]` | 1 (2/3) or 2 (1/3) |
| remove(1) | i=0 | `[1,2]` | true |
| getRandom | idx∈{0,1} | `[1,2]` | 1 or 2, 1/2 each |

---

## Approach 2 — Slice + Index Map (Optimal)

### Intuition
GetRandom needs a contiguous array; Remove needs O(1) deletion. Keep all elements in a slice `nums`, and keep `idx[val]` = the **set** of positions holding `val`. To remove one copy of `val`, take any of its positions, overwrite it with the **last** element of the slice, fix that moved element's index bookkeeping, then shrink. Swap-with-last makes deletion O(1); the per-value index sets keep everything consistent under duplicates.

### Algorithm
1. **Insert(val):** append `val` to `nums`; add index `len-1` to `idx[val]`; return `true` iff `idx[val]` now has size 1 (val was absent before).
2. **Remove(val):** if `idx[val]` empty → `false`. Pick some index `i` in `idx[val]`; let `last = len(nums)-1`, `lastVal = nums[last]`. Set `nums[i] = lastVal`; remove `i` from `idx[val]`; if `i != last`, remove `last` from `idx[lastVal]` and add `i` to it. Shrink `nums`. Return `true`.
3. **GetRandom:** `nums[rand.Intn(len(nums))]`.

### Complexity
- **Time:** Insert O(1) avg, Remove O(1) avg, GetRandom O(1) — hash ops and swap-with-last are O(1) amortized.
- **Space:** O(n) — the slice plus the index sets.

### Code
```go
func (c *IndexMapCollection) Insert(val int) bool {
	if c.idx[val] == nil {
		c.idx[val] = map[int]struct{}{} // first time we ever see val
	}
	c.nums = append(c.nums, val)           // store the new copy at the end
	c.idx[val][len(c.nums)-1] = struct{}{} // record its position
	return len(c.idx[val]) == 1            // size 1 ⇒ val was absent before
}

func (c *IndexMapCollection) Remove(val int) bool {
	positions := c.idx[val]
	if len(positions) == 0 {
		return false // no copy of val present
	}
	var i int
	for p := range positions { // grab an arbitrary index holding val
		i = p
		break
	}
	last := len(c.nums) - 1
	lastVal := c.nums[last]

	c.nums[i] = lastVal   // overwrite removed slot with the last element
	delete(c.idx[val], i) // i no longer holds val

	if i != last { // the moved element lands at a new index i
		delete(c.idx[lastVal], last)   // no longer at `last`
		c.idx[lastVal][i] = struct{}{} // now at i
	}

	c.nums = c.nums[:last] // drop the tail slot
	if len(c.idx[val]) == 0 {
		delete(c.idx, val) // keep the map clean
	}
	return true
}

func (c *IndexMapCollection) GetRandom() int {
	return c.nums[c.rng.Intn(len(c.nums))]
}
```

### Dry Run
Operations `insert(1), insert(1), insert(2), getRandom, remove(1), getRandom`.

| Op | nums after | idx after | return |
|----|-----------|-----------|--------|
| insert(1) | `[1]` | `{1:{0}}` | true (size 1) |
| insert(1) | `[1,1]` | `{1:{0,1}}` | false (size 2) |
| insert(2) | `[1,1,2]` | `{1:{0,1}, 2:{2}}` | true |
| getRandom | `[1,1,2]` | — | 1 (2/3) or 2 (1/3) |
| remove(1): i=1, last=2, lastVal=2 → nums[1]=2, drop `1` from idx[1], move 2 from idx 2→1 | `[1,2]` | `{1:{0}, 2:{1}}` | true |
| getRandom | `[1,2]` | — | 1 or 2, 1/2 each |

(The exact `i` chosen from the set is arbitrary; if `i=0` were picked, `nums[0]=2` then shrink → `[2,1]`, with idx updated symmetrically. Both are valid.)

---

## Key Takeaways
- The **O(1) Insert/Delete/GetRandom** pattern = contiguous slice for random indexing + hash map for O(1) location + **swap-with-last** for O(1) deletion.
- Duplicates upgrade the map from `val → index` (problem #380) to `val → set of indices`. Any occurrence can be removed; you must fix the moved element's entry.
- Careful with the **`i == last`** edge case: when the removed element *is* the last one, don't re-add it to a map entry you just deleted.
- `getRandom` is multiplicity-weighted for free because duplicates occupy multiple slice slots.

---

## Related Problems
- LeetCode #380 — Insert Delete GetRandom O(1) (no duplicates; the base version)
- LeetCode #146 — LRU Cache (hash map + linked structure for O(1) ops)
- LeetCode #895 — Maximum Frequency Stack (value → structure bookkeeping)
- LeetCode #705 — Design HashSet
