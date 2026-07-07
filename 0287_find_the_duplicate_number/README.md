# 0287 — Find the Duplicate Number

> LeetCode #287 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Binary Search, Bit Manipulation

---

## Problem Statement

Given an array of integers `nums` containing `n + 1` integers where each integer is in the range `[1, n]` inclusive.

There is only **one repeated number** in `nums`, return *this repeated number*.

You must solve the problem **without** modifying the array `nums` and using only constant extra space.

**Example 1:**

```
Input: nums = [1,3,4,2,2]
Output: 2
```

**Example 2:**

```
Input: nums = [3,1,3,4,2]
Output: 3
```

**Example 3:**

```
Input: nums = [3,3,3,3,3]
Output: 3
```

**Constraints:**

- `1 <= n <= 10^5`
- `nums.length == n + 1`
- `1 <= nums[i] <= n`
- All the integers in `nums` appear only once except for precisely one integer which appears two or more times.

**Follow up:**

- How can we prove that at least one duplicate number must exist in `nums`?
- Can you solve the problem in linear runtime complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Floyd's Cycle Detection (Tortoise & Hare)** — reading `nums` as a functional graph `i → nums[i]` creates a cycle whose entrance is the duplicate → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Binary Search on the Answer** — search the value range `[1, n]` using a pigeonhole counting predicate → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Hashing** — a set records first sightings for the O(n)-space baseline → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — bring duplicates adjacent (on a copy, to respect the no-modify rule) → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force nested scan | O(n²) | O(1) | Tiny inputs; no extra space but too slow |
| 2 | Hash Set | O(n) | O(n) | Fast and trivial when extra space is allowed |
| 3 | Sort copy + adjacent scan | O(n log n) | O(n) | When a copy is acceptable; conceptually simple |
| 4 | Binary Search + count | O(n log n) | O(1) | O(1) space, no array modification; pigeonhole trick |
| 5 | Floyd's Cycle Detection | O(n) | O(1) | The intended optimal: O(n) time, O(1) space, no modification |

---

## Approach 1 — Brute Force (Nested Scan)

### Intuition

The most literal reading: a duplicate is one value living at two different indices. Compare all pairs; the value shared by a matching pair is the answer.

### Algorithm

1. For each index `i`, scan all `j > i`.
2. If `nums[i] == nums[j]`, return `nums[i]`.

### Complexity

- **Time:** O(n²) — every pair compared.
- **Space:** O(1) — no extra structure; array untouched.

### Code

```go
func bruteForce(nums []int) int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] { // same value at two indices → duplicate
				return nums[i]
			}
		}
	}
	return -1 // unreachable given the problem guarantees a duplicate exists
}
```

### Dry Run

Example 1: `nums = [1,3,4,2,2]`.

| i | nums[i] | j scanned | match? |
|---|---------|-----------|--------|
| 0 | 1 | 1,2,3,4 → 3,4,2,2 | no |
| 1 | 3 | 2,3,4 → 4,2,2 | no |
| 2 | 4 | 3,4 → 2,2 | no |
| 3 | 2 | 4 → 2 | **yes** → return 2 |

Result: `2` ✔

---

## Approach 2 — Hash Set

### Intuition

Walk once; the first value seen twice is the duplicate. A set gives O(1) membership.

### Algorithm

1. Keep a set `seen`.
2. For each value: if in `seen`, return it; otherwise add it.

### Complexity

- **Time:** O(n) — single pass.
- **Space:** O(n) — the set (violates the O(1) follow-up but is simplest).

### Code

```go
func hashSet(nums []int) int {
	seen := make(map[int]struct{}, len(nums)) // set of already-observed values
	for _, v := range nums {
		if _, ok := seen[v]; ok { // second sighting → duplicate
			return v
		}
		seen[v] = struct{}{} // record first sighting
	}
	return -1
}
```

### Dry Run

Example 1: `nums = [1,3,4,2,2]`.

| step | v | seen before? | seen after |
|------|---|--------------|------------|
| 1 | 1 | no | {1} |
| 2 | 3 | no | {1,3} |
| 3 | 4 | no | {1,3,4} |
| 4 | 2 | no | {1,3,4,2} |
| 5 | 2 | **yes** → return 2 | — |

Result: `2` ✔

---

## Approach 3 — Sort Then Adjacent Compare

### Intuition

After sorting, duplicate values are neighbours. Sort a **copy** so the original array is not modified (the constraint forbids in-place mutation).

### Algorithm

1. Copy `nums`, sort the copy.
2. Scan adjacent pairs; equal neighbours reveal the duplicate.

### Complexity

- **Time:** O(n log n) — the sort dominates.
- **Space:** O(n) — the copy (in-place sort would violate the no-modify rule).

### Code

```go
func sortScan(nums []int) int {
	c := make([]int, len(nums))
	copy(c, nums)  // copy so the original array is left untouched
	sort.Ints(c)   // duplicates become adjacent
	for i := 1; i < len(c); i++ {
		if c[i] == c[i-1] { // neighbours equal → duplicate
			return c[i]
		}
	}
	return -1
}
```

### Dry Run

Example 1: `nums = [1,3,4,2,2]` → sorted copy `[1,2,2,3,4]`.

| i | c[i-1] | c[i] | equal? |
|---|--------|------|--------|
| 1 | 1 | 2 | no |
| 2 | 2 | 2 | **yes** → return 2 |

Result: `2` ✔

---

## Approach 4 — Binary Search on Value Range (Pigeonhole)

### Intuition

Binary-search over **values**, not indices. For a candidate `mid`, count how many array elements are ≤ `mid`. Without a duplicate in `[1..mid]`, exactly `mid` elements would be ≤ `mid`. If the count exceeds `mid`, the pigeonhole principle forces the duplicate into `[1..mid]`; otherwise it lives in `[mid+1..n]`. Converge on the smallest value whose count overshoots.

### Algorithm

1. `lo = 1`, `hi = n` (where `n = len(nums) - 1`).
2. `mid = (lo+hi)/2`; `count = #{v : v ≤ mid}`.
3. If `count > mid`, the duplicate is in `[lo, mid]` → `hi = mid`; else `lo = mid + 1`.
4. Loop until `lo == hi`, the duplicate.

### Complexity

- **Time:** O(n log n) — `log n` search steps, each an O(n) count.
- **Space:** O(1) — array never modified, constant extra space.

### Code

```go
func binarySearchCount(nums []int) int {
	lo, hi := 1, len(nums)-1 // value range [1, n]
	for lo < hi {
		mid := lo + (hi-lo)/2
		count := 0
		for _, v := range nums { // how many values are ≤ mid?
			if v <= mid {
				count++
			}
		}
		if count > mid { // too many small values → duplicate ≤ mid
			hi = mid
		} else { // duplicate is on the high side
			lo = mid + 1
		}
	}
	return lo // lo == hi is the duplicate
}
```

### Dry Run

Example 1: `nums = [1,3,4,2,2]`, `n = 4`.

| lo | hi | mid | count(≤mid) | count > mid? | new range |
|----|----|-----|-------------|--------------|-----------|
| 1 | 4 | 2 | {1,2,2} = 3 | 3 > 2 yes | hi = 2 |
| 1 | 2 | 1 | {1} = 1 | 1 > 1 no | lo = 2 |
| 2 | 2 | — | — | loop ends | return 2 |

Result: `2` ✔

---

## Approach 5 — Floyd's Cycle Detection (Optimal)

### Intuition

Interpret the array as a linked list: from index `i`, follow to index `nums[i]`. Because every value is in `[1, n]` and there are `n+1` slots, this "next pointer" walk starting at index `0` must eventually loop — and the node where two chains merge is exactly the repeated value (two indices point to the same value node). Floyd's tortoise-and-hare finds a meeting point inside the cycle; restarting one pointer at the head and advancing both one step at a time makes them meet at the cycle's **entrance**, which is the duplicate.

### Algorithm

1. `slow = nums[0]`, `fast = nums[nums[0]]`.
2. Advance `slow` by 1 and `fast` by 2 until they meet.
3. Reset `slow = 0`; advance both one step until equal.
4. That value is the duplicate.

### Complexity

- **Time:** O(n) — linear pointer moves.
- **Space:** O(1) — two integer pointers; array never modified.

### Code

```go
func floydCycle(nums []int) int {
	// Phase 1: find an intersection point inside the cycle.
	slow, fast := nums[0], nums[nums[0]]
	for slow != fast {
		slow = nums[slow]       // one step
		fast = nums[nums[fast]] // two steps
	}
	// Phase 2: find the cycle entrance = duplicate value.
	slow = 0
	for slow != fast {
		slow = nums[slow] // both now move one step at a time
		fast = nums[fast]
	}
	return slow // entrance node == duplicated value
}
```

### Dry Run

Example 1: `nums = [1,3,4,2,2]`. Pointers hold *values* (which double as next indices).

Phase 1 (slow +1, fast +2):

| step | slow | fast | meet? |
|------|------|------|-------|
| init | nums[0]=1 | nums[nums[0]]=nums[1]=3 | no |
| 1 | nums[1]=3 | nums[nums[3]]=nums[2]=4 | no |
| 2 | nums[3]=2 | nums[nums[4]]=nums[2]=4 | no |
| 3 | nums[2]=4 | nums[nums[4]]=nums[2]=4 | **yes** (both 4) |

Phase 2 (reset slow = 0, both +1):

| step | slow | fast | equal? |
|------|------|------|--------|
| init | 0 | 4 | no |
| 1 | nums[0]=1 | nums[4]=2 | no |
| 2 | nums[1]=3 | nums[2]=4 | no |
| 3 | nums[3]=2 | nums[4]=2 | **yes** → return 2 |

Result: `2` ✔

---

## Key Takeaways

- **An array with values in `[1,n]` and length `n+1` is secretly a linked list with a cycle.** The `i → nums[i]` mapping guarantees a repeated node, and Floyd finds its entrance in O(n) time / O(1) space without touching the array — the intended solution.
- **Binary-search the answer, not the index.** When a monotone counting predicate ("how many elements ≤ x") separates the answer space, binary search over values gives O(n log n) with O(1) space.
- **Pigeonhole** proves a duplicate must exist: `n+1` items in `n` value-buckets.
- Know the trade-off ladder: O(n²)/O(1) → O(n)/O(n) hash → O(n log n)/O(1) binary search → O(n)/O(1) Floyd. Interviewers often walk you up it.

---

## Related Problems

- LeetCode #142 — Linked List Cycle II (Floyd entrance, same math)
- LeetCode #141 — Linked List Cycle (Floyd detection)
- LeetCode #268 — Missing Number (index ↔ value bijection)
- LeetCode #41 — First Missing Positive (values as indices)
- LeetCode #645 — Set Mismatch (find the duplicate and the missing)
