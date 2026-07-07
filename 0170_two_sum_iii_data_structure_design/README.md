# 0170 — Two Sum III - Data structure design

> LeetCode #170 · Difficulty: Easy (Premium) 🔒
> **Categories:** Array, Hash Table, Two Pointers, Design, Data Stream

---

## Problem Statement

Design a data structure that accepts a stream of integers and checks if it has a pair of integers that sum up to a particular value.

Implement the `TwoSum` class:

- `TwoSum()` Initializes the `TwoSum` object, with an empty array initially.
- `void add(int number)` Adds `number` to the data structure.
- `boolean find(int value)` Returns `true` if there exists any pair of numbers whose sum is equal to `value`, otherwise, it returns `false`.

**Example 1:**

```
Input
["TwoSum", "add", "add", "add", "find", "find"]
[[], [1], [3], [5], [4], [7]]
Output
[null, null, null, null, true, false]

Explanation
TwoSum twoSum = new TwoSum();
twoSum.add(1);   // [] --> [1]
twoSum.add(3);   // [1] --> [1,3]
twoSum.add(5);   // [1,3] --> [1,3,5]
twoSum.find(4);  // 1 + 3 = 4, return true
twoSum.find(7);  // No two integers sum up to 7, return false
```

**Constraints:**

- `-10^5 <= number <= 10^5`
- `-2^31 <= value <= 2^31 - 1`
- At most `10^4` calls will be made to `add` and `find`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| LinkedIn   | ★★★★★ Very High  | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Design / Data Structures** — the real question is choosing where to spend work: at `add` (ingestion) or `find` (query), driven by the expected call mix → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Hash Map** — a value → count map gives O(1) `add` and an O(distinct) complement scan on `find`; counts (not a plain set) are what make the self-pair case (`x + x = value`) decidable → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers** — a sorted store answers `find` with the #167 converging scan → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Binary Search** — keeps the Approach 2 store sorted via lower-bound insertion → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (unsorted slice) | add O(1) · find O(n²) | O(n) | Baseline only; find collapses once n grows |
| 2 | Sorted Slice + Two Pointers | add O(n) · find O(n) | O(n) | When you also need ordered data (min/max/range queries) |
| 3 | Hash Map of Counts (Optimal Balance) | add O(1) · find O(n) | O(n) | Default answer — best for add-heavy or mixed workloads |
| 4 | Precomputed Pair Sums | add O(n) · find O(1) | O(n²) | find ≫ add and memory is cheap |

---

## Approach 1 — Brute Force (Unsorted Slice, Check All Pairs)

### Intuition

Do nothing clever at ingestion: append every number to a slice. `find` then owes the full debt — test every pair `(i, j)` with `i < j`. Distinct *indices* (not distinct values) is exactly the problem's pairing rule, so duplicates are handled for free.

### Algorithm

1. `add(number)`: append to the slice.
2. `find(value)`: for every `i < j`, if `nums[i] + nums[j] == value` return `true`.
3. Return `false` if no pair matches.

### Complexity

- **Time:** `add` O(1) amortized (append); `find` O(n²) — all ~n²/2 pairs in the worst case (up to ~5·10⁷ checks with 10⁴ adds).
- **Space:** O(n) — the slice of added numbers.

### Code

```go
type TwoSumBrute struct {
	nums []int // every added number, in arrival order
}

// NewTwoSumBrute initializes the object with an empty container.
func NewTwoSumBrute() *TwoSumBrute {
	return &TwoSumBrute{nums: []int{}}
}

// Add stores number in the data structure. Time O(1).
func (t *TwoSumBrute) Add(number int) {
	t.nums = append(t.nums, number) // just remember it
}

// Find reports whether any pair of stored numbers sums to value. Time O(n^2).
func (t *TwoSumBrute) Find(value int) bool {
	for i := 0; i < len(t.nums)-1; i++ {
		for j := i + 1; j < len(t.nums); j++ {
			// Two distinct stored elements (indices differ) forming the sum.
			if t.nums[i]+t.nums[j] == value {
				return true
			}
		}
	}
	return false // no pair matched
}
```

### Dry Run

Example 1: `add(1), add(3), add(5), find(4), find(7)`.

| Step | Operation | nums | (i, j) pairs checked | Result |
|------|-----------|------|----------------------|--------|
| 1 | add(1) | [1] | — | — |
| 2 | add(3) | [1,3] | — | — |
| 3 | add(5) | [1,3,5] | — | — |
| 4 | find(4) | [1,3,5] | (0,1): 1+3 = 4 ✓ | `true` |
| 5 | find(7) | [1,3,5] | (0,1): 4 ✗ · (0,2): 6 ✗ · (1,2): 8 ✗ | `false` |

Output `[null, null, null, null, true, false]` ✔

---

## Approach 2 — Sorted Slice + Two Pointers

### Intuition

Import the lesson of #167: on sorted data, a pair-sum query needs only one converging two-pointer sweep. So keep the store sorted at all times — `add` binary-searches the insertion point (lower bound) and splices the number in, paying O(n) for the shift; `find` becomes the O(n) pointer scan with no extra memory.

### Algorithm

1. `add(number)`:
   1. `i = lowerBound(nums, number)` via binary search.
   2. Grow the slice by one, shift `nums[i:]` right, write `number` at `i`.
2. `find(value)`:
   1. `left = 0`, `right = n-1`.
   2. While `left < right`: sum equal → `true`; sum < value → `left++`; sum > value → `right--`.
   3. Return `false` when the pointers meet.

### Complexity

- **Time:** `add` O(n) — O(log n) search + O(n) element shift; `find` O(n) — single sweep.
- **Space:** O(n) — the sorted slice itself.

### Code

```go
type TwoSumSorted struct {
	nums []int // all added numbers, maintained in non-decreasing order
}

// NewTwoSumSorted initializes the object with an empty container.
func NewTwoSumSorted() *TwoSumSorted {
	return &TwoSumSorted{nums: []int{}}
}

// Add inserts number keeping the slice sorted. Time O(n).
func (t *TwoSumSorted) Add(number int) {
	// sort.SearchInts = lower bound: first index whose element >= number.
	i := sort.SearchInts(t.nums, number)
	t.nums = append(t.nums, 0)     // grow by one slot
	copy(t.nums[i+1:], t.nums[i:]) // shift the tail right to open index i
	t.nums[i] = number             // drop the new number into place
}

// Find reports whether any pair sums to value via two pointers. Time O(n).
func (t *TwoSumSorted) Find(value int) bool {
	left, right := 0, len(t.nums)-1
	for left < right {
		sum := t.nums[left] + t.nums[right]
		switch {
		case sum == value:
			return true // found a pair of distinct elements
		case sum < value:
			left++ // need a bigger sum → advance the small end
		default:
			right-- // need a smaller sum → retreat the large end
		}
	}
	return false // pointers met without hitting the target
}
```

### Dry Run

Example 1: `add(1), add(3), add(5), find(4), find(7)`.

| Step | Operation | nums (sorted) | left | right | sum | Action |
|------|-----------|---------------|------|-------|-----|--------|
| 1 | add(1) | [1] | — | — | — | insert at 0 |
| 2 | add(3) | [1,3] | — | — | — | insert at 1 |
| 3 | add(5) | [1,3,5] | — | — | — | insert at 2 |
| 4 | find(4) | [1,3,5] | 0 | 2 | 1+5 = 6 > 4 | right-- → 1 |
| 5 | find(4) cont. | [1,3,5] | 0 | 1 | 1+3 = 4 == 4 | return `true` |
| 6 | find(7) | [1,3,5] | 0 | 2 | 1+5 = 6 < 7 | left++ → 1 |
| 7 | find(7) cont. | [1,3,5] | 1 | 2 | 3+5 = 8 > 7 | right-- → 1; pointers meet → `false` |

Output `[null, null, null, null, true, false]` ✔

---

## Approach 3 — Hash Map of Counts (Optimal Balance)

### Intuition

Keep a frequency map `value → times added`. `add` is a single increment. `find(value)` walks the **distinct** values: each `x` needs the complement `value - x` in the map — and if the complement is `x` itself, `x` must have been added at least twice (a pair means two separate additions, not one element used twice). Counting distinct values also means a million copies of the same number cost one map entry. This add-O(1)/find-O(n) profile is the standard accepted solution.

### Algorithm

1. `add(number)`: `counts[number]++`.
2. `find(value)`: for each key `x` in `counts`:
   1. `need = value - x`.
   2. If `need == x` and `counts[x] >= 2` → `true`.
   3. If `need != x` and `counts[need] > 0` → `true`.
3. No key succeeds → `false`.

### Complexity

- **Time:** `add` O(1) average (one map increment); `find` O(n) where n = distinct values stored (each key does O(1) lookups).
- **Space:** O(n) — one map entry per distinct value, regardless of duplicates.

### Code

```go
type TwoSumHashMap struct {
	counts map[int]int // value → how many times it was added
}

// NewTwoSumHashMap initializes the object with an empty container.
func NewTwoSumHashMap() *TwoSumHashMap {
	return &TwoSumHashMap{counts: map[int]int{}}
}

// Add records one more occurrence of number. Time O(1).
func (t *TwoSumHashMap) Add(number int) {
	t.counts[number]++
}

// Find reports whether two stored occurrences sum to value. Time O(n).
func (t *TwoSumHashMap) Find(value int) bool {
	for x := range t.counts {
		need := value - x // complement that would complete the pair
		if need == x {
			// Pairing x with itself needs two separate additions of x.
			if t.counts[x] >= 2 {
				return true
			}
		} else if t.counts[need] > 0 {
			// Distinct complement present at least once.
			return true
		}
	}
	return false // no value has a usable complement
}
```

### Dry Run

Example 1: `add(1), add(3), add(5), find(4), find(7)`.

| Step | Operation | counts | Key x examined | need = value−x | Check | Result |
|------|-----------|--------|----------------|----------------|-------|--------|
| 1 | add(1) | {1:1} | — | — | — | — |
| 2 | add(3) | {1:1, 3:1} | — | — | — | — |
| 3 | add(5) | {1:1, 3:1, 5:1} | — | — | — | — |
| 4 | find(4) | {1:1, 3:1, 5:1} | x = 1 | 3 | need ≠ x; counts[3] = 1 > 0 ✓ | `true` |
| 5 | find(7) | {1:1, 3:1, 5:1} | x = 1 | 6 | counts[6] = 0 ✗ | continue |
| 6 | find(7) cont. | — | x = 3 | 4 | counts[4] = 0 ✗ | continue |
| 7 | find(7) cont. | — | x = 5 | 2 | counts[2] = 0 ✗ | keys exhausted → `false` |

Output `[null, null, null, null, true, false]` ✔ (map iteration order varies in Go; any order reaches the same answers.)

---

## Approach 4 — Precomputed Pair Sums (Fast Find)

### Intuition

Invert the trade-off: if the workload is `find`-dominated (e.g., rare inserts, constant querying), precompute the answers. Every `add(number)` pairs the newcomer with each already-stored element and records those sums in a set; `find` is then a single O(1) membership test. Memory is the casualty: up to one set entry per pair — O(n²).

### Algorithm

1. `add(number)`:
   1. For every stored `x`, insert `x + number` into the `sums` set.
   2. Append `number` to the stored list (duplicates kept — they create genuine pairs like `0+0`).
2. `find(value)`: return `sums[value]`.

### Complexity

- **Time:** `add` O(n) — pairs the new number with all n existing ones; `find` O(1) — one hash lookup.
- **Space:** O(n²) — the set can hold every pairwise sum (~5·10⁷ for 10⁴ adds; usually far fewer since `number` spans only 2·10⁵+1 values, bounding distinct sums at ~4·10⁵).

### Code

```go
type TwoSumPairSums struct {
	nums []int        // all added numbers (duplicates kept: they form real pairs)
	sums map[int]bool // every sum achievable by two distinct additions
}

// NewTwoSumPairSums initializes the object with an empty container.
func NewTwoSumPairSums() *TwoSumPairSums {
	return &TwoSumPairSums{nums: []int{}, sums: map[int]bool{}}
}

// Add records number and all new pair sums it creates. Time O(n).
func (t *TwoSumPairSums) Add(number int) {
	// Every existing element pairs with the newcomer exactly once.
	for _, x := range t.nums {
		t.sums[x+number] = true
	}
	t.nums = append(t.nums, number)
}

// Find reports whether value is an achievable pair sum. Time O(1).
func (t *TwoSumPairSums) Find(value int) bool {
	return t.sums[value] // set membership answers the query directly
}
```

### Dry Run

Example 1: `add(1), add(3), add(5), find(4), find(7)`.

| Step | Operation | New sums recorded | nums | sums | Result |
|------|-----------|-------------------|------|------|--------|
| 1 | add(1) | none (store empty) | [1] | {} | — |
| 2 | add(3) | 1+3 = 4 | [1,3] | {4} | — |
| 3 | add(5) | 1+5 = 6, 3+5 = 8 | [1,3,5] | {4,6,8} | — |
| 4 | find(4) | — | [1,3,5] | {4,6,8} | 4 ∈ sums → `true` |
| 5 | find(7) | — | [1,3,5] | {4,6,8} | 7 ∉ sums → `false` |

Output `[null, null, null, null, true, false]` ✔

---

## Key Takeaways

- **Design problems are workload questions.** State the add:find call ratio before choosing: add-heavy → count map (O(1)/O(n)); find-heavy → precomputed sums (O(n)/O(1)); balanced with ordering needs → sorted store. Saying this trade-off out loud *is* the interview answer.
- **Counts, not a set.** The classic bug is `find(0)` returning `true` after a single `add(0)`. A pair requires two separate additions, so store frequencies and demand `counts[x] >= 2` when the complement equals the element itself.
- **Iterate distinct values, not raw additions** — duplicate-heavy streams then cost nothing extra on `find`.
- Bounded key range (`|number| ≤ 10^5`) means pair sums span only ~4·10^5 distinct values — a reason the "O(n²) space" of Approach 4 is less scary in practice, and a hint that arrays could replace hash maps.

---

## Related Problems

- LeetCode #1 — Two Sum (one-shot array version: hash map of complements)
- LeetCode #167 — Two Sum II - Input Array Is Sorted (the two-pointer engine used by Approach 2)
- LeetCode #653 — Two Sum IV - Input is a BST (same pair query over a tree)
- LeetCode #346 — Moving Average from Data Stream (stream + query design pattern)
- LeetCode #155 — Min Stack (design: precompute at write time to make reads O(1))
