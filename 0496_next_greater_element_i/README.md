# 0496 — Next Greater Element I

> LeetCode #496 · Difficulty: Easy
> **Categories:** Array, Hash Table, Stack, Monotonic Stack

---

## Problem Statement

The **next greater element** of some element `x` in an array is the **first greater** element that is **to the right** of `x` in the same array.

You are given two **distinct 0-indexed** integer arrays `nums1` and `nums2`, where `nums1` is a subset of `nums2`.

For each `0 <= i < nums1.length`, find the index `j` such that `nums1[i] == nums2[j]` and determine the **next greater element** of `nums2[j]` in `nums2`. If there is no next greater element, then the answer for this query is `-1`.

Return *an array* `ans` *of length* `nums1.length` *such that* `ans[i]` *is the **next greater element** as described above.*

**Example 1:**

```
Input: nums1 = [4,1,2], nums2 = [1,3,4,2]
Output: [-1,3,-1]
Explanation: The next greater element for each value of nums1 is as follows:
- 4 is underlined in nums2 = [1,3,4,2]. There is no next greater element, so the answer is -1.
- 1 is underlined in nums2 = [1,3,4,2]. The next greater element is 3.
- 2 is underlined in nums2 = [1,3,4,2]. There is no next greater element, so the answer is -1.
```

**Example 2:**

```
Input: nums1 = [2,4], nums2 = [1,2,3,4]
Output: [3,-1]
Explanation: The next greater element for each value of nums1 is as follows:
- 2 is underlined in nums2 = [1,2,3,4]. The next greater element is 3.
- 4 is underlined in nums2 = [1,2,3,4]. There is no next greater element, so the answer is -1.
```

**Constraints:**

- `1 <= nums1.length <= nums2.length <= 1000`
- `0 <= nums1[i], nums2[i] <= 10^4`
- All integers in `nums1` and `nums2` are **unique**.
- All the integers of `nums1` also appear in `nums2`.

**Follow up:** Could you find an `O(nums1.length + nums2.length)` solution?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Monotonic Stack** — sweeping `nums2` once while keeping a *decreasing* stack of values that are still waiting for a larger neighbour lets each element discover its next-greater in amortised O(1); this "next greater to the right" pattern is the canonical monotonic-stack use case → see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)
- **Hash Map** — because `nums1` is a subset of `nums2`, precomputing every value's answer into a `value → nextGreater` map turns each query into an O(1) lookup, decoupling the two arrays → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(m·n) | O(1) | Tiny inputs; when clarity beats speed and n ≤ ~1000 |
| 2 | Monotonic Stack + Hash Map (Optimal) | O(m + n) | O(n) | The intended answer; meets the follow-up's linear bound |

*(m = len(nums1), n = len(nums2))*

---

## Approach 1 — Brute Force

### Intuition

Translate the definition directly. For each query value, find where it sits in `nums2`, then scan every position to its right and stop at the first strictly-greater number. If the scan runs off the end, the answer is `-1`. No cleverness — just two nested loops.

### Algorithm

1. For every value `v` in `nums1` (index `i`):
   1. Default `result[i] = -1`.
   2. Linear-scan `nums2` to find the index `j` where `nums2[j] == v`.
   3. Walk `k` from `j+1` to the end; the first `nums2[k] > v` becomes `result[i]`; break.
2. Return `result`.

### Complexity

- **Time:** O(m·n) — each of the `m` queries may locate its value and then scan the whole of `nums2`, so up to `m·n` comparisons.
- **Space:** O(1) — only loop indices beyond the required output slice.

### Code

```go
func bruteForce(nums1 []int, nums2 []int) []int {
	result := make([]int, len(nums1)) // one answer slot per query
	for i, v := range nums1 {         // handle each query value independently
		result[i] = -1 // default: assume no greater element exists
		j := 0
		for j < len(nums2) && nums2[j] != v { // locate v inside nums2
			j++
		}
		// walk to the right of v looking for the first strictly-greater number
		for k := j + 1; k < len(nums2); k++ {
			if nums2[k] > v {
				result[i] = nums2[k] // first greater to the right wins
				break                // stop at the very first one
			}
		}
	}
	return result
}
```

### Dry Run

Example 1: `nums1 = [4,1,2]`, `nums2 = [1,3,4,2]`.

| i | v | index j in nums2 | scan right of j | first `> v` | result[i] |
|---|---|------------------|-----------------|-------------|-----------|
| 0 | 4 | 2 (`nums2[2]=4`) | `nums2[3]=2` | none | -1 |
| 1 | 1 | 0 (`nums2[0]=1`) | `nums2[1]=3` | 3 | 3 |
| 2 | 2 | 3 (`nums2[3]=2`) | (end of array) | none | -1 |

Result: `[-1, 3, -1]` ✔

---

## Approach 2 — Monotonic Stack + Hash Map (Optimal)

### Intuition

Solve the harder problem once — the next-greater element for **every** value of `nums2` — then read off the answers for `nums1`. Maintain a stack of `nums2` values that have not yet found a greater neighbour to their right; kept in processing order this stack is **strictly decreasing** from bottom to top. When a new value `cur` appears, it is exactly the next-greater element for every stacked value smaller than it, so pop them all and pair each with `cur`. Each value is pushed once and popped at most once, giving a linear sweep. A `value → nextGreater` map then answers each `nums1` query in O(1), which is legal because `nums1 ⊆ nums2`.

### Algorithm

1. Create an empty map `nextGreater` and an empty `stack`.
2. For each `cur` in `nums2`:
   1. While `stack` is non-empty and `cur > stack.top`: pop `top`, set `nextGreater[top] = cur`.
   2. Push `cur`.
3. Values still on the stack have no greater element to the right — leave them out of the map.
4. For each `v` in `nums1`: `result = nextGreater[v]` if present, else `-1`.

### Complexity

- **Time:** O(m + n) — the `nums2` sweep is O(n) (each element pushed/popped once), and each of the `m` queries is an O(1) map lookup.
- **Space:** O(n) — the map holds up to `n` pairings and the stack up to `n` values.

### Code

```go
func monotonicStack(nums1 []int, nums2 []int) []int {
	nextGreater := make(map[int]int, len(nums2)) // value -> next greater element
	stack := make([]int, 0, len(nums2))          // decreasing stack of "unanswered" values

	for _, cur := range nums2 {
		// cur resolves every smaller value sitting on top of the stack
		for len(stack) > 0 && cur > stack[len(stack)-1] {
			top := stack[len(stack)-1]   // the value that was waiting
			stack = stack[:len(stack)-1] // pop it
			nextGreater[top] = cur       // cur is its first greater-to-the-right
		}
		stack = append(stack, cur) // cur now waits for its own next greater
	}
	// values remaining on the stack never found a greater element — omit them

	result := make([]int, len(nums1))
	for i, v := range nums1 {
		if g, ok := nextGreater[v]; ok { // O(1) lookup for this query
			result[i] = g
		} else {
			result[i] = -1 // no greater element existed
		}
	}
	return result
}
```

### Dry Run

Example 1: `nums2 = [1,3,4,2]` (build the map), then answer `nums1 = [4,1,2]`.

| Step | cur | pops (top → set nextGreater) | stack after push | map so far |
|------|-----|------------------------------|------------------|------------|
| 1 | 1 | — | `[1]` | `{}` |
| 2 | 3 | pop 1 → `nextGreater[1]=3` | `[3]` | `{1:3}` |
| 3 | 4 | pop 3 → `nextGreater[3]=4` | `[4]` | `{1:3, 3:4}` |
| 4 | 2 | (2 < 4, no pop) | `[4,2]` | `{1:3, 3:4}` |

Leftover stack `[4,2]` → values 4 and 2 have no next-greater. Now answer queries:

| v | in map? | result |
|---|---------|--------|
| 4 | no | -1 |
| 1 | yes → 3 | 3 |
| 2 | no | -1 |

Result: `[-1, 3, -1]` ✔

---

## Key Takeaways

- **"Next greater / smaller element to the right"** is the flagship monotonic-stack pattern: keep a stack of unresolved candidates, and let each incoming element resolve everything it dominates. Amortised O(n) because every element enters and leaves the stack once.
- **Precompute for the superset, index with a map.** When queries come from a subset of a larger array, solve the general version once over the superset and store answers in a hash map so each query is O(1). This decouples query count from work.
- Stack **direction encodes the question**: a *decreasing* stack surfaces next-*greater* elements; flip the comparison for next-*smaller*.
- The values left on the stack at the end are precisely those with **no** qualifying neighbour — a free way to detect the "-1" cases.

---

## Related Problems

- LeetCode #503 — Next Greater Element II (circular array; same stack, indices mod n)
- LeetCode #556 — Next Greater Element III (next greater permutation of digits)
- LeetCode #739 — Daily Temperatures (next warmer day = next greater to the right)
- LeetCode #901 — Online Stock Span (monotonic stack over a stream)
- LeetCode #84 — Largest Rectangle in Histogram (next smaller on both sides)
