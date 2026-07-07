# 0398 — Random Pick Index

> LeetCode #398 · Difficulty: Medium
> **Categories:** Hash Table, Math, Reservoir Sampling, Randomized

---

## Problem Statement

Given an integer array `nums` with possible **duplicates**, randomly output the index of a given `target` number. You can assume that the given `target` number must exist in the array.

Implement the `Solution` class:

- `Solution(int[] nums)` Initializes the object with the array `nums`.
- `int pick(int target)` Picks a random index `i` from `nums` where `nums[i] == target`. If there are multiple valid i's, then each index should have an **equal probability** of returning.

**Example 1:**

```
Input
["Solution", "pick", "pick", "pick"]
[[[1, 2, 3, 3, 3]], [3], [1], [3]]
Output
[null, 4, 0, 2]

Explanation
Solution solution = new Solution([1, 2, 3, 3, 3]);
solution.pick(3); // It should return either index 2, 3, or 4 randomly. Each index should have equal probability of returning.
solution.pick(1); // It should return 0. Since in the array only nums[0] is equal to 1.
solution.pick(3); // It should return either index 2, 3, or 4 randomly. Each index should have equal probability of returning.
```

**Constraints:**

- `1 <= nums.length <= 2 * 10^4`
- `-2^31 <= nums[i] <= 2^31 - 1`
- `target` is an integer from `nums`.
- At most `10^4` calls will be made to `pick`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Reservoir Sampling** — the optimal O(1)-space method keeps a size-1 reservoir, replacing the current pick by the c-th match with probability `1/c`, yielding a uniform choice in one scan → see [`/dsa/reservoir_sampling.md`](/dsa/reservoir_sampling.md)
- **Hash Map (grouping)** — the space-for-speed baseline groups all indices by value up front for O(1) picks → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Design of Data Structures** — this is a class-design problem: a constructor plus a `pick` method with a probabilistic contract → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Time (pick) | Space | When to use |
|---|----------|-------------|-------|-------------|
| 1 | Hash Map of Indices | O(1) | O(n) | Many `pick` calls, memory is cheap; fastest per query |
| 2 | Reservoir Sampling (Optimal Space) | O(n) | O(1) extra | Memory-constrained / streaming; the intended interview answer |

---

## Approach 1 — Hash Map of Indices

### Intuition

`pick(target)` needs a uniform index among all positions equal to `target`. If we precompute, for each value, the list of indices where it appears, then a pick is just one uniform random draw over that list — O(1) per call after an O(n) build.

### Algorithm

1. **Constructor:** scan `nums` once, appending each index `i` to `map[nums[i]]`.
2. **pick(target):** fetch `positions = map[target]`; return `positions[rand.Intn(len(positions))]`.

### Complexity

- **Time:** constructor O(n); each `pick` O(1).
- **Space:** O(n) — every index is stored, grouped by value.

### Code

```go
type SolutionHashMap struct {
	idx map[int][]int // value -> all indices where it occurs
}

func NewSolutionHashMap(nums []int) *SolutionHashMap {
	m := make(map[int][]int, len(nums))
	for i, v := range nums {
		m[v] = append(m[v], i) // group indices by their value
	}
	return &SolutionHashMap{idx: m}
}

func (s *SolutionHashMap) Pick(target int) int {
	positions := s.idx[target]              // all indices holding target
	return positions[rand.Intn(len(positions))] // uniform choice among them
}
```

### Dry Run

Example: `nums = [1,2,3,3,3]`.

| Phase | State / action |
|-------|----------------|
| build | `map = {1:[0], 2:[1], 3:[2,3,4]}` |
| pick(3) | `positions = [2,3,4]`, `rand.Intn(3)` → say 1 → return `positions[1] = 3` (any of 2/3/4 equally likely) |
| pick(1) | `positions = [0]`, `rand.Intn(1) = 0` → return `0` |
| pick(3) | `positions = [2,3,4]`, uniform draw → e.g. `2` |

Each index in `[2,3,4]` is returned with probability exactly `1/3`. ✔

---

## Approach 2 — Reservoir Sampling (Optimal Space)

### Intuition

Store only the raw array. On `pick(target)`, stream through it and count matches. Maintain a single "chosen" index; when the c-th match appears, overwrite the choice with probability `1/c`. By induction, after the full scan each match has been retained with probability exactly `1/count` — the size-1 reservoir sampling guarantee — so the result is uniform with O(1) extra memory.

**Why `1/c`?** The c-th element is kept with prob `1/c`; any earlier element survives with prob `(its earlier keep) · (1 − 1/(k+1)) · … = 1/c` after normalisation, so all `c` candidates end equally likely.

### Algorithm

1. `count = 0`, `result = -1`.
2. For each `i` with `nums[i] == target`:
   1. `count++`.
   2. With probability `1/count` (i.e. `rand.Intn(count) == 0`), set `result = i`.
3. Return `result`.

### Complexity

- **Time:** constructor O(1); each `pick` O(n) — one scan of the array.
- **Space:** O(1) extra — only the input reference plus two scalars.

### Code

```go
type SolutionReservoir struct {
	nums []int // reference to the input; no per-value preprocessing
}

func NewSolutionReservoir(nums []int) *SolutionReservoir {
	return &SolutionReservoir{nums: nums}
}

func (s *SolutionReservoir) Pick(target int) int {
	count := 0    // how many matches seen so far
	result := -1  // currently chosen index
	for i, v := range s.nums {
		if v != target {
			continue // ignore non-matching positions
		}
		count++
		// The c-th match wins the "seat" with probability 1/c, keeping the
		// distribution uniform over all matches seen so far.
		if rand.Intn(count) == 0 {
			result = i
		}
	}
	return result
}
```

### Dry Run

Example: `nums = [1,2,3,3,3]`, `pick(3)`.

| i | nums[i] | match? | count | rand.Intn(count)==0? | result |
|---|---------|--------|-------|----------------------|--------|
| 0 | 1 | no | 0 | — | -1 |
| 1 | 2 | no | 0 | — | -1 |
| 2 | 3 | yes | 1 | `Intn(1)=0` → yes (prob 1/1) | 2 |
| 3 | 3 | yes | 2 | prob 1/2 → maybe | 2 or 3 |
| 4 | 3 | yes | 3 | prob 1/3 → maybe | 2/3/4 |

Final `result ∈ {2,3,4}`, each with probability exactly `1/3`. The `main()` uniformity check (30000 picks) yields counts ≈ 9968 / 10012 / 10020 — evenly spread. ✔

---

## Key Takeaways

- **Reservoir sampling of size 1** is the canonical way to pick one uniform element from a stream in O(1) space: keep the c-th item with probability `1/c`.
- **Hash-map preprocessing** trades O(n) memory for O(1) picks — pick it when there are many queries and memory is abundant; pick reservoir sampling when memory is tight or the data is streamed.
- The `rand.Intn(count) == 0` idiom *is* "keep with probability `1/count`" — a compact, allocation-free way to express the replacement rule.
- Verify randomized solutions empirically: run thousands of picks and confirm the empirical distribution is roughly uniform (as the demo does).

---

## Related Problems

- LeetCode #382 — Linked List Random Node (size-1 reservoir over a list)
- LeetCode #384 — Shuffle an Array (uniform randomness / Fisher–Yates)
- LeetCode #528 — Random Pick with Weight (prefix sums + binary search)
- LeetCode #710 — Random Pick with Blacklist (index remapping)
