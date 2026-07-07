# 0228 — Summary Ranges

> LeetCode #228 · Difficulty: Easy
> **Categories:** Array, Two Pointers

---

## Problem Statement

You are given a **sorted unique** integer array `nums`.

A **range** `[a,b]` is the set of all integers from `a` to `b` (inclusive).

Return *the **smallest sorted** list of ranges that **cover all the numbers in the array exactly***. That is, each element of `nums` is covered by exactly one of the ranges, and there is no integer `x` such that `x` is in one of the ranges but not in `nums`.

Each range `[a,b]` in the list should be output as:
- `"a->b"` if `a != b`
- `"a"` if `a == b`

**Example 1:**
```
Input: nums = [0,1,2,4,5,7]
Output: ["0->2","4->5","7"]
Explanation: The ranges are:
[0,2] --> "0->2"
[4,5] --> "4->5"
[7,7] --> "7"
```

**Example 2:**
```
Input: nums = [0,2,3,4,6,8,9]
Output: ["0","2->4","6","8->9"]
Explanation: The ranges are:
[0,0] --> "0"
[2,4] --> "2->4"
[6,6] --> "6"
[8,9] --> "8->9"
```

**Constraints:**
- `0 <= nums.length <= 20`
- `-2³¹ <= nums[i] <= 2³¹ - 1`
- All the values of `nums` are **unique**.
- `nums` is sorted in ascending order.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Yandex     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — an anchor at each range's start and a probe extending the run of consecutive numbers → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Array Scanning** — a single linear pass detecting where consecutiveness breaks → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **String Building** — formatting each range as `"a"` or `"a->b"` → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Single pass, track start | O(n) | O(1) extra | Simplest; emit range when consecutiveness breaks |
| 2 | Two pointers | O(n) | O(1) extra | Makes the "extend window" structure explicit |

---

## Approach 1 — Single Pass, Track Range Start

### Intuition
The array is sorted and duplicate-free, so a "range" is a maximal run of
consecutive integers. Walk left to right remembering where the current run
started. The run continues while `nums[i] == nums[i-1]+1`. The instant that
fails (or we reach the end), the run `[start .. nums[i-1]]` is complete — format
it and begin a new run at `nums[i]`.

### Algorithm
1. If `nums` is empty, return an empty list.
2. Set `start = nums[0]`.
3. For `i` from 1 to `n` (inclusive): if `i == n` **or** `nums[i] != nums[i-1]+1`,
   the run ending at `nums[i-1]` is done:
   - If `start == nums[i-1]`, append `"start"`; else append `"start->nums[i-1]"`.
   - If `i < n`, set `start = nums[i]`.
4. Return the collected strings.

### Complexity
- **Time:** O(n) — one pass; each range formatted once.
- **Space:** O(1) extra — ignoring the required output list.

### Code
```go
func summaryRanges(nums []int) []string {
	res := []string{}
	n := len(nums)
	if n == 0 {
		return res // no numbers → no ranges
	}

	start := nums[0] // first element of the run currently being built
	for i := 1; i <= n; i++ {
		// A run ends when we run off the end OR the next value is not consecutive.
		if i == n || nums[i] != nums[i-1]+1 {
			if start == nums[i-1] { // single-element range "a"
				res = append(res, strconv.Itoa(start))
			} else { // multi-element range "a->b"
				res = append(res, strconv.Itoa(start)+"->"+strconv.Itoa(nums[i-1]))
			}
			if i < n { // begin the next run at the current (non-consecutive) value
				start = nums[i]
			}
		}
	}
	return res
}
```

### Dry Run
`nums = [0,1,2,4,5,7]`, `start = 0`:

| i | nums[i] | nums[i-1] | Break? | Emitted | new start |
|---|---------|-----------|--------|---------|-----------|
| 1 | 1       | 0         | no (1==0+1) | —          | 0 |
| 2 | 2       | 1         | no (2==1+1) | —          | 0 |
| 3 | 4       | 2         | yes (4≠2+1) | "0->2"     | 4 |
| 4 | 5       | 4         | no (5==4+1) | —          | 4 |
| 5 | 7       | 5         | yes (7≠5+1) | "4->5"     | 7 |
| 6 | (end)   | 7         | yes (i==n)  | "7"        | — |

Result: `["0->2","4->5","7"]`. ✅

---

## Approach 2 — Two Pointers (Explicit Window Extension)

### Intuition
Same runs, framed as a window. Fix the left edge `i` at a range's start, then
push a right edge `j` forward while `nums[j+1] == nums[j]+1`. When `j` can go no
further, `[i..j]` is a maximal range; format it and restart with `i = j+1`. This
makes the "extend the window" structure explicit.

### Algorithm
1. Set `i = 0`.
2. While `i < n`: set `j = i`; advance `j` while `j+1 < n` and `nums[j+1] == nums[j]+1`.
3. If `i == j`, emit `"nums[i]"`; else emit `"nums[i]->nums[j]"`.
4. Set `i = j + 1` and repeat.
5. Return the list.

### Complexity
- **Time:** O(n) — `i` and `j` together advance at most `n` steps.
- **Space:** O(1) extra — ignoring the output list.

### Code
```go
func summaryRangesTwoPointers(nums []int) []string {
	res := []string{}
	n := len(nums)
	i := 0
	for i < n {
		j := i // extend the right edge as far as consecutiveness holds
		for j+1 < n && nums[j+1] == nums[j]+1 {
			j++
		}
		if i == j { // window is a single element
			res = append(res, strconv.Itoa(nums[i]))
		} else { // window spans nums[i]..nums[j]
			res = append(res, strconv.Itoa(nums[i])+"->"+strconv.Itoa(nums[j]))
		}
		i = j + 1 // jump past the finished range
	}
	return res
}
```

### Dry Run
`nums = [0,1,2,4,5,7]`:

| Outer i | Inner j stops at | Window | Emitted | next i |
|---------|------------------|--------|---------|--------|
| 0       | 2 (nums 0,1,2)   | [0..2] | "0->2"  | 3      |
| 3       | 4 (nums 4,5)     | [4..5] | "4->5"  | 5      |
| 5       | 5 (nums 7)       | [7..7] | "7"     | 6      |

Result: `["0->2","4->5","7"]`. ✅

---

## Key Takeaways
- A maximal run of consecutive sorted-unique integers is detected by the simple
  test `nums[i] == nums[i-1] + 1`.
- Iterating one past the end (`i <= n`) lets a single branch flush the final
  range without duplicating logic after the loop.
- The two-pointer framing (anchor + probe) is the same algorithm expressed as an
  explicit sliding window — useful vocabulary for interviews.
- Format single-element ranges as `"a"` and multi-element as `"a->b"`.

---

## Related Problems
- LeetCode #163 — Missing Ranges (complement: report the gaps)
- LeetCode #56 — Merge Intervals (coalescing overlapping ranges)
- LeetCode #57 — Insert Interval (range manipulation)
- LeetCode #352 — Data Stream as Disjoint Intervals (streaming version)
