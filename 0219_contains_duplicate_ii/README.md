# 0219 — Contains Duplicate II

> LeetCode #219 · Difficulty: Easy
> **Categories:** Array, Hash Table, Sliding Window

---

## Problem Statement

Given an integer array `nums` and an integer `k`, return `true` if there are two **distinct indices** `i` and `j` in the array such that `nums[i] == nums[j]` and `abs(i - j) <= k`.

**Example 1:**

```
Input: nums = [1,2,3,1], k = 3
Output: true
```

**Example 2:**

```
Input: nums = [1,0,1,1], k = 1
Output: true
```

**Example 3:**

```
Input: nums = [1,2,3,1,2,3], k = 2
Output: false
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`
- `0 <= k <= 10^5`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Palantir   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Hash Set** — track the last index (or the window contents) per value for O(1) lookups → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sliding Window** — the optimal approach keeps a fixed-size set of the last `k` values and slides it across the array → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n·k) | O(1) | Small `k`; no extra memory |
| 2 | Hash Map of Last Index | O(n) | O(n) | Simple one-pass; stores all values |
| 3 | Sliding-Window Hash Set (Optimal) | O(n) | O(min(n,k)) | Best space — set never exceeds `k` entries |

---

## Approach 1 — Brute Force

### Intuition

We want two equal values whose indices differ by at most `k`. For each index `i`, any partner `j` must lie within `k` positions ahead, so we only scan the next `k` entries instead of the whole tail. The first equal nearby pair answers the question.

### Algorithm

1. For each `i` from `0` to `n-1`.
2.   Scan `j` from `i+1` up to `min(i+k, n-1)`.
3.   If `nums[i] == nums[j]`, return `true`.
4. If nothing matched, return `false`.

### Complexity

- **Time:** O(n·k) — each `i` inspects at most `k` neighbours.
- **Space:** O(1) — no auxiliary structures.

### Code

```go
func bruteForce(nums []int, k int) bool {
	n := len(nums)
	for i := 0; i < n; i++ {
		// Only indices within k of i can satisfy |i-j| <= k.
		for j := i + 1; j <= i+k && j < n; j++ {
			if nums[i] == nums[j] { // equal values close enough in index
				return true
			}
		}
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,1], k = 3`.

| i | j range | pairs checked | match? |
|---|---------|---------------|--------|
| 0 | 1..3 | (1,2),(1,3),(1,1) | at j=3: nums[0]=1 == nums[3]=1 → **true** |

Result: `true` ✔

---

## Approach 2 — Hash Map of Last Index

### Intuition

For any value, the *closest* previous occurrence is simply its most recent one. So we only ever need the last index each value was seen at. On re-encountering a value, if the gap to its stored index is `<= k`, we found a valid pair; otherwise we overwrite the stored index with the current (nearer) one and continue.

### Algorithm

1. Keep a map `value → last index`.
2. For each `i`: if `nums[i]` is in the map and `i - lastIndex <= k`, return `true`.
3. Set `map[nums[i]] = i` (always overwrite with the newest index).
4. If the loop ends, return `false`.

### Complexity

- **Time:** O(n) — single pass, O(1) average per map operation.
- **Space:** O(n) — up to `n` distinct values are stored.

### Code

```go
func hashMapLastIndex(nums []int, k int) bool {
	lastIndex := make(map[int]int, len(nums)) // value → most recent index
	for i, v := range nums {
		if j, ok := lastIndex[v]; ok && i-j <= k { // seen before AND within k
			return true
		}
		lastIndex[v] = i // record/refresh the newest index of v
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,1], k = 3`.

| i | v | lastIndex before | seen & i-j<=k? | action |
|---|---|------------------|----------------|--------|
| 0 | 1 | {} | no | set {1:0} |
| 1 | 2 | {1:0} | no | set {1:0,2:1} |
| 2 | 3 | {1:0,2:1} | no | set {…,3:2} |
| 3 | 1 | {1:0,…} | j=0, 3-0=3 <= 3 → **yes** | **return true** |

Result: `true` ✔

---

## Approach 3 — Sliding-Window Hash Set (Optimal)

### Intuition

Instead of storing every value's index, keep only the values within the last `k` indices in a set — that window is exactly the set of candidate partners for the current element. If `nums[i]` is already in the window, a near-duplicate exists. After adding `nums[i]`, if the window would span more than `k` indices, evict the value that just slid out (the one at index `i-k`), keeping the set size bounded by `k`.

### Algorithm

1. Maintain a set `window` of the last `k` values.
2. For each `i`: if `nums[i] ∈ window`, return `true`.
3. Insert `nums[i]`; if `len(window) > k`, delete `nums[i-k]`.
4. If the loop ends, return `false`.

### Complexity

- **Time:** O(n) — one pass, O(1) average per set operation.
- **Space:** O(min(n, k)) — the window never holds more than `k` values.

### Code

```go
func slidingWindowSet(nums []int, k int) bool {
	window := make(map[int]struct{}) // values within the last k indices
	for i, v := range nums {
		if _, ok := window[v]; ok { // v already present in the k-window
			return true
		}
		window[v] = struct{}{} // add current value
		if len(window) > k {   // window grew beyond k values → evict oldest
			delete(window, nums[i-k]) // remove the value leaving the window
		}
	}
	return false
}
```

### Dry Run

Example 2: `nums = [1,0,1,1], k = 1`.

| i | v | window before | v in window? | after add | evict? | window after |
|---|---|---------------|--------------|-----------|--------|--------------|
| 0 | 1 | {} | no | {1} | size 1 ≤ 1, no | {1} |
| 1 | 0 | {1} | no | {1,0} | size 2 > 1 → delete nums[0]=1 | {0} |
| 2 | 1 | {0} | no | {0,1} | size 2 > 1 → delete nums[1]=0 | {1} |
| 3 | 1 | {1} | **yes** | — | — | **return true** |

Result: `true` ✔

---

## Key Takeaways

- **"Most recent index" beats "all indices"** — for a distance test you only ever need the closest prior occurrence, which is the last one seen.
- **Fixed-size sliding window = bounded memory.** Evicting the element at index `i-k` right after inserting keeps the set at ≤ `k` entries, the tightest space you can achieve here.
- **`|i-j| <= k` becomes a window of width `k`** — the same window idea scales to Contains Duplicate III where it combines with a bucketing trick on values.
- Watch the boundary `j <= i+k && j < n` in brute force so you don't index out of range.

---

## Related Problems

- LeetCode #217 — Contains Duplicate (no index constraint)
- LeetCode #220 — Contains Duplicate III (near-duplicate in both value and index)
- LeetCode #3 — Longest Substring Without Repeating Characters (sliding window + set)
- LeetCode #438 — Find All Anagrams in a String (fixed-size window)
