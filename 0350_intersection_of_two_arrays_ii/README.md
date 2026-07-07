# 0350 — Intersection of Two Arrays II

> LeetCode #350 · Difficulty: Easy
> **Categories:** Array, Hash Table, Two Pointers, Binary Search, Sorting, Counting

---

## Problem Statement

Given two integer arrays `nums1` and `nums2`, return an array of their **intersection**. Each element in the result must appear as many times as it shows in **both** arrays and you may return the result in **any order**.

**Example 1:**
```
Input: nums1 = [1,2,2,1], nums2 = [2,2]
Output: [2,2]
```

**Example 2:**
```
Input: nums1 = [4,9,5], nums2 = [9,4,9,8,4]
Output: [4,9]
Explanation: [9,4] is also accepted.
```

**Constraints:**
- `1 <= nums1.length, nums2.length <= 1000`
- `0 <= nums1[i], nums2[i] <= 1000`

**Follow-up:**
- What if the given array is already sorted? How would you optimize your algorithm?
- What if `nums1`'s size is small compared to `nums2`'s size? Which algorithm is better?
- What if elements of `nums2` are stored on disk, and the memory is limited such that you cannot load all elements into the memory at once?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Hash Map (frequency counting)** — the count of each shared value is `min` of its counts in the two arrays → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers (merge walk)** — sorted arrays merge to emit one output per matched pair → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — enables the two-pointer approach and answers the "already sorted" follow-up → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Nested Scan with Used Marks (Brute Force) | O(n·m) | O(m) | Tiny inputs |
| 2 | Frequency Map (Optimal, unsorted) | O(n + m) | O(min(n,m)) | Default; great when one array is small |
| 3 | Sort + Two Pointers | O(n log n + m log m) | O(1) extra | Already sorted / disk-stream follow-ups |

---

## Approach 1 — Nested Scan with Used Marks (Brute Force)

### Intuition
Unlike #349, **duplicates count**: a value shared twice must appear twice in the answer. So each element of `nums1` should consume one matching, not-yet-consumed element of `nums2`. A boolean `used` array over `nums2` prevents claiming the same slot twice.

### Algorithm
1. Create `used`, a bool slice over `nums2`.
2. For each `x` in `nums1`, scan `nums2` for the first `j` with `nums2[j] == x` and `!used[j]`; if found, append `x` and set `used[j] = true`.

### Complexity
- **Time:** O(n·m) — each element of `nums1` may scan all of `nums2`.
- **Space:** O(m) — the `used` marks (plus the result).

### Code
```go
func bruteForce(nums1, nums2 []int) []int {
	used := make([]bool, len(nums2)) // which nums2 slots are already matched
	res := []int{}
	for _, x := range nums1 {
		for j := 0; j < len(nums2); j++ {
			if !used[j] && nums2[j] == x { // an unclaimed equal element
				res = append(res, x) // pair them up
				used[j] = true       // this slot is now consumed
				break
			}
		}
	}
	return res
}
```

### Dry Run
`nums1 = [1,2,2,1]`, `nums2 = [2,2]`, `used = [f,f]`:

| x | scan nums2 | match slot | used after | res |
|---|------------|-----------|-----------|-----|
| 1 | 2≠1, 2≠1 | none | [f,f] | [] |
| 2 | slot0=2, unused | j=0 | [t,f] | [2] |
| 2 | slot0 used, slot1=2 unused | j=1 | [t,t] | [2,2] |
| 1 | both slots used/≠ | none | [t,t] | [2,2] |

Result: **[2,2]**.

---

## Approach 2 — Frequency Map (Optimal, unsorted)

### Intuition
The multiplicity of a shared value in the result is `min(count_in_nums1, count_in_nums2)`. Build a frequency map of one array; walk the other, and for each value with a positive remaining count, emit it and decrement. When a count reaches zero the value is exhausted — which is exactly the `min` behaviour. Counting the **smaller** array minimises map memory (answering the "nums1 is small" follow-up).

### Algorithm
1. Ensure `nums1` is the smaller array (swap if needed); build `count[v]` over it.
2. For each `y` in `nums2`: if `count[y] > 0`, append `y` and decrement `count[y]`.

### Complexity
- **Time:** O(n + m) — one pass to count, one to consume.
- **Space:** O(min(n,m)) — a map over the smaller array.

### Code
```go
func hashMap(nums1, nums2 []int) []int {
	// Count the smaller array to minimise map memory.
	if len(nums1) > len(nums2) {
		nums1, nums2 = nums2, nums1
	}
	count := make(map[int]int) // value → remaining available occurrences
	for _, x := range nums1 {
		count[x]++
	}
	res := []int{}
	for _, y := range nums2 {
		if count[y] > 0 { // still have an unmatched occurrence of y
			res = append(res, y)
			count[y]-- // consume it
		}
	}
	return res
}
```

### Dry Run
`nums1 = [1,2,2,1]`, `nums2 = [2,2]`. `nums1` is larger, so it stays as the counted array (swap only happens when `len(nums1) > len(nums2)`, which is true here → count is built over `nums2 = [2,2]` and we probe `nums1`). After the swap: counted = `[2,2]` → `count = {2:2}`, probed = `[1,2,2,1]`:

| y (from probed) | count[y] > 0? | action | count | res |
|-----------------|---------------|--------|-------|-----|
| 1 | count[1]=0 → no | — | {2:2} | [] |
| 2 | count[2]=2 → yes | append 2, dec | {2:1} | [2] |
| 2 | count[2]=1 → yes | append 2, dec | {2:0} | [2,2] |
| 1 | no | — | {2:0} | [2,2] |

Result: **[2,2]**.

---

## Approach 3 — Sort + Two Pointers

### Intuition
After sorting, equal values line up. A merge walk advances the smaller side; on a match it emits the value **once** and advances **both** pointers, so a value shared k times produces exactly k outputs (the `min` of the two run lengths). This also answers the follow-ups: sorted inputs skip the map entirely, and if `nums2` lives on disk you can stream it past a sorted, in-memory `nums1`.

### Algorithm
1. Sort both arrays.
2. `i = j = 0`. While both in range:
   - `a[i] < b[j]` → `i++`.
   - `a[i] > b[j]` → `j++`.
   - equal → append the value and advance **both** pointers.

### Complexity
- **Time:** O(n log n + m log m) — the two sorts dominate; the merge is O(n + m).
- **Space:** O(1) extra beyond the sorted copies and the output.

### Code
```go
func twoPointers(nums1, nums2 []int) []int {
	a := append([]int(nil), nums1...) // copy to avoid mutating inputs
	b := append([]int(nil), nums2...)
	sort.Ints(a)
	sort.Ints(b)

	res := []int{}
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++ // advance the smaller side
		case a[i] > b[j]:
			j++
		default: // matched pair: emit once, consume from both
			res = append(res, a[i])
			i++
			j++
		}
	}
	return res
}
```

### Dry Run
`nums1 = [1,2,2,1]` → `a = [1,1,2,2]`; `nums2 = [2,2]` → `b = [2,2]`:

| i | j | a[i] | b[j] | comparison | action | res |
|---|---|------|------|------------|--------|-----|
| 0 | 0 | 1 | 2 | a<b | i++ | [] |
| 1 | 0 | 1 | 2 | a<b | i++ | [] |
| 2 | 0 | 2 | 2 | equal | append 2, i++, j++ | [2] |
| 3 | 1 | 2 | 2 | equal | append 2, i++, j++ | [2,2] |
| 4 | 2 | — | — | out of range | stop | [2,2] |

Result: **[2,2]**.

---

## Key Takeaways
- **Multiplicity ⇒ counts, not sets.** The distinguishing move from #349 is replacing the hash **set** with a hash **map of counts** and decrementing on each match; the answer count of a value is `min` of its two counts.
- **Count the smaller array** to bound memory — this is the direct answer to the "nums1 is small" follow-up.
- **Sorted inputs ⇒ two-pointer merge** with O(1) extra space; advancing both pointers on a match is what preserves multiplicity.
- **Disk-resident nums2:** sort once, then stream `nums2` from disk against an in-memory sorted `nums1`, holding only a single element of `nums2` at a time.

---

## Related Problems
- LeetCode #349 — Intersection of Two Arrays (set version, unique values)
- LeetCode #1 — Two Sum (hash-map lookup)
- LeetCode #4 — Median of Two Sorted Arrays (two-pointer merge)
- LeetCode #88 — Merge Sorted Array (two-pointer merge in place)
