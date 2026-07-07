# 0349 — Intersection of Two Arrays

> LeetCode #349 · Difficulty: Easy
> **Categories:** Array, Hash Table, Two Pointers, Binary Search, Sorting, Set

---

## Problem Statement

Given two integer arrays `nums1` and `nums2`, return an array of their **intersection**. Each element in the result must be **unique** and you may return the result in **any order**.

**Example 1:**
```
Input: nums1 = [1,2,2,1], nums2 = [2,2]
Output: [2]
```

**Example 2:**
```
Input: nums1 = [4,9,5], nums2 = [9,4,9,8,4]
Output: [9,4]
Explanation: [4,9] is also accepted.
```

**Constraints:**
- `1 <= nums1.length, nums2.length <= 1000`
- `0 <= nums1[i], nums2[i] <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Hash Set** — O(1) membership testing to check "is this value in the other array?" → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers (merge walk)** — after sorting, a merge-style scan finds shared values without extra space → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — the two-pointer approach relies on both arrays being ordered → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Nested Scan (Brute Force) | O(n·m) | O(min(n,m)) | Tiny inputs; no set allowed |
| 2 | Hash Set (Optimal, unsorted) | O(n + m) | O(n) | Default fastest answer |
| 3 | Sort + Two Pointers | O(n log n + m log m) | O(1) extra | Sorted output / limited memory |

---

## Approach 1 — Nested Scan (Brute Force)

### Intuition
The intersection is "values present in both arrays." The most literal check is: for each value in `nums1`, linearly search `nums2`. Since each common value must appear only once, track which values we have already emitted in a set.

### Algorithm
1. For each `x` in `nums1`: skip if already emitted, else linear-search `nums2`.
2. On a hit, append `x` to the result and mark it emitted.

### Complexity
- **Time:** O(n·m) — every element of `nums1` may scan all of `nums2`.
- **Space:** O(min(n,m)) — the emitted-set and result.

### Code
```go
func bruteForce(nums1, nums2 []int) []int {
	emitted := make(map[int]bool) // values already placed in the result
	res := []int{}
	for _, x := range nums1 {
		if emitted[x] { // already output this common value
			continue
		}
		for _, y := range nums2 { // linear search nums2 for x
			if x == y {
				res = append(res, x) // x is in both arrays
				emitted[x] = true    // never emit it again
				break
			}
		}
	}
	return res
}
```

### Dry Run
`nums1 = [1,2,2,1]`, `nums2 = [2,2]`:

| x | emitted[x]? | found in nums2? | action | res |
|---|-------------|-----------------|--------|-----|
| 1 | no | scan [2,2] → no | — | [] |
| 2 | no | 2==2 → yes | append 2, mark | [2] |
| 2 | yes | skip | — | [2] |
| 1 | no | scan → no | — | [2] |

Result: **[2]**.

---

## Approach 2 — Hash Set (Optimal, unsorted)

### Intuition
Membership testing is what a hash set is for. Dump `nums1` into a set, then a single pass over `nums2` asks "is this value in `nums1`?" in O(1). Deleting the value on the first hit keeps the result duplicate-free even if `nums2` repeats it.

### Algorithm
1. Build a set from `nums1`.
2. For each `y` in `nums2`: if `y` is in the set, append `y` and delete it from the set.

### Complexity
- **Time:** O(n + m) — one pass to build the set, one to probe it.
- **Space:** O(n) — the set of `nums1`'s distinct values.

### Code
```go
func hashSet(nums1, nums2 []int) []int {
	set := make(map[int]bool) // distinct values of nums1
	for _, x := range nums1 {
		set[x] = true
	}
	res := []int{}
	for _, y := range nums2 {
		if set[y] { // y appears in nums1 → it is in the intersection
			res = append(res, y)
			delete(set, y) // remove so a repeated y in nums2 is not re-added
		}
	}
	return res
}
```

### Dry Run
`nums1 = [1,2,2,1]`, `nums2 = [2,2]`:

| step | set | y | y in set? | action | res |
|------|-----|---|-----------|--------|-----|
| build | {1, 2} | — | — | — | [] |
| probe | {1, 2} | 2 | yes | append 2, delete 2 | [2] |
| probe | {1} | 2 | no (deleted) | — | [2] |

Result: **[2]**.

---

## Approach 3 — Sort + Two Pointers

### Intuition
Once both arrays are sorted, a merge-style walk finds equal values: advance the pointer at the smaller value; when they match, record it and skip past all copies of that value in **both** arrays so it is emitted only once. No hash structure, and the output comes out sorted.

### Algorithm
1. Copy and sort `nums1` and `nums2`.
2. `i = j = 0`. While both in range:
   - `a[i] < b[j]` → `i++`.
   - `a[i] > b[j]` → `j++`.
   - equal → append the value, then skip all its duplicates in both arrays.

### Complexity
- **Time:** O(n log n + m log m) — dominated by the two sorts; the merge walk is O(n + m).
- **Space:** O(1) extra beyond the sorted copies and output — just the two pointers.

### Code
```go
func twoPointers(nums1, nums2 []int) []int {
	a := append([]int(nil), nums1...) // copy so we do not mutate the inputs
	b := append([]int(nil), nums2...)
	sort.Ints(a)
	sort.Ints(b)

	res := []int{}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++ // a's value too small; advance a
		case a[i] > b[j]:
			j++ // b's value too small; advance b
		default: // a[i] == b[j]: a common value
			res = append(res, a[i])
			val := a[i]
			for i < len(a) && a[i] == val { // skip all copies in a
				i++
			}
			for j < len(b) && b[j] == val { // skip all copies in b
				j++
			}
		}
	}
	return res
}
```

### Dry Run
`nums1 = [1,2,2,1]` → sorted `a = [1,1,2,2]`; `nums2 = [2,2]` → sorted `b = [2,2]`:

| i | j | a[i] | b[j] | comparison | action | res |
|---|---|------|------|------------|--------|-----|
| 0 | 0 | 1 | 2 | a<b | i++ | [] |
| 1 | 0 | 1 | 2 | a<b | i++ | [] |
| 2 | 0 | 2 | 2 | equal | append 2, skip 2s: i→4, j→2 | [2] |
| 4 | 2 | — | — | i out of range | stop | [2] |

Result: **[2]**.

---

## Key Takeaways
- **"Values in both" ⇒ set membership.** The hash-set approach is the canonical O(n+m) answer; storing the smaller array keeps space minimal.
- **Uniqueness is enforced by removal:** delete from the set (or skip duplicates in the two-pointer walk) so a value is emitted exactly once.
- **Sorted inputs unlock the two-pointer merge** — no extra hashing and O(1) auxiliary space, at the cost of sorting.
- Contrast with #350 (Intersection II), where **multiplicity matters** and you count with a frequency map instead of a set.

---

## Related Problems
- LeetCode #350 — Intersection of Two Arrays II (keep duplicates by multiplicity)
- LeetCode #1 — Two Sum (hash-set / hash-map membership)
- LeetCode #202 — Happy Number (set to detect membership/cycles)
- LeetCode #4 — Median of Two Sorted Arrays (two-pointer merge of sorted arrays)
