# 0164 — Maximum Gap

> LeetCode #164 · Difficulty: Hard
> **Categories:** Array, Sorting, Bucket Sort, Radix Sort

---

## Problem Statement

Given an integer array `nums`, return the maximum difference between two successive elements in its sorted form. If the array contains less than two elements, return `0`.

You must write an algorithm that runs in **linear time** and uses **linear extra space**.

**Example 1:**
```
Input: nums = [3,6,9,1]
Output: 3
Explanation: The sorted form of the array is [1,3,6,9], either (3,6) or (6,9) has the maximum difference 3.
```

**Example 2:**
```
Input: nums = [10]
Output: 0
Explanation: The array contains less than 2 elements, therefore return 0.
```

**Constraints:**
- `1 <= nums.length <= 10^5`
- `0 <= nums[i] <= 10^9`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★☆☆ Medium    | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Microsoft | ★★☆☆☆ Low       | 2023          |
| Apple     | ★★☆☆☆ Low       | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting (non-comparison sorts)** — radix sort and bucket sort beat the Ω(n log n) comparison-sort lower bound because the values are bounded integers → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Math / Pigeonhole Principle** — n numbers create n−1 sorted gaps averaging (max−min)/(n−1), so the maximum gap is at least that; this bounds where the answer can live → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach                          | Time       | Space | When to use                                                      |
|---|-----------------------------------|------------|-------|-------------------------------------------------------------------|
| 1 | Brute Force (Successor Search)    | O(n²)      | O(1)  | Never in practice; shows the definition without sorting           |
| 2 | Sort and Scan                     | O(n log n) | O(n)  | Real-world default when the linear-time constraint is waived      |
| 3 | Radix Sort                        | O(n)       | O(n)  | Meets the bound by fully sorting bounded integers digit by digit  |
| 4 | Bucket Sort + Pigeonhole (Optimal)| O(n)       | O(n)  | The intended answer — linear time *without* fully sorting         |

---

## Approach 1 — Brute Force (Successor Search)

### Intuition
In the sorted form, the element that follows `v` is the *smallest value strictly greater than* `v` — its successor. The gap after `v` is `successor − v`. So for every element, find its successor with a full scan and keep the largest gap. No sorting at all, just the definition. Duplicates produce sorted gaps of 0, which can never beat a positive maximum, so treating equal values as "not a successor" is safe (an all-equal array correctly yields 0).

### Algorithm
1. If `n < 2`, return `0`.
2. For every element `v` in `nums`:
   1. Scan all elements `w`, keeping the smallest `w` with `w > v` (the successor).
   2. If a successor exists, update `maxGap = max(maxGap, successor − v)`.
3. Return `maxGap`.

### Complexity
- **Time:** O(n²) — each of the n elements triggers a full O(n) successor scan; ~10¹⁰ steps at n = 10⁵, far too slow.
- **Space:** O(1) — a few scalar variables.

### Code
```go
func bruteForce(nums []int) int {
	if len(nums) < 2 {
		return 0 // fewer than two elements → no successive pair exists
	}
	maxGap := 0
	for _, v := range nums { // v plays the role of "left element of a sorted pair"
		successor := -1 // smallest value strictly greater than v; -1 = none found
		for _, w := range nums {
			if w > v && (successor == -1 || w < successor) {
				successor = w // tighter successor candidate
			}
		}
		if successor != -1 && successor-v > maxGap { // v is not the maximum value
			maxGap = successor - v
		}
	}
	return maxGap
}
```

### Dry Run
`nums = [3,6,9,1]` (Example 1):

| v | values > v | successor | gap = successor − v | maxGap after |
|---|------------|-----------|---------------------|--------------|
| 3 | 6, 9       | 6         | 3                   | 3            |
| 6 | 9          | 9         | 3                   | 3            |
| 9 | none       | −1        | — (skip)            | 3            |
| 1 | 3, 6, 9    | 3         | 2                   | 3            |

Return **3** ✅ — exactly the adjacent gaps of the sorted form `[1,3,6,9]`.

---

## Approach 2 — Sort and Scan

### Intuition
The problem is *defined* on the sorted array, so the obvious solution is to actually sort and then take the largest adjacent difference in one pass. Correct and simple — but a comparison sort costs Θ(n log n), violating the required linear bound. This approach is the honest baseline that motivates the two linear approaches below.

### Algorithm
1. If `n < 2`, return `0`.
2. Copy `nums` (keep the caller's slice intact) and sort the copy ascending.
3. Scan `i = 1 … n−1`, tracking `maxGap = max(maxGap, arr[i] − arr[i−1])`.
4. Return `maxGap`.

### Complexity
- **Time:** O(n log n) — dominated by the comparison sort; the scan is O(n).
- **Space:** O(n) — the defensive copy (sorting in place would make it O(log n) for sort internals).

### Code
```go
func sortAndScan(nums []int) int {
	if len(nums) < 2 {
		return 0 // no pair to compare
	}
	arr := make([]int, len(nums)) // copy so the caller's slice stays untouched
	copy(arr, nums)
	sort.Ints(arr) // produce the sorted form
	maxGap := 0
	for i := 1; i < len(arr); i++ {
		if gap := arr[i] - arr[i-1]; gap > maxGap { // adjacent sorted difference
			maxGap = gap
		}
	}
	return maxGap
}
```

### Dry Run
`nums = [3,6,9,1]` (Example 1). After sorting: `arr = [1,3,6,9]`:

| i | arr[i−1] | arr[i] | gap | maxGap after |
|---|----------|--------|-----|--------------|
| 1 | 1        | 3      | 2   | 2            |
| 2 | 3        | 6      | 3   | 3            |
| 3 | 6        | 9      | 3   | 3            |

Return **3** ✅.

---

## Approach 3 — Radix Sort

### Intuition
The linear-time requirement rules out *comparison* sorting, not sorting altogether. Because the values are bounded non-negative integers (`0 ≤ nums[i] ≤ 10⁹ < 2³²`), LSD radix sort sorts them in a constant number of stable counting-sort passes — 4 passes of 8-bit digits cover all 32 bits. Total cost O(4·(n + 256)) = O(n). Then the answer is the same adjacent-difference scan as Approach 2.

### Algorithm
1. If `n < 2`, return `0`. Copy the input into `arr`; allocate a same-size buffer `buf`.
2. For `shift` in {0, 8, 16, 24}:
   1. Histogram the 256 possible values of the digit `(v >> shift) & 0xFF`.
   2. Exclusive prefix-sum the histogram into `starts[d]` — the first output slot for digit `d`.
   3. Scatter every element (in current order) to `buf[starts[digit]++]` — stability preserves the order established by lower digits.
   4. Swap `arr` and `buf`.
3. `arr` is now fully sorted; scan adjacent pairs for the maximum gap.

### Complexity
- **Time:** O(d·(n + b)) with d = 4 passes and b = 256 buckets → O(n); ~4·10⁵ scatter operations for n = 10⁵.
- **Space:** O(n + b) — the scatter buffer plus the fixed 256-entry count/start tables.

### Code
```go
func radixSort(nums []int) int {
	if len(nums) < 2 {
		return 0 // no pair to compare
	}
	arr := make([]int, len(nums)) // working copy (also keeps input untouched)
	copy(arr, nums)
	buf := make([]int, len(arr))             // stable-scatter destination for each pass
	for shift := 0; shift < 32; shift += 8 { // 4 passes cover all 32 bits
		var counts [256]int
		for _, v := range arr {
			counts[(v>>shift)&0xFF]++ // histogram of the current 8-bit digit
		}
		pos := 0
		var starts [256]int
		for d := 0; d < 256; d++ { // exclusive prefix sums → first slot per digit
			starts[d] = pos
			pos += counts[d]
		}
		for _, v := range arr { // stable scatter: equal digits keep their order
			d := (v >> shift) & 0xFF
			buf[starts[d]] = v
			starts[d]++
		}
		arr, buf = buf, arr // sorted-by-this-digit buffer becomes the input
	}
	maxGap := 0
	for i := 1; i < len(arr); i++ {
		if gap := arr[i] - arr[i-1]; gap > maxGap { // adjacent sorted difference
			maxGap = gap
		}
	}
	return maxGap
}
```

### Dry Run
`nums = [3,6,9,1]` (Example 1). All values fit in the lowest 8 bits, so pass 1 (`shift = 0`) does all the work and passes 2–4 are stable identity shuffles:

| pass (shift) | digit of 3,6,9,1 | non-zero counts        | starts (relevant) | scatter result `buf` |
|--------------|------------------|------------------------|-------------------|----------------------|
| 0            | 3, 6, 9, 1       | c[1]=1 c[3]=1 c[6]=1 c[9]=1 | s[1]=0 s[3]=1 s[6]=2 s[9]=3 | 3→slot1, 6→slot2, 9→slot3, 1→slot0 ⇒ `[1,3,6,9]` |
| 8, 16, 24    | all 0            | c[0]=4                 | s[0]=0            | order preserved ⇒ `[1,3,6,9]` |

Final scan of `[1,3,6,9]`: gaps 2, 3, 3 → return **3** ✅.

---

## Approach 4 — Bucket Sort + Pigeonhole (Optimal)

### Intuition
You do not need the full sorted order — only the largest adjacent gap. Pigeonhole argument: n numbers spanning `[min, max]` create n−1 sorted gaps whose *average* is `(max−min)/(n−1)`, so the **maximum gap is at least `ceil((max−min)/(n−1))`**. Now partition the span into buckets exactly that wide: any two values inside the same bucket differ by at most `bucketSize − 1 < maxGap`, so the answer can never be an intra-bucket gap. Therefore only the jumps **between** buckets matter — specifically `min(next non-empty bucket) − max(previous non-empty bucket)` — and those require only each bucket's min and max, never a sort.

### Algorithm
1. If `n < 2`, return `0`. One pass for global `min` and `max`; if `min == max`, return `0` (all gaps are 0).
2. `bucketSize = ceil((max − min) / (n − 1))` (≥ 1); `bucketCount = (max − min) / bucketSize + 1`.
3. For each value `v`: bucket index `b = (v − min) / bucketSize`; update `bucketMin[b]` and `bucketMax[b]`.
4. Sweep buckets left to right with `prevMax` = max of the last non-empty bucket (start at global `min`, which lives in bucket 0):
   1. Skip empty buckets (the eventual gap simply spans across them).
   2. Candidate gap = `bucketMin[b] − prevMax`; keep the largest.
   3. Set `prevMax = bucketMax[b]`.
5. Return the largest candidate.

### Complexity
- **Time:** O(n) — three linear passes: min/max, bucket fill, and a sweep over at most n+1 buckets.
- **Space:** O(n) — the `bucketMin`/`bucketMax` arrays (≤ n+1 entries each).

### Code
```go
func bucketPigeonhole(nums []int) int {
	n := len(nums)
	if n < 2 {
		return 0 // no pair to compare
	}
	minV, maxV := nums[0], nums[0]
	for _, v := range nums { // one pass for the global extremes
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if minV == maxV {
		return 0 // all elements equal → every sorted gap is 0
	}
	// Ceil division keeps bucketSize ≥ 1 and guarantees the answer is inter-bucket.
	bucketSize := (maxV - minV + n - 2) / (n - 1)
	bucketCount := (maxV-minV)/bucketSize + 1 // enough buckets to cover [minV, maxV]
	bucketMin := make([]int, bucketCount)     // per-bucket minimum
	bucketMax := make([]int, bucketCount)     // per-bucket maximum
	for i := range bucketMin {
		bucketMin[i] = -1 // -1 marks an empty bucket (values are ≥ 0 by constraint)
		bucketMax[i] = -1
	}
	for _, v := range nums {
		b := (v - minV) / bucketSize // bucket index of this value
		if bucketMin[b] == -1 || v < bucketMin[b] {
			bucketMin[b] = v
		}
		if bucketMax[b] == -1 || v > bucketMax[b] {
			bucketMax[b] = v
		}
	}
	maxGap := 0
	prevMax := minV // max of the last non-empty bucket seen (bucket 0 holds minV)
	for b := 0; b < bucketCount; b++ {
		if bucketMin[b] == -1 {
			continue // empty bucket — the gap simply spans across it
		}
		if gap := bucketMin[b] - prevMax; gap > maxGap { // inter-bucket gap
			maxGap = gap
		}
		prevMax = bucketMax[b] // this bucket's max feeds the next gap
	}
	return maxGap
}
```

### Dry Run
`nums = [3,6,9,1]` (Example 1): `n = 4`, `minV = 1`, `maxV = 9`; `bucketSize = ceil(8/3) = (8 + 2) / 3 = 3`; `bucketCount = 8/3 + 1 = 3`. Bucket of `v` is `(v−1)/3`:

**Fill phase:**

| v | b = (v−1)/3 | bucketMin after       | bucketMax after       |
|---|-------------|-----------------------|-----------------------|
| 3 | 0           | [3, −1, −1]           | [3, −1, −1]           |
| 6 | 1           | [3, 6, −1]            | [3, 6, −1]            |
| 9 | 2           | [3, 6, 9]             | [3, 6, 9]             |
| 1 | 0           | [**1**, 6, 9]         | [3, 6, 9]             |

**Sweep phase** (`prevMax` starts at `minV = 1`):

| b | bucket (min, max) | gap = bucketMin[b] − prevMax | maxGap after | prevMax after |
|---|-------------------|------------------------------|--------------|---------------|
| 0 | (1, 3)            | 1 − 1 = 0                    | 0            | 3             |
| 1 | (6, 6)            | 6 − 3 = 3                    | 3            | 6             |
| 2 | (9, 9)            | 9 − 6 = 3                    | 3            | 9             |

Return **3** ✅. (Example 2 `[10]`: `n < 2` → **0** ✅.)

---

## Key Takeaways

- **The pigeonhole lower bound `maxGap ≥ ceil((max−min)/(n−1))` is the whole trick**: pick buckets of exactly that width and the answer is forced to be an *inter*-bucket gap, so per-bucket (min, max) is all the "sorting" you need.
- "Linear time" requirements on integer data are a signal to reach past comparison sorts: counting sort, radix sort, and bucket techniques all sidestep the Ω(n log n) lower bound by exploiting bounded values.
- Empty buckets are not a problem — they are the *point*: wide gaps manifest as runs of empty buckets, and the `prevMax` sweep bridges them naturally.
- Use ceil division `(a + b − 1) / b` (here `(max−min + n−2) / (n−1)`) to keep `bucketSize ≥ 1` and the correctness proof airtight; a floor-based size needs a separate `max(1, …)` guard.
- Sentinel choices must respect the data range: `-1` marks empty buckets only because values are guaranteed non-negative; with arbitrary values use a separate `used []bool`.
- The successor-search brute force is a useful sanity oracle for randomized testing: it computes sorted-adjacent gaps without sorting, so it cannot share a sorting bug with the fast versions.

---

## Related Problems

- LeetCode #220 — Contains Duplicate III (bucket width chosen so answers are intra/inter-bucket by construction)
- LeetCode #539 — Minimum Time Difference (bucket/counting over a bounded value domain)
- LeetCode #912 — Sort an Array (practice ground for radix/counting sort implementations)
- LeetCode #41 — First Missing Positive (linear-time answer via value-indexed placement, same "beat the sort" spirit)
- LeetCode #274 — H-Index (counting-sort style bucketing to replace comparison sorting)
