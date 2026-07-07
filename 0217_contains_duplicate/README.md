# 0217 — Contains Duplicate

> LeetCode #217 · Difficulty: Easy
> **Categories:** Array, Hash Table, Sorting

---

## Problem Statement

Given an integer array `nums`, return `true` if any value appears **at least twice** in the array, and return `false` if every element is distinct.

**Example 1:**

```
Input: nums = [1,2,3,1]
Output: true
Explanation: The element 1 occurs at the indices 0 and 3.
```

**Example 2:**

```
Input: nums = [1,2,3,4]
Output: false
Explanation: All elements are distinct.
```

**Example 3:**

```
Input: nums = [1,1,1,3,3,4,3,2,4,2]
Output: true
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Set** — the optimal one-pass solution stores seen values in a set and reports the first collision → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — an alternate O(n log n) approach where duplicates become adjacent → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Tiny arrays only; illustrates the definition |
| 2 | Sorting | O(n log n) | O(n) (or O(1) in place) | When you cannot afford extra hash memory but can sort |
| 3 | Hash Set (Optimal) | O(n) | O(n) | The default answer — fastest, one pass |

---

## Approach 1 — Brute Force

### Intuition

A duplicate is, by definition, two different positions holding equal values. The most direct check compares every element with every element that comes after it. The first equal pair proves a duplicate exists.

### Algorithm

1. For each index `i` from `0` to `n-1`.
2.   For each index `j` from `i+1` to `n-1`: if `nums[i] == nums[j]`, return `true`.
3. If no pair matched, return `false`.

### Complexity

- **Time:** O(n²) — in the worst case (all distinct) every one of the ~n²/2 pairs is compared.
- **Space:** O(1) — no auxiliary structures.

### Code

```go
func bruteForce(nums []int) bool {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ { // only pairs with j > i, no self-compare
			if nums[i] == nums[j] { // found two equal values at different indices
				return true
			}
		}
	}
	return false // no matching pair anywhere
}
```

### Dry Run

Example 1: `nums = [1,2,3,1]`.

| i | j | nums[i] | nums[j] | equal? |
|---|---|---------|---------|--------|
| 0 | 1 | 1 | 2 | no |
| 0 | 2 | 1 | 3 | no |
| 0 | 3 | 1 | 1 | **yes → return true** |

Result: `true` ✔

---

## Approach 2 — Sorting

### Intuition

If some value appears more than once, sorting the array places all its copies contiguously. Then the global question "is any value repeated?" collapses to the local question "are any two neighbors equal?", answerable in one linear scan.

### Algorithm

1. Copy `nums` (to avoid mutating the caller's slice) and sort the copy.
2. Scan `i` from `1` to `n-1`: if `arr[i] == arr[i-1]`, return `true`.
3. Otherwise return `false`.

### Complexity

- **Time:** O(n log n) — the sort dominates; the adjacency scan is O(n).
- **Space:** O(n) for the defensive copy — O(1) extra if sorting the input in place is permitted.

### Code

```go
func sorting(nums []int) bool {
	arr := make([]int, len(nums)) // copy so the caller's slice is untouched
	copy(arr, nums)
	sort.Ints(arr) // equal values become adjacent after sorting
	for i := 1; i < len(arr); i++ {
		if arr[i] == arr[i-1] { // adjacent equal → duplicate
			return true
		}
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,1]` → sorted `arr = [1,1,2,3]`.

| i | arr[i-1] | arr[i] | equal? |
|---|----------|--------|--------|
| 1 | 1 | 1 | **yes → return true** |

Result: `true` ✔

---

## Approach 3 — Hash Set (Optimal)

### Intuition

Make a single pass while remembering everything seen so far. The first time the current value is already in the "seen" set, that value must have appeared earlier — a duplicate. Set membership is O(1) average, so the whole scan is linear.

### Algorithm

1. Create an empty set `seen`.
2. For each `v` in `nums`: if `v ∈ seen`, return `true`; else insert `v`.
3. If the loop completes, return `false`.

### Complexity

- **Time:** O(n) — one pass, O(1) average per lookup/insert.
- **Space:** O(n) — up to `n` distinct keys held in the set.

### Code

```go
func hashSet(nums []int) bool {
	seen := make(map[int]struct{}, len(nums)) // struct{} is a zero-byte "present" marker
	for _, v := range nums {
		if _, ok := seen[v]; ok { // v already recorded → second occurrence
			return true
		}
		seen[v] = struct{}{} // record v as seen
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,1]`.

| Step | v | seen before | v in seen? | action |
|------|---|-------------|------------|--------|
| 1 | 1 | {} | no | insert 1 → {1} |
| 2 | 2 | {1} | no | insert 2 → {1,2} |
| 3 | 3 | {1,2} | no | insert 3 → {1,2,3} |
| 4 | 1 | {1,2,3} | **yes** | **return true** |

Result: `true` ✔

---

## Key Takeaways

- **`map[T]struct{}` is Go's idiomatic set** — the zero-byte value type signals "membership only, no payload".
- **Trade space for time:** the hash set converts an O(n²) pairwise check into an O(n) pass by paying O(n) memory.
- **Sorting turns "any duplicate anywhere" into "adjacent duplicates"** — a recurring reduction (see also removing duplicates from sorted arrays).
- Return early on the first hit; there is no need to scan the rest once a duplicate is confirmed.

---

## Related Problems

- LeetCode #219 — Contains Duplicate II (duplicate within index distance `k`)
- LeetCode #220 — Contains Duplicate III (near-duplicate in value and index)
- LeetCode #26 — Remove Duplicates from Sorted Array (adjacency after sorting)
- LeetCode #287 — Find the Duplicate Number (exactly one duplicate, O(1) space)
- LeetCode #442 — Find All Duplicates in an Array
