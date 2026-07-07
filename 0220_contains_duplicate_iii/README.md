# 0220 — Contains Duplicate III

> LeetCode #220 · Difficulty: Hard
> **Categories:** Array, Sliding Window, Sorting, Bucket Sort, Ordered Set

---

## Problem Statement

You are given an integer array `nums` and two integers `indexDiff` and `valueDiff`.

Find a pair of indices `(i, j)` such that:

- `i != j`,
- `abs(i - j) <= indexDiff`,
- `abs(nums[i] - nums[j]) <= valueDiff`, and

Return `true` if such a pair exists, or `false` otherwise.

**Example 1:**

```
Input: nums = [1,2,3,1], indexDiff = 3, valueDiff = 0
Output: true
Explanation: We can choose (i, j) = (0, 3).
We satisfy the three conditions:
i != j --> 0 != 3
abs(i - j) <= indexDiff --> abs(0 - 3) <= 3
abs(nums[i] - nums[j]) <= valueDiff --> abs(1 - 1) <= 0
```

**Example 2:**

```
Input: nums = [1,5,9,1,5,9], indexDiff = 2, valueDiff = 3
Output: false
Explanation: After trying all the possible pairs (i, j), we cannot satisfy the three conditions, so we return false.
```

**Constraints:**

- `2 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`
- `1 <= indexDiff <= nums.length`
- `0 <= valueDiff <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window** — the `abs(i-j) <= indexDiff` bound restricts every check to a moving window of the last `indexDiff` elements → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Ordered Set / Binary Search** — inside the window we binary-search a sorted structure for a neighbour within `valueDiff` → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Bucketing (Bucket Sort idea)** — the optimal O(n) approach maps values into buckets of width `valueDiff+1` so a same-bucket collision is an instant hit → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Hash Map** — buckets are stored as `bucketId → value` for O(1) access → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n·indexDiff) | O(1) | Small `indexDiff`; no extra memory |
| 2 | Sliding Window + Sorted Structure | O(n·log(indexDiff)) search (O(n·indexDiff) with a plain slice) | O(min(n,indexDiff)) | Clean "ordered set" formulation; near-optimal with a balanced BST |
| 3 | Bucketing (Optimal) | O(n) | O(min(n,indexDiff)) | The best asymptotics — constant work per element |

---

## Approach 1 — Brute Force

### Intuition

We want indices `i < j` with `j - i <= indexDiff` and `|nums[i] - nums[j]| <= valueDiff`. Directly test each `i` against the up-to-`indexDiff` elements immediately after it. The index bound caps the inner loop, and the first pair that also meets the value bound is our answer.

### Algorithm

1. For each `i`, scan `j` from `i+1` to `min(i+indexDiff, n-1)`.
2. If `abs(nums[i] - nums[j]) <= valueDiff`, return `true`.
3. If nothing qualifies, return `false`.

### Complexity

- **Time:** O(n·indexDiff) — each `i` inspects at most `indexDiff` neighbours.
- **Space:** O(1) — no auxiliary storage.

### Code

```go
func bruteForce(nums []int, indexDiff int, valueDiff int) bool {
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j <= i+indexDiff && j < n; j++ { // only indices within indexDiff
			if abs(nums[i]-nums[j]) <= valueDiff { // values close enough too
				return true
			}
		}
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,1], indexDiff = 3, valueDiff = 0`.

| i | j | nums[i] | nums[j] | \|diff\| | ≤ valueDiff(0)? |
|---|---|---------|---------|--------|-----------------|
| 0 | 1 | 1 | 2 | 1 | no |
| 0 | 2 | 1 | 3 | 2 | no |
| 0 | 3 | 1 | 1 | 0 | **yes → return true** |

Result: `true` ✔

---

## Approach 2 — Sliding Window + Sorted Structure

### Intuition

The index bound says only the last `indexDiff` values matter, so keep them in a window. Within the window we need a value in `[nums[i]-valueDiff, nums[i]+valueDiff]`. If the window is kept **sorted**, binary search finds the smallest element `>= nums[i]-valueDiff`; if that element is also `<= nums[i]+valueDiff`, it is a valid partner. This mirrors a balanced-BST / TreeSet `ceiling` query. As `i` advances, drop the value that leaves the index window.

### Algorithm

1. Maintain a sorted slice `window` of the last `indexDiff` values.
2. For each `i`: binary-search the first element `>= nums[i]-valueDiff`.
3.   If it exists and is `<= nums[i]+valueDiff`, return `true`.
4.   Insert `nums[i]` into `window`, preserving sorted order.
5.   If `i >= indexDiff`, remove `nums[i-indexDiff]` from `window`.
6. If the loop ends, return `false`.

### Complexity

- **Time:** O(n·log(indexDiff)) for the searches; with a plain slice the insert/delete shifts cost O(indexDiff), giving O(n·indexDiff) overall. A balanced BST / skip list keeps it at O(n·log(indexDiff)).
- **Space:** O(min(n, indexDiff)) — the window.

### Code

```go
func slidingWindowSorted(nums []int, indexDiff int, valueDiff int) bool {
	window := make([]int, 0) // sorted values of the last indexDiff indices
	for i, v := range nums {
		// Find first element >= v - valueDiff (the lowest acceptable neighbour).
		pos := sort.SearchInts(window, v-valueDiff)
		// If such an element exists and is also <= v + valueDiff, it is within range.
		if pos < len(window) && window[pos] <= v+valueDiff {
			return true
		}
		// Insert v into the sorted window at its correct position.
		ins := sort.SearchInts(window, v)                 // where v belongs
		window = append(window, 0)                        // grow by one
		copy(window[ins+1:], window[ins:])                // shift right to open a gap
		window[ins] = v                                   // place v
		// Evict the value that is now out of the index window.
		if i >= indexDiff {
			out := nums[i-indexDiff]                       // value leaving the window
			del := sort.SearchInts(window, out)           // find it (guaranteed present)
			window = append(window[:del], window[del+1:]...)
		}
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,1], indexDiff = 3, valueDiff = 0`.

| i | v | window before | search v−0=v | ceiling ≤ v+0? | insert / evict |
|---|---|---------------|--------------|----------------|----------------|
| 0 | 1 | [] | none | no | window=[1] |
| 1 | 2 | [1] | first ≥ 2 → none (pos=1) | no | window=[1,2] |
| 2 | 3 | [1,2] | first ≥ 3 → none | no | window=[1,2,3] |
| 3 | 1 | [1,2,3] | first ≥ 1 → window[0]=1, 1 ≤ 1 | **yes** | **return true** |

Result: `true` ✔

---

## Approach 3 — Bucketing (Optimal)

### Intuition

Slice the value axis into buckets each `valueDiff+1` wide, `bucket(v) = floor(v / (valueDiff+1))`. Any two values sharing a bucket differ by at most `valueDiff` — so a second value landing in an occupied bucket is an **immediate** yes. Values within `valueDiff` that don't share a bucket must lie in an **adjacent** bucket (`b-1` or `b+1`), and there we verify the difference explicitly. Because a same-bucket collision returns instantly, we only ever need to store one value per bucket. Restrict the live buckets to the current index window by deleting the bucket of the value that slides out. Negative values need a floor division so the bucket boundaries stay consistent.

### Algorithm

1. `width = valueDiff + 1`; `bucket(v) = floor(v / width)`.
2. For each `i` with value `v`, `b = bucket(v)`:
3.   If bucket `b` occupied → return `true`.
4.   If bucket `b-1` occupied and `|v - its value| <= valueDiff` → return `true`.
5.   If bucket `b+1` occupied and `|v - its value| <= valueDiff` → return `true`.
6.   Store `v` in bucket `b`.
7.   If `i >= indexDiff`, delete the bucket of `nums[i-indexDiff]`.
8. If the loop ends, return `false`.

### Complexity

- **Time:** O(n) — each index does O(1) bucket lookups/updates.
- **Space:** O(min(n, indexDiff)) — at most `indexDiff+1` buckets are alive at once.

### Code

```go
func bucketing(nums []int, indexDiff int, valueDiff int) bool {
	buckets := make(map[int]int) // bucket id → the single value stored there
	width := valueDiff + 1       // bucket width so one bucket spans valueDiff+1 values

	// getBucket computes a floor-division bucket id that also works for negatives.
	getBucket := func(v int) int {
		if v >= 0 {
			return v / width
		}
		return (v+1)/width - 1 // arithmetic floor division for negative v
	}

	for i, v := range nums {
		b := getBucket(v)
		if _, ok := buckets[b]; ok { // same bucket → guaranteed within valueDiff
			return true
		}
		// Adjacent buckets can still hold a value within valueDiff — verify.
		if x, ok := buckets[b-1]; ok && abs(v-x) <= valueDiff {
			return true
		}
		if x, ok := buckets[b+1]; ok && abs(v-x) <= valueDiff {
			return true
		}
		buckets[b] = v // store v (at most one value per bucket in the window)
		if i >= indexDiff {
			// The value at index i-indexDiff is leaving the window; drop its bucket.
			delete(buckets, getBucket(nums[i-indexDiff]))
		}
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1,2,3,1], indexDiff = 3, valueDiff = 0` → `width = 1`, so `bucket(v) = v`.

| i | v | b | bucket b occupied? | b−1 / b+1 within valueDiff? | action |
|---|---|---|--------------------|-----------------------------|--------|
| 0 | 1 | 1 | no | b0 empty, b2 empty | buckets={1:1} |
| 1 | 2 | 2 | no | b1=1, \|2−1\|=1 > 0 no; b3 empty | buckets={1:1,2:2} |
| 2 | 3 | 3 | no | b2=2, \|3−2\|=1 > 0 no; b4 empty | buckets={1:1,2:2,3:3} |
| 3 | 1 | 1 | **yes (holds 1)** | — | **return true** |

Result: `true` ✔ (values 1 and 1, indices 0 and 3, within index and value bounds).

---

## Key Takeaways

- **Two coupled constraints → window + ordered lookup.** The index bound becomes a sliding window; the value bound becomes a range/ceiling query inside it.
- **Bucket width `valueDiff+1` is the trick.** Same bucket ⇒ difference `≤ valueDiff` for free; only adjacent buckets need an explicit check. This is what drops the value query from `log` to `O(1)`.
- **One value per bucket is enough** because a same-bucket collision short-circuits immediately.
- **Floor division for negatives.** `v/width` truncates toward zero in Go, which breaks bucket boundaries for negative values; `(v+1)/width - 1` gives the true floor. This is exactly what the `[-3,3]` test guards.
- With a balanced BST / TreeSet (Approach 2) the value query is `O(log indexDiff)`; Go's standard library has no ordered set, so the slice version trades that for `O(indexDiff)` shifts — the bucketing approach sidesteps the issue entirely.

---

## Related Problems

- LeetCode #217 — Contains Duplicate (no constraints)
- LeetCode #219 — Contains Duplicate II (index bound only)
- LeetCode #2841 — Maximum Sum of Almost Unique Subarray (fixed window + hash)
- LeetCode #480 — Sliding Window Median (ordered structure over a window)
